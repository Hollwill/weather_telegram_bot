package sqlstore

import (
	"context"
	"database/sql"

	"github.com/HollWill/weather_telegram_bot/db/models"

	"github.com/jmoiron/sqlx"
)

const UserTable = "users"

func CreateTables(db *sqlx.DB) {
	db.Exec(`CREATE TABLE IF NOT EXISTS ` + UserTable + ` (
    id int not null primary key,
    name varchar(250) not null,
    lat real,
    long real,
    crontab varchar(250) not null
  )`)
}

type SqlxDatabase interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

func FindUserByID(ctx context.Context, db SqlxDatabase, id int) (*models.User, error) {
	p := new(models.User)
	sql := `SELECT * FROM ` + UserTable + ` WHERE id=$1`
	err := db.GetContext(ctx, p, sql, id)
	return p, err
}

func SaveUser(ctx context.Context, db SqlxDatabase, u *models.User) error {
	sql := `INSERT INTO ` + UserTable + `(id, name, lat, long, crontab) VALUES ($1, $2, $3, $4, $5) ON CONFLICT(id) DO UPDATE SET name=EXCLUDED.name, lat=EXCLUDED.lat, long=EXCLUDED.long, crontab=EXCLUDED.crontab RETURNING id`
	var lastId int
	stmt, err := db.PreparexContext(ctx, sql)
	if err != nil {
		return err
	}
	err = stmt.GetContext(ctx, &lastId, u.ID, u.Name, u.Latitude, u.Longitude, u.CronTab)

	return err
}

func GetAllUsers(ctx context.Context, db SqlxDatabase) ([]models.User, error) {
	var users []models.User
	sql := `SELECT * FROM ` + UserTable
	err := db.SelectContext(ctx, &users, sql)

	return users, err
}
