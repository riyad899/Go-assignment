package models

import "time"

type Reservation struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	UserID       uint        `gorm:"not null" json:"user_id"`
	ZoneID       uint        `gorm:"not null" json:"zone_id"`
	LicensePlate string      `gorm:"not null;type:varchar(15)" json:"license_plate"`
	Status       string      `gorm:"type:varchar(20);default:'active';not null" json:"status"` // active, completed, cancelled
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	
	// Associations
	User         *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Zone         *ParkingZone `gorm:"foreignKey:ZoneID" json:"zone,omitempty"`
}
