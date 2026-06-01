package migrate

import (
	"fmt"

	"github.com/fernoe1/appointment-telegram-bot/internal/domain"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	if err := db.AutoMigrate(&domain.Appointment{}); err != nil {
		return fmt.Errorf("migrate.Run->AutoMigrate: %w", err)
	}

	return nil
}
