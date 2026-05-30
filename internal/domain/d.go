package domain

const (
	MaxAppointmentsPerDay = 2
	AppointmentDateLayout = "2006-01-02"
)

type Appointment struct {
	TID         int64  `gorm:"primaryKey"`
	CID         int64  `gorm:"uniqueIndex;not null"`
	Username    string `gorm:"size:255"`
	PhoneNumber string `gorm:"size:64;not null"`

	Date string `gorm:"size:10;index;not null"`
	Hour int    `gorm:"not null"`
}
