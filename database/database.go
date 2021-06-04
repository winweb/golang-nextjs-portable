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

func DbOpen(filename string) (*gorm.DB, error) {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             4 * time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Silent,    // Log level
			IgnoreRecordNotFoundError: true,             // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,            // Disable color
		},
	)

	db, err := gorm.Open(sqlite.Open(filename+"?_journal_mode=WAL&_synchronous=NORMAL&mode=rwc&cache=shared&_busy_timeout=40000"), &gorm.Config{
		//PrepareStmt: true,
		SkipDefaultTransaction: true,
		CreateBatchSize: 10,
		Logger: newLogger,
		//Logger: logger.Discard, //disable logger print Slow SQL and happening errors
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("Opening db file: ", filename)

	sdb, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	err = sdb.Ping()
	if err != nil {
		log.Fatal(err)
	}

	sdb.SetMaxIdleConns(1)
	sdb.SetMaxOpenConns(1)
	//sdb.SetConnMaxLifetime(time.Nanosecond)

	sqlite3Config(sdb)

	db.AutoMigrate(&models.People{})

	sdb.Exec("CREATE UNIQUE INDEX IF NOT EXISTS people_id_index ON peoples (id);")

	return db, nil
}

func sqlite3Config(db *sql.DB) {

	db.Exec("PRAGMA page_size= 65535;")
	db.Exec("PRAGMA cache_size= 8000;")
	db.Exec("PRAGMA mmap_size = 30000000000;")

	/*
	db.Exec("PRAGMA journal_mode = WAL;")
	db.Exec("PRAGMA locking_mode = EXCLUSIVE;")
	db.Exec("PRAGMA busy_timeout = 20000;")
	//db.Exec("PRAGMA cache = shared;")

	//db.Exec("PRAGMA cache_spill = ON;")
	//db.Exec("PRAGMA count_changes = OFF;")
	db.Exec("PRAGMA encoding = \"UTF-8\";")
	//db.Exec("PRAGMA journal_mode = delete;")
	//db.Exec("PRAGMA locking_mode = EXCLUSIVE;")
	//db.Exec("PRAGMA main.synchronous=NORMAL;")
	//db.Exec("PRAGMA page_size = 4096;")
	//db.Exec("PRAGMA shrink_memory;")
	db.Exec("PRAGMA synchronous = NORMAL;")
	db.Exec("PRAGMA temp_store = MEMORY;")
	//db.Exec("legacy_file_format=OFF;")

	//db.Exec("PRAGMA mmap_size=1099511627776;")
	//db.Exec("PRAGMA threads = 5;")
	//db.Exec("PRAGMA wal_autocheckpoint = 1638400;")
	*/
}