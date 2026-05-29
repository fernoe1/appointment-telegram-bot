package pkg

import (
	"fmt"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewGormDB(url string) (*gorm.DB, error) {
	db, err := gorm.Open(
		sqlite.Dialector{
			DriverName: "libsql",
			DSN:        url,
		},
		&gorm.Config{},
	)
	if err != nil {
		return nil, fmt.Errorf("pkg.NewGormDB->gorm.Open: %w", err)
	}

	return db, nil
}
