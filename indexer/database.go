// Package indexer provides functionality for indexing blockchain transactions and events
package indexer

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/bmatsuo/lmdb-go/lmdb"

	indexerTypes "github.com/elys-network/elys/indexer/types"
)

// LMDBManager handles LMDB operations for storing and retrieving records.
// It maintains five databases:
// - recordDB: Stores the actual transaction and event records
// - addressDB: Maps addresses to record indices for efficient lookups
// - recordCountDB: Tracks the total number of records in the system
// - txHashDB: Maps transaction hashes to record indices to prevent duplicates
// - eventIDDB: Maps event IDs to record indices to prevent duplicate events
type LMDBManager struct {
	env              *lmdb.Env  // LMDB environment handle
	recordDB         lmdb.DBI   // Database for storing transaction/event records
	addressDB        lmdb.DBI   // Database mapping addresses to record indices
	recordCountDB    lmdb.DBI   // Database tracking total record count
	txHashDB         lmdb.DBI   // Database mapping tx hashes to record indices
	eventIDDB        lmdb.DBI   // Database mapping event IDs to record indices
	path             string     // File system path to the LMDB data files
	totalIndexLength *uint64    // Pointer to the current total number of records
	indexMutex       sync.Mutex // Mutex to protect index operations
}

// NewLMDBManager creates and initializes a new LMDB manager instance.
// It sets up the database environment, creates necessary subdatabases,
// and loads or initializes the record count.
func NewLMDBManager(path string, totalIndexLength *uint64) (*LMDBManager, error) {
	// Ensure the database directory exists
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	// Initialize LMDB environment
	env, err := lmdb.NewEnv()
	if err != nil {
		return nil, err
	}

	// Configure environment to support 5 named databases
	if err := env.SetMaxDBs(5); err != nil {
		return nil, err
	}

	// Set initial database size to 1GB
	if err := env.SetMapSize(1 << 30); err != nil {
		return nil, err
	}

	// Open the environment with read-write permissions
	if err := env.Open(path, 0, 0644); err != nil {
		return nil, err
	}

	manager := &LMDBManager{
		env:              env,
		path:             path,
		totalIndexLength: totalIndexLength,
	}

	// Initialize the databases within a transaction
	err = env.Update(func(txn *lmdb.Txn) error {
		var err error
		// Create main record storage database
		if manager.recordDB, err = txn.OpenDBI("records", lmdb.Create); err != nil {
			return err
		}
		// Create address index database with duplicate key support
		if manager.addressDB, err = txn.OpenDBI("addresses", lmdb.Create|lmdb.DupSort); err != nil {
			return err
		}
		// Create record count tracking database
		if manager.recordCountDB, err = txn.OpenDBI("recordcount", lmdb.Create); err != nil {
			return err
		}
		// Create tx hash tracking database
		if manager.txHashDB, err = txn.OpenDBI("txhashes", lmdb.Create); err != nil {
			return err
		}
		// Create event ID tracking database
		if manager.eventIDDB, err = txn.OpenDBI("eventids", lmdb.Create); err != nil {
			return err
		}

		// Load existing record count or initialize to 0
		countBytes, err := txn.Get(manager.recordCountDB, []byte("count"))
		if err == nil {
			*totalIndexLength = binary.LittleEndian.Uint64(countBytes)
		} else if lmdb.IsNotFound(err) {
			*totalIndexLength = 0
			// Initialize count to 0 in database
			initBytes := make([]byte, 8)
			binary.LittleEndian.PutUint64(initBytes, 0)
			if err := txn.Put(manager.recordCountDB, []byte("count"), initBytes, 0); err != nil {
				return fmt.Errorf("failed to initialize count: %v", err)
			}
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

// CheckAndResizeIfNeeded monitors database usage and automatically increases
// the size when available space drops below 20%. It doubles the current size
// when more space is needed and handles the resize operation gracefully.
func (m *LMDBManager) CheckAndResizeIfNeeded() error {
	info, err := m.env.Info()
	if err != nil {
		return err
	}

	// Calculate current space usage
	usedSpace := uint64(info.LastPNO) * uint64(os.Getpagesize())
	availableSpace := uint64(info.MapSize) - usedSpace

	// Resize if less than 20% space remains
	if availableSpace < uint64(info.MapSize)/5 {
		newSize := info.MapSize * 2
		if err := m.env.SetMapSize(newSize); err != nil {
			// If direct resize fails, attempt recovery by recreating environment
			m.env.Close()
			if env, err := lmdb.NewEnv(); err == nil {
				if err := env.SetMaxDBs(4); err == nil {
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

// ProcessNewTx wraps a transaction in a GenericRecord and processes it.
// It provides a convenient way to index individual transactions.
func (m *LMDBManager) ProcessNewTx(tx indexerTypes.GenericTransaction, address string) error {
	record := indexerTypes.GenericRecord{
		Transaction: &tx,
	}
	return m.ProcessRecord(record, address)
}

// ProcessNewEvent wraps an event in a GenericRecord and processes it.
// It provides a convenient way to index individual events.
func (m *LMDBManager) ProcessNewEvent(event indexerTypes.GenericEvent, address string) error {
	record := indexerTypes.GenericRecord{
		Event: &event,
	}
	return m.ProcessRecord(record, address)
}

// ProcessRecord stores a new record (transaction or event) in the database.
// It updates the record count, stores the record data, and maintains address indices
// for both the main address and any included addresses.
// Included addresses are like recievers, so if someone recieved 100 tokens they would be Included.
func (m *LMDBManager) ProcessRecord(record indexerTypes.GenericRecord, address string) error {
	// Check for duplicate transaction if this is a transaction record
	if record.IsTransaction() {
		txHash := record.Transaction.BaseTransaction.TxHash
		exists := false
		err := m.env.View(func(txn *lmdb.Txn) error {
			_, err := txn.Get(m.txHashDB, []byte(txHash))
			if err == nil {
				exists = true
			} else if !lmdb.IsNotFound(err) {
				return err
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("error checking tx hash: %v", err)
		}
		if exists {
			return fmt.Errorf("transaction %s has already been processed", txHash)
		}
	}

	// Check for duplicate event if this is an event record
	if record.IsEvent() {
		eventID := record.Event.BaseEvent.EventID
		exists := false
		err := m.env.View(func(txn *lmdb.Txn) error {
			_, err := txn.Get(m.eventIDDB, []byte(eventID))
			if err == nil {
				exists = true
			} else if !lmdb.IsNotFound(err) {
				return err
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("error checking event ID: %v", err)
		}
		if exists {
			return fmt.Errorf("event %s has already been processed", eventID)
		}
	}

	// Ensure database has enough space
	if err := m.CheckAndResizeIfNeeded(); err != nil {
		return err
	}

	// Lock index operations
	m.indexMutex.Lock()
	defer m.indexMutex.Unlock()

	return m.env.Update(func(txn *lmdb.Txn) error {
		// Get current count from database to ensure consistency
		countBytes, err := txn.Get(m.recordCountDB, []byte("count"))
		if err != nil && !lmdb.IsNotFound(err) {
			return fmt.Errorf("error reading count: %v", err)
		}

		var count uint64
		if err == nil {
			count = binary.LittleEndian.Uint64(countBytes)
		}

		// Increment count
		count++
		*m.totalIndexLength = count

		// Store new count
		newCountBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(newCountBytes, count)
		if err := txn.Put(m.recordCountDB, []byte("count"), newCountBytes, 0); err != nil {
			return fmt.Errorf("error storing new count: %v", err)
		}

		// Serialize and store the record
		recordBytes, err := json.Marshal(record)
		if err != nil {
			return err
		}

		indexBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(indexBytes, count)
		if err := txn.Put(m.recordDB, indexBytes, recordBytes, 0); err != nil {
			return err
		}

		// Store tx hash mapping if this is a transaction
		if record.IsTransaction() {
			txHash := record.Transaction.BaseTransaction.TxHash
			if err := txn.Put(m.txHashDB, []byte(txHash), indexBytes, 0); err != nil {
				return fmt.Errorf("error storing tx hash mapping: %v", err)
			}
		}

		// Store event ID mapping if this is an event
		if record.IsEvent() {
			eventID := record.Event.BaseEvent.EventID
			if err := txn.Put(m.eventIDDB, []byte(eventID), indexBytes, 0); err != nil {
				return fmt.Errorf("error storing event ID mapping: %v", err)
			}
		}

		// Get included addresses based on record type
		var includedAddresses []string
		if record.IsTransaction() {
			includedAddresses = record.Transaction.BaseTransaction.IncludedAddresses
		} else if record.IsEvent() {
			includedAddresses = record.Event.BaseEvent.IncludedAddresses
		}

		// Create a map to track unique addresses
		uniqueAddresses := make(map[string]string)

		// Add main address if not empty
		if address != "" {
			uniqueAddresses[address] = address
		}

		// Add included addresses if not empty and not already present
		for _, addr := range includedAddresses {
			if addr != "" {
				uniqueAddresses[addr] = addr
			}
		}

		// Push the index to each address's store
		for _, addr := range uniqueAddresses {
			if err := txn.Put(m.addressDB, []byte(addr), indexBytes, 0); err != nil {
				return err
			}
		}

		return nil
	})
}

// GetRecordCount returns the current total number of records in the database
func (m *LMDBManager) GetRecordCount() uint64 {
	return *m.totalIndexLength
}

// GetRecordByIndex retrieves a specific record by its index number.
// Returns the record and any error encountered during retrieval.
func (m *LMDBManager) GetRecordByIndex(index uint64) (indexerTypes.GenericRecord, error) {
	var record indexerTypes.GenericRecord
	err := m.env.View(func(txn *lmdb.Txn) error {
		indexBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(indexBytes, index)
		recordBytes, err := txn.Get(m.recordDB, indexBytes)
		if err != nil {
			return err
		}
		return json.Unmarshal(recordBytes, &record)
	})
	return record, err
}

// GetRecordsByAddress retrieves all records associated with a given address.
// This includes both records where the address is the main address and where
// it appears in the included addresses list.
func (m *LMDBManager) GetRecordsByAddress(address string) ([]indexerTypes.GenericRecord, error) {
	var records []indexerTypes.GenericRecord
	err := m.env.View(func(txn *lmdb.Txn) error {
		cursor, err := txn.OpenCursor(m.addressDB)
		if err != nil {
			return fmt.Errorf("error opening cursor: %v", err)
		}
		defer cursor.Close()

		// Position cursor at first record for this address
		_, value, err := cursor.Get([]byte(address), nil, lmdb.SetKey)
		if lmdb.IsNotFound(err) {
			return nil // No records found for this address
		} else if err != nil {
			return fmt.Errorf("error in initial cursor.Get: %v", err)
		}

		// Iterate through all records for this address
		for {
			index := binary.BigEndian.Uint64(value)
			record, err := m.GetRecordByIndex(index)
			if err != nil {
				return fmt.Errorf("error getting record by index %d: %v", index, err)
			}
			records = append(records, record)

			// Move to next record with same address
			_, value, err = cursor.Get(nil, nil, lmdb.NextDup)
			if lmdb.IsNotFound(err) {
				// Reached the end of the records
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

	return records, nil
}

// Close properly shuts down the LMDB environment and releases resources
func (m *LMDBManager) Close() error {
	m.env.Close()
	return nil
}
