package data

import (
	"database/sql"
	"errors"

	"gorm.io/gorm"
)

func IsRowNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, sql.ErrNoRows)
}
