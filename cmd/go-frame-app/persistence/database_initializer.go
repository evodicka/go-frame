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
	// Db is the global BoltDB database instance.
	Db *bolt.DB
)

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lmsgprefix|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lmsgprefix|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lmsgprefix|log.Lshortfile)

	InfoLogger.Println("Opening Database Connection")
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		ErrorLogger.Fatal(err)
	} else {
		Db = db
		InfoLogger.Println("Database connection established")
		initBuckets()
	}
}

func initBuckets() {
	err := Db.Update(initImageBuckets)
	if err != nil {
		ErrorLogger.Fatal(err)
	}
	err = Db.Update(initStatusBuckets)
	if err != nil {
		ErrorLogger.Fatal(err)
	}
	err = Db.Update(initConfigBucket)
	if err != nil {
		ErrorLogger.Fatal(err)
	}
}

// Close closes the connection to the BoltDB database.
//
// Returns:
//   - error: An error if closing the database fails.
func Close() error {
	return Db.Close()
}
