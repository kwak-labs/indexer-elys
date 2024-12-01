package indexer

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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
	EstakingTypes "github.com/elys-network/elys/indexer/txs/estaking"
	indexerTypes "github.com/elys-network/elys/indexer/types"
)

// AppI defines the interface that the app must implement
type AppI interface {
	InterfaceRegistry() types.InterfaceRegistry
}

// queueItem represents a transaction to be processed by the worker
type queueItem struct {
	ctx               sdk.Context
	proc              indexerTypes.Processor
	includedAddresses []string
}

// eventItem represents an event to be processed by the event worker
// An event are things like liquidations that happen automatically
type eventItem struct {
	ctx       sdk.Context
	eventType string
	proc      indexerTypes.EventProcessor
	addresses []string
	id        string
}

// Global variables for managing the indexer state
var (
	txChan           chan queueItem // Channel for queuing transactions
	eventChan        chan eventItem // Channel for queuing events
	database         *LMDBManager   // Database manager instance
	totalIndexLength uint64         // Total number of indexed items
	once             sync.Once      // Ensures Init is called only once
	workerDone       chan struct{}  // Channel to signal worker completion
	eventWorkerDone  chan struct{}  // Channel to signal event worker completion
	app              AppI           // Application interface instance
	workerReady      sync.WaitGroup // WaitGroup for worker initialization
	dbReady          chan struct{}  // Channel to signal database readiness
)

