package defaultdb

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type SqlPersistence struct {
	db *sql.DB
}

func (dmp SqlPersistence) Find(key string) (string, bool) {
	row, err := dmp.db.Query("SELECT val FROM node WHERE id=?", key)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var val string
		row.Scan(&val)
		return val, true
	}
	return "", false
}

func (dmp SqlPersistence) Add(key, value string) error {
	_, err := dmp.db.Exec("INSERT OR REPLACE INTO node (id, val)  VALUES (?, ?);", key, value)
	return err
}

func NewSqlDB(filename string) SqlPersistence {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		file.Close()
	}

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS node ("id" TEXT NOT NULL PRIMARY KEY,"val" TEXT NOT NULL);CREATE UNIQUE INDEX IF NOT EXISTS id_idx ON node(id);`) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	return SqlPersistence{db: db}
}
