package domain

import (
	"database/sql"
	"time"
)

type Base struct {
	ID         int
	CreateDate time.Time
	UpdateDate time.Time
	DeleteDate sql.NullTime
}
