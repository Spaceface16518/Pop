package store

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	// "github.com/lib/pq"
	"log"
)

const (
	TableName       = "suggestions"
	NameColumnName  = "name"
	CountColumnName = "count"
)

type Store interface {
	Save(mem map[string]int) error
	Load() (map[string]int, error)
	Close() error
}

var psqlStatementBuilder sq.StatementBuilderType = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// DataStore is the top/package level variable that is essentially an interface with the application database.
var DataStore *DataStoreWrapper

type DataStoreWrapper struct {
	db *sql.DB
}

// NewDataStore wraps a database pointer in a Store interface compatible structure
func NewDataStore(dbPointer *sql.DB) *DataStoreWrapper {
	return &DataStoreWrapper{dbPointer}
}

func (dataStore *DataStoreWrapper) Save(memRepr map[string]int) error {
	tx, err := dataStore.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if err = saveOnTx(memRepr, tx); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (dataStore *DataStoreWrapper) Load() (map[string]int, error) {
	tx, err := dataStore.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var suggestions map[string]int
	suggestions, err = loadFromTx(tx)
	return suggestions, err
}

func (dataStore *DataStoreWrapper) Close() error {
	return dataStore.db.Close()
}

// InitTable initializes the "suggestions" table using "CREATE TABLE IF NOT EXISTS". This function should be run before executing other actions against the database.
func (dataStore *DataStoreWrapper) InitTable() error {
	_, err := dataStore.db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(%s TEXT NOT NULL UNIQUE, %v INTEGER NOT NULL)", TableName, NameColumnName, CountColumnName))
	return err
}

// SetStore intializes the package level variable DataStore to a specified value. This is for use in test mocking, and can go in the init() of the main package.
func SetStore(s *DataStoreWrapper) {
	DataStore = s
}

// func SetStore(s Store) {
// 	DataStore = s
// }

func saveOnTx(memRepr map[string]int, tx *sql.Tx) error {
	s := psqlStatementBuilder.Insert(TableName).Columns(NameColumnName, CountColumnName)
	for name, count := range memRepr {
		s = s.Values(name, count)
	}
	queryString, queryArgs, err := s.ToSql()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(queryString, queryArgs...); err != nil {
		return err
	}
	return nil
}

func loadFromTx(tx *sql.Tx) (map[string]int, error) {
	s := psqlStatementBuilder.Select(NameColumnName, CountColumnName).From(TableName)
	queryString, _, err := s.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(queryString)
	if err != nil {
		return nil, err
	}

	suggestions := map[string]int{}
	for rows.Next() {
		var name string
		var count int

		if err := rows.Scan(&name, &count); err != nil {
			return nil, err
		}

		suggestions[name] += count
	}
	return suggestions, nil
}
