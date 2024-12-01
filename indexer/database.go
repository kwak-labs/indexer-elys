package indexer

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"

	"github.com/bmatsuo/lmdb-go/lmdb"

	indexerTypes "github.com/elys-network/elys/indexer/types"
)

// LMDBManager handles LMDB operations for storing and retrieving transactions
type LMDBManager struct {
	env              *lmdb.Env
	eventDB          lmdb.DBI
	addressDB        lmdb.DBI
	eventCountDB     lmdb.DBI
	path             string
	totalIndexLength *uint64
}

// NewLMDBManager creates a new LMDB manager
func NewLMDBManager(path string, totalIndexLength *uint64) (*LMDBManager, error) {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	// Set up LMDB environment
	env, err := lmdb.NewEnv()
	if err != nil {
		return nil, err
	}

	if err := env.SetMaxDBs(3); err != nil {
		return nil, err
	}

	// Start with 1GB, we'll increase it later if needed
	if err := env.SetMapSize(1 << 30); err != nil {
		return nil, err
	}

	if err := env.Open(path, 0, 0644); err != nil {
		return nil, err
	}

	manager := &LMDBManager{env: env, path: path, totalIndexLength: totalIndexLength}

	// Initialize databases
	err = env.Update(func(txn *lmdb.Txn) error {
		var err error
		if manager.eventDB, err = txn.OpenDBI("txs", lmdb.Create); err != nil {
			return err
		}
		if manager.addressDB, err = txn.OpenDBI("addresses", lmdb.Create|lmdb.DupSort); err != nil {
			return err
		}
		if manager.eventCountDB, err = txn.OpenDBI("txcount", lmdb.Create); err != nil {
			return err
		}

		// Get the current transaction count
		countBytes, err := txn.Get(manager.eventCountDB, []byte("count"))
		if err == nil {
			*totalIndexLength = binary.LittleEndian.Uint64(countBytes)
		} else if lmdb.IsNotFound(err) {
			*totalIndexLength = 0
		} else {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return manager, nil
}

// CheckAndResizeIfNeeded increases the database size if it's getting full
func (m *LMDBManager) CheckAndResizeIfNeeded() error {
	info, err := m.env.Info()
	if err != nil {
		return err
	}

	usedSpace := uint64(info.LastPNO) * uint64(os.Getpagesize())
	availableSpace := uint64(info.MapSize) - usedSpace

	// Double the size if less than 20% is available
	if availableSpace < uint64(info.MapSize)/5 {
		newSize := info.MapSize * 2
		if err := m.env.SetMapSize(newSize); err != nil {
			// If resizing fails, close and reopen the environment
			m.env.Close()
			if env, err := lmdb.NewEnv(); err == nil {
				if err := env.SetMaxDBs(3); err == nil {
					if err := env.SetMapSize(newSize); err == nil {
						if err := env.Open(m.path, 0, 0644); err == nil {
							m.env = env
							fmt.Printf("Resized database to %d bytes\n", newSize)
							return nil
						}
					}
				}
			}
			return err
		}
	}

	return nil
}

// ProcessNewTx adds a new transaction to the database
func (m *LMDBManager) ProcessNewTx(tx indexerTypes.GenericTransaction, address string) error {
	if err := m.CheckAndResizeIfNeeded(); err != nil {
		return err
	}

	return m.env.Update(func(txn *lmdb.Txn) error {
		// Increment the total index length

		fmt.Printf("Before increment: %d\n", *m.totalIndexLength)
		*m.totalIndexLength++
		fmt.Printf("After increment: %d\n", *m.totalIndexLength)
		count := *m.totalIndexLength

		// Store the new count in the database
		countBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(countBytes, count)
		if err := txn.Put(m.eventCountDB, []byte("count"), countBytes, 0); err != nil {
			return fmt.Errorf("error storing new count: %v", err)
		}

		fmt.Printf("New Count: %d\n", count)

		// Store the transaction
		txBytes, err := json.Marshal(tx)
		if err != nil {
			return err
		}

		indexBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(indexBytes, count)
		if err := txn.Put(m.eventDB, indexBytes, txBytes, 0); err != nil {
			return err
		}

		// Update the address index for the main address
		if err := txn.Put(m.addressDB, []byte(address), indexBytes, 0); err != nil {
			return err
		}

		// Update the address index for all included addresses
		for _, includedAddress := range tx.BaseTransaction.IncludedAddresses {
			if err := txn.Put(m.addressDB, []byte(includedAddress), indexBytes, 0); err != nil {
				return err
			}
		}

		return nil
	})
}

func (m *LMDBManager) ProcessNewEvent(event indexerTypes.GenericEvent, address string) error {
	if err := m.CheckAndResizeIfNeeded(); err != nil {
		return err
	}

	return m.env.Update(func(txn *lmdb.Txn) error {
		// Increment the total index length
		*m.totalIndexLength++
		count := *m.totalIndexLength

		// Store the new count
		countBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(countBytes, count)
		if err := txn.Put(m.eventCountDB, []byte("count"), countBytes, 0); err != nil {
			return fmt.Errorf("error storing new count: %v", err)
		}

		// Store the event
		eventBytes, err := json.Marshal(event)
		if err != nil {
			return err
		}

		indexBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(indexBytes, count)
		if err := txn.Put(m.eventDB, indexBytes, eventBytes, 0); err != nil {
			return err
		}

		// Update the address index for the main address
		if err := txn.Put(m.addressDB, []byte(address), indexBytes, 0); err != nil {
			return err
		}

		// Update the address index for all included addresses
		for _, includedAddress := range event.BaseEvent.IncludedAddresses {
			if err := txn.Put(m.addressDB, []byte(includedAddress), indexBytes, 0); err != nil {
				return err
			}
		}

		return nil
	})
}

// GetTxCount returns the total number of transactions
func (m *LMDBManager) GetTxCount() uint64 {
	return *m.totalIndexLength
}

// GetTxByIndex retrieves a transaction by its index
func (m *LMDBManager) GetTxByIndex(index uint64) (indexerTypes.GenericTransaction, error) {
	var tx indexerTypes.GenericTransaction
	err := m.env.View(func(txn *lmdb.Txn) error {
		indexBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(indexBytes, index)
		txBytes, err := txn.Get(m.eventDB, indexBytes)
		if err != nil {
			return err
		}
		return json.Unmarshal(txBytes, &tx)
	})
	return tx, err
}

// GetTxsByAddress retrieves all transactions for a given address
func (m *LMDBManager) GetTxsByAddress(address string) ([]indexerTypes.GenericTransaction, error) {
	var txs []indexerTypes.GenericTransaction
	err := m.env.View(func(txn *lmdb.Txn) error {
		cursor, err := txn.OpenCursor(m.addressDB)
		if err != nil {
			return fmt.Errorf("error opening cursor: %v", err)
		}
		defer cursor.Close()

		_, value, err := cursor.Get([]byte(address), nil, lmdb.SetKey)
		if lmdb.IsNotFound(err) {
			return nil // No transactions found for this address
		} else if err != nil {
			return fmt.Errorf("error in initial cursor.Get: %v", err)
		}

		for {
			index := binary.BigEndian.Uint64(value)
			tx, err := m.GetTxByIndex(index)
			if err != nil {
				return fmt.Errorf("error getting transaction by index %d: %v", index, err)
			}
			txs = append(txs, tx)

			_, value, err = cursor.Get(nil, nil, lmdb.NextDup)
			if lmdb.IsNotFound(err) {
				// Reached the end of the transactions
				break
			} else if err != nil {
				return fmt.Errorf("error in cursor.Get for NextDup: %v", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return txs, nil
}

// Close shuts down the LMDB environment
func (m *LMDBManager) Close() error {
	m.env.Close()
	return nil
}
