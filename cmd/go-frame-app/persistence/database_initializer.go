package persistence

import (
	"log"
	"os"

	bolt "go.etcd.io/bbolt"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lmsgprefix|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lmsgprefix|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lmsgprefix|log.Lshortfile)
}

// Storage handles the database connection and operations.
type Storage struct {
	Db *bolt.DB
}

// NewStorage opens a connection to the BoltDB database and initializes buckets.
//
// Parameters:
//   - path: The file path to the database.
//
// Returns:
//   - *Storage: The storage instance.
//   - error: An error if opening the database fails.
func NewStorage(path string) (*Storage, error) {
	InfoLogger.Println("Opening Database Connection")
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	storage := &Storage{Db: db}
	InfoLogger.Println("Database connection established")

	if err := storage.initBuckets(); err != nil {
		db.Close()
		return nil, err
	}

	return storage, nil
}

func (s *Storage) initBuckets() error {
	err := s.Db.Update(initImageBuckets)
	if err != nil {
		return err
	}
	err = s.Db.Update(initStatusBuckets)
	if err != nil {
		return err
	}
	err = s.Db.Update(initConfigBucket)
	if err != nil {
		return err
	}
	return nil
}

// Close closes the connection to the BoltDB database.
//
// Returns:
//   - error: An error if closing the database fails.
func (s *Storage) Close() error {
	return s.Db.Close()
}
