package indexer

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"

	"github.com/cosmos/gogoproto/proto"
	CommitmentTypes "github.com/elys-network/elys/indexer/txs/commitments"
	indexerTypes "github.com/elys-network/elys/indexer/types"
)

// AppI defines the interface that the app must implement
type AppI interface {
	InterfaceRegistry() types.InterfaceRegistry
}

type queueItem struct {
	ctx               sdk.Context
	proc              indexerTypes.Processor
	includedAddresses []string
}

type eventItem struct {
	ctx       sdk.Context
	eventType string
	proc      indexerTypes.EventProcessor
	addresses []string
}

var (
	txChan           chan queueItem
	eventChan        chan eventItem
	database         *LMDBManager
	totalIndexLength uint64
	once             sync.Once
	workerDone       chan struct{}
	eventWorkerDone  chan struct{}
	app              AppI
	workerReady      sync.WaitGroup
	dbReady          chan struct{}
)

// Init initializes the indexer with a single worker and stores the app interface
func Init(a AppI) {
	once.Do(func() {
		app = a
		dbReady = make(chan struct{})
		workerReady.Add(2) // Add one more for event worker

		go initDatabase()

		txChan = make(chan queueItem, 10000)
		eventChan = make(chan eventItem, 10000)
		workerDone = make(chan struct{})
		eventWorkerDone = make(chan struct{})

		go func() {
			<-dbReady // Wait for the database to be ready
			go worker()
			go eventWorker()
			workerReady.Done() // Signal that the workers are ready
			workerReady.Done()
		}()

		// Wait for both the database and the workers to be ready
		<-dbReady
		workerReady.Wait()
	})
}

func initDatabase() {
	db, err := NewLMDBManager("./lmdb-data", &totalIndexLength)
	if err != nil {
		panic(err)
	}
	database = db
	data, err := db.GetTxsByAddress("elys1093h5gs0wz3rrm78zdrfzul2mdh654d95mhnj9")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Retrieved %d transactions\n", len(data))
		for _, tx := range data {
			fmt.Println(tx)
			transactionType, TransactionData, err := ParseTransaction(tx)

			if err != nil {
				fmt.Println(err)
			}

			switch transactionType {
			case "/elys.commitment.MsgStake":
				if stakeData, ok := TransactionData.(CommitmentTypes.MsgStake); ok {
					fmt.Printf("Stake amount: %s %s\n", stakeData.Token.Amount, stakeData.Token.Denom)
				}

			}
		}
	}

	close(dbReady) // Signal that the database is ready
}

// StopIndexer stops the indexer worker
func StopIndexer() {
	close(txChan)
	close(eventChan)
	<-workerDone
	<-eventWorkerDone
}

// worker processes transactions from the channel
func worker() {
	defer close(workerDone)
	for item := range txChan {
		processTransactionInternal(item.ctx, item.proc, item.includedAddresses)
	}
}

// eventWorker processes background events from the channel
func eventWorker() {
	defer close(eventWorkerDone)
	for event := range eventChan {
		processEventInternal(event)
	}
}

// QueueTransaction sends the transaction context and processor to the worker
func QueueTransaction(ctx sdk.Context, proc indexerTypes.Processor, addresses []string) {
	item := queueItem{
		ctx:               ctx,
		proc:              proc,
		includedAddresses: addresses,
	}

	select {
	case txChan <- item:
		fmt.Println("Processing")
	default:
		fmt.Println("Transaction indexer channel is full, waiting to enqueue...")
		txChan <- item // This will block until there's space in the channel
	}
}

// QueueEvent sends background events to the event worker
func QueueEvent(ctx sdk.Context, eventType string, proc indexerTypes.EventProcessor, addresses []string) {
	event := eventItem{
		ctx:       ctx,
		eventType: eventType,
		proc:      proc,
		addresses: addresses,
	}

	select {
	case eventChan <- event:
		fmt.Printf("Processing event: %s\n", eventType)
	default:
		fmt.Println("Event channel is full, waiting to enqueue...")
		eventChan <- event
	}
}

func processEventInternal(event eventItem) {
	baseEvent := indexerTypes.BaseEvent{
		BlockTime:   event.ctx.BlockTime(),
		EventType:   event.eventType,
		BlockHeight: event.ctx.BlockHeight(),
	}

	res, err := event.proc.Process(database, baseEvent)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)
}

func processTransactionInternal(ctx sdk.Context, proc indexerTypes.Processor, includingAddresses []string) {
	txBytes := ctx.TxBytes()
	if len(txBytes) == 0 {
		fmt.Println("No transaction bytes found in context")
		return
	}

	txChecksum := sha256.Sum256(txBytes)
	txHash := hex.EncodeToString(txChecksum[:])

	blockHeight := ctx.BlockHeight()
	blockTime := ctx.BlockTime()
	gasUsed := ctx.GasMeter().GasConsumed()

	txConfig := tx.NewTxConfig(codec.NewProtoCodec(app.InterfaceRegistry()), tx.DefaultSignModes)
	decodedTx, err := txConfig.TxDecoder()(txBytes)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	msg := decodedTx.GetMsgs()[0]

	// Get signer address, handling both legacy and non-legacy messages
	var sender sdk.AccAddress
	if legacyMsg, ok := msg.(sdk.LegacyMsg); ok {
		sender = legacyMsg.GetSigners()[0]
	} else {
		// For non-legacy messages, get signer from tx signers
		sigTx, ok := decodedTx.(signing.SigVerifiableTx)
		if !ok {
			fmt.Println("tx does not implement SigVerifiableTx")
			return
		}
		signers, err := sigTx.GetSigners()
		if err != nil {
			fmt.Println("failed to get signers:", err)
			return
		}
		if len(signers) == 0 {
			fmt.Println("no signers found")
			return
		}
		sender = signers[0]
	}

	feeTx, ok := decodedTx.(sdk.FeeTx)
	if !ok {
		fmt.Println("tx is not a sdk.FeeTx")
		return
	}

	memoTx, ok := decodedTx.(sdk.TxWithMemo)
	if !ok {
		fmt.Println("tx is not a sdk.TxWithMemo")
		return
	}

	memo := memoTx.GetMemo()
	fees := feeTx.GetFee()
	gasLimit := feeTx.GetGas()

	var feeDetails []indexerTypes.FeeDetail
	for _, fee := range fees {
		feeDetails = append(feeDetails, indexerTypes.FeeDetail{
			Amount: fee.Amount.String(),
			Denom:  fee.Denom,
		})
	}

	baseTx := indexerTypes.BaseTransaction{
		BlockTime:         blockTime,
		Author:            sender.String(),
		IncludedAddresses: includingAddresses,
		BlockHeight:       blockHeight,
		TxHash:            txHash,
		TxType:            "/" + proto.MessageName(msg),
		Fees:              feeDetails,
		GasUsed:           strconv.FormatUint(gasUsed, 10),
		GasLimit:          strconv.FormatUint(gasLimit, 10),
		Memo:              memo,
	}

	fmt.Println(baseTx)

	res, err := proc.Process(database, baseTx)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)
}

// retryProcessing attempts to reprocess a transaction after a delay
func retryProcessing(ctx sdk.Context, proc indexerTypes.Processor, includingAddresses []string) {
	go func() {
		time.Sleep(5 * time.Second) // Wait for 5 seconds before retrying
		QueueTransaction(ctx, proc, includingAddresses)
	}()
}