// Init initializes the indexer with a single worker and stores the app interface
func Init(a AppI) {
	once.Do(func() {
		app = a
		dbReady = make(chan struct{})
		workerReady.Add(2) // Add one more for event worker

		go initDatabase()

		// Initialize channels with buffer sizes
		txChan = make(chan queueItem, 10000)
		eventChan = make(chan eventItem, 10000)
		workerDone = make(chan struct{})
		eventWorkerDone = make(chan struct{})

		// Start workers after database is ready
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

// initDatabase initializes the LMDB database and performs test queries
func initDatabase() {
	var err error
	database, err = NewLMDBManager("./lmdb-data", &totalIndexLength)
	if err != nil {
		panic(err)
	}

	// Decode every data
	count := database.GetRecordCount()
	fmt.Printf("\n=== Decoding %d records ===\n", count)

	for i := uint64(1); i <= count; i++ {
		record, err := database.GetRecordByIndex(i)
		if err != nil {
			fmt.Printf("Error getting record %d: %v\n", i, err)
			continue
		}

		txType, data, err := ParseRecord(record)
		if err != nil {
			fmt.Printf("Error parsing record %d: %v\n", i, err)
			continue
		}

		// Convert to JSON and format like console.log
		jsonBytes, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling data for record %d: %v\n", i, err)
			continue
		}

		fmt.Printf("\nRecord #%d - Type: %s\n%s\n", i, txType, string(jsonBytes))
	}

	fmt.Println("\n=== Finished decoding all records ===")

	// Test query to verify database functionality
	data, err := database.GetRecordsByAddress("elys1093h5gs0wz3rrm78zdrfzul2mdh654d95mhnj9")
	if err != nil {
		panic(fmt.Errorf("database test query failed: %v", err))
	}

	fmt.Printf("Retrieved %d transactions\n", len(data))
	for _, tx := range data {
		transactionType, TransactionData, err := ParseRecord(tx)
		if err != nil {
			panic(fmt.Errorf("failed to parse record: %v", err))
		}

		// Handle specific transaction types
		switch transactionType {
		case "/elys.commitment.MsgStake":
			if stakeData, ok := TransactionData.(CommitmentTypes.MsgStake); ok {
				fmt.Printf("Stake amount: %s %s\n", stakeData.Token.Amount, stakeData.Token.Denom)
			}
		case "/elys.estaking.MsgWithdrawAllRewards":
			if rewardData, ok := TransactionData.(EstakingTypes.MsgWithdrawAllRewards); ok {
				fmt.Printf("Validators: %v\n", rewardData.Validators)
				for _, token := range rewardData.Amount {
					fmt.Printf("Reward Amount: %s %s\n", token.Amount, token.Denom)
				}
			}
		}
	}

	close(dbReady) // Signal that the database is ready
}

// StopIndexer gracefully stops the indexer workers
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

	// Try to queue transaction, wait if channel is full
	select {
	case txChan <- item:
		fmt.Println("Processing")
	default:
		fmt.Println("Transaction indexer channel is full, waiting to enqueue...")
		txChan <- item // This will block until there's space in the channel
	}
}

// QueueEvent sends background events to the event worker
func QueueEvent(ctx sdk.Context, eventType string, proc indexerTypes.EventProcessor, addresses []string, id string) {
	event := eventItem{
		ctx:       ctx,
		eventType: eventType,
		proc:      proc,
		addresses: addresses,
		id:        id,
	}

	// Try to queue event, wait if channel is full
	select {
	case eventChan <- event:
		fmt.Printf("Processing event: %s\n", eventType)
	default:
		fmt.Println("Event channel is full, waiting to enqueue...")
		eventChan <- event
	}
}

// processEventInternal handles the processing of a single event
func processEventInternal(event eventItem) {
	baseEvent := indexerTypes.BaseEvent{
		EventID:           event.id,
		IncludedAddresses: event.addresses,
		BlockTime:         event.ctx.BlockTime(),
		EventType:         event.eventType,
		BlockHeight:       event.ctx.BlockHeight(),
	}

	_, err := event.proc.Process(database, baseEvent)
	if err != nil {
		panic(fmt.Errorf("failed to process event: %v", err))
	}
}

// processTransactionInternal handles the processing of a single transaction
func processTransactionInternal(ctx sdk.Context, proc indexerTypes.Processor, includingAddresses []string) {
	// Get transaction bytes and calculate hash
	txBytes := ctx.TxBytes()
	if len(txBytes) == 0 {
		panic("no transaction bytes found in context")
	}

	txChecksum := sha256.Sum256(txBytes)
	txHash := hex.EncodeToString(txChecksum[:])

	// Get block information
	blockHeight := ctx.BlockHeight()
	blockTime := ctx.BlockTime()
	gasUsed := ctx.GasMeter().GasConsumed()

	// Decode transaction
	txConfig := tx.NewTxConfig(codec.NewProtoCodec(app.InterfaceRegistry()), tx.DefaultSignModes)
	decodedTx, err := txConfig.TxDecoder()(txBytes)
	if err != nil {
		panic(fmt.Errorf("failed to decode transaction: %v", err))
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
			panic("tx does not implement SigVerifiableTx")
		}
		signers, err := sigTx.GetSigners()
		if err != nil {
			panic(fmt.Errorf("failed to get signers: %v", err))
		}
		if len(signers) == 0 {
			panic("no signers found")
		}
		sender = signers[0]
	}

	// Extract fee information
	feeTx, ok := decodedTx.(sdk.FeeTx)
	if !ok {
		panic("tx is not a sdk.FeeTx")
	}

	// Extract memo information
	memoTx, ok := decodedTx.(sdk.TxWithMemo)
	if !ok {
		panic("tx is not a sdk.TxWithMemo")
	}

	memo := memoTx.GetMemo()
	fees := feeTx.GetFee()
	gasLimit := feeTx.GetGas()

	// Convert fees to FeeDetail structs
	var feeDetails []indexerTypes.FeeDetail
	for _, fee := range fees {
		feeDetails = append(feeDetails, indexerTypes.FeeDetail{
			Amount: fee.Amount.String(),
			Denom:  fee.Denom,
		})
	}

	// Create base transaction
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

	// Process the transaction
	_, err = proc.Process(database, baseTx)
	if err != nil {
		fmt.Printf("failed to process transaction: %v", err)
	}
}

// retryProcessing attempts to reprocess a transaction after a delay
func retryProcessing(ctx sdk.Context, proc indexerTypes.Processor, includingAddresses []string) {
	go func() {
		time.Sleep(5 * time.Second) // Wait for 5 seconds before retrying
		QueueTransaction(ctx, proc, includingAddresses)
	}()
}
