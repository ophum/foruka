package datastore

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DataStore struct {
	store *sql.DB
}

func NewDataStore(path string) *DataStore {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/database.sqlite", path))
	if err != nil {
		return nil
	}

	_, err = db.Exec(
		`create table if not exists "kvs" ("key", "value")`,
	)

	ds := &DataStore{
		store: db,
	}
	return ds
}

func (d *DataStore) get(key string) (string, error) {
	row := d.store.QueryRow(
		`select value from kvs where key=?`,
		key,
	)
	var value string
	err := row.Scan(&value)

	if err != nil {
		return "", err
	}
	return value, nil
}

func (d *DataStore) Get(key string) string {
	v, _ := d.get(key)
	return v
}

func (d *DataStore) Del(key string) {
	d.store.Exec(
		`delete from kvs where key=?`,
		key,
	)
}

func (d *DataStore) Set(key, value string) {
	var s string
	if _, err := d.get(key); err == sql.ErrNoRows {
		s = `insert into kvs(key, value) values(?, ?)`
		_, _ = d.store.Exec(
			s,
			key,
			value,
		)
	} else {
		s = `update kvs set value=? where key=?`
		_, _ = d.store.Exec(
			s,
			value,
			key,
		)
	}
}

func (d *DataStore) Gets(search string) map[string]string {
	rows, _ := d.store.Query(
		`select key, value from kvs where key like ?`,
		search,
	)

	defer rows.Close()

	res := map[string]string{}
	for rows.Next() {
		var key string
		var value string

		if err := rows.Scan(&key, &value); err != nil {
			return nil
		}
		res[key] = value
	}

	return res
}
