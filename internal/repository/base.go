package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type BaseRepo[T any] struct {
	DB        *sqlx.DB
	TableName string
}

func (br *BaseRepo[T]) GetAll() ([]T, error) {
	var entries []T
	query := fmt.Sprintf("SELECT * FROM %s", br.TableName)
	err := br.DB.Select(&entries, query)
	return entries, err
}

func (br *BaseRepo[T]) FindById(id int64) (T, error) {
	var entry T
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", br.TableName)
	err := br.DB.Get(&entry, query, id)
	return entry, err
}

func (br *BaseRepo[T]) Persist(query string, args ...any) (int64, error) {
	result, err := br.DB.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
