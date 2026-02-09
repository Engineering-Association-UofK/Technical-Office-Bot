package repository

import (
	"fmt"

	"github.com/abdulrahim-m/Technical-Office-Bot/internal/models"
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

func (br *BaseRepo[T]) FineById(id int) (T, error) {
	var entry T
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", br.TableName)
	err := br.DB.Get(&entry, query, id)
	return entry, err
}

type AdminRepository struct {
	BaseRepo[models.FeedbackModel]
}

func NewAdminRepository(db *sqlx.DB) *AdminRepository {
	return &AdminRepository{
		BaseRepo: BaseRepo[models.FeedbackModel]{
			DB:        db,
			TableName: "admins",
		},
	}
}

// You can still add specific methods that only Admins have!
func (r *AdminRepository) FindByEmail(email string) (*models.FeedbackModel, error) {
	var a models.FeedbackModel
	err := r.DB.Get(&a, "SELECT * FROM admins WHERE email = ?", email)
	return &a, err
}
