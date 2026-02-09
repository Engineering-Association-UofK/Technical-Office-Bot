package models

import "time"

type AdminOTP struct {
	ID        int       `db:"id"`
	AdminID   int       `db:"admin_id"`
	Code      string    `db:"code"`
	CreatedAt time.Time `db:"created_at"`
}
