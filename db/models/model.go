package models

type User struct {
	ID        int     `db:"id"      bson:"id"`
	Name      string  `db:"name"    bson:"name"`
	Latitude  float32 `db:"lat"     bson:"lat"`
	Longitude float32 `db:"long"    bson:"long"`
	CronTab   string  `db:"crontab" bson:"crontab"`
}
