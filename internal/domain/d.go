package domain

const (
	MaxAppointmentsPerDay = 2
	AppointmentDateLayout = "2006-01-02"
)

type User struct {
	ID          uint         `gorm:"primaryKey"`
	TelegramID  int64        `gorm:"uniqueIndex;not null"`
	Username    string       `gorm:"size:255"`
	PhoneNumber string       `gorm:"size:64;not null"`
	Appointment *Appointment `gorm:"constraint:OnDelete:CASCADE;"`
}

type BookedDay struct {
	ID           uint          `gorm:"primaryKey"`
	Date         string        `gorm:"size:10;uniqueIndex;not null"`
	Full         bool          `gorm:"default:false"`
	Appointments []Appointment `gorm:"constraint:OnDelete:CASCADE;"`
}

type Appointment struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint `gorm:"uniqueIndex;not null"`
	User        User
	BookedDayID uint `gorm:"index;not null"`
	BookedDay   BookedDay
	Hour        int `gorm:"not null"`
}
