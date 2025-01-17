package database

import (
	"database/sql"
	"gorm.io/gorm/logger"
	"os"
	"time"

	"github.com/dstotijn/golang-nextjs-portable/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func DbOpen(filename string) (*gorm.DB, *sql.DB, error) {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             4 * time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error,    // Log level
			IgnoreRecordNotFoundError: true,             // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,            // Disable color
		},
	)

	db, err := gorm.Open(sqlite.Open(filename+"?_journal_mode=WAL&_txlock=exclusive&_synchronous=NORMAL&mode=rwc&cache=shared&_cache_size=10000&_busy_timeout=40000&_loc=UTC"), &gorm.Config{
		PrepareStmt: true,
		SkipDefaultTransaction: true,
		CreateBatchSize: 10,
		Logger: newLogger,
	})
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	log.Println("Opening db file: ", filename)

	sdb, err := db.DB()
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	err = sdb.Ping()
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	sdb.SetMaxIdleConns(1)
	sdb.SetMaxOpenConns(1)

	_, err = sdb.Exec("PRAGMA page_size= 65535;")
	_, err = sdb.Exec("PRAGMA mmap_size = 30000000000;")

	err = db.AutoMigrate(&models.People{}, &models.CreditCard{})

	return db, sdb, nil
}