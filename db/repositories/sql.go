package repositories

import (
	"context"

	"github.com/HollWill/weather_telegram_bot/db/models"
	"github.com/HollWill/weather_telegram_bot/db/sqlstore"
	"github.com/jmoiron/sqlx"
)

type sqlRepository interface {
	getDB() *sqlx.DB
}

type SqlUserRepository struct {
	db *sqlx.DB
}

func NewSqlUserRepository(db *sqlx.DB) *SqlUserRepository {
	return &SqlUserRepository{db: db}
}

func (r *SqlUserRepository) getDB() *sqlx.DB {
	return r.db
}

func (r *SqlUserRepository) Save(ctx context.Context, u *models.User) error {
	db := r.getDB()
	return sqlstore.SaveUser(ctx, db, u)
}

func (r *SqlUserRepository) FindByID(ctx context.Context, id int) (*models.User, error) {
	db := r.getDB()
	return sqlstore.FindUserByID(ctx, db, id)
}

func (r *SqlUserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	db := r.getDB()
	return sqlstore.GetAllUsers(ctx, db)
}
