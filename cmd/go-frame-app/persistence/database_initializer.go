package persistence

import (
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	Db            *bolt.DB
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
}

func Close() error {
	return Db.Close()
}
