package models

import "time"

type DefaultModel struct {
	ID        int       `gorm:"column:id;type:integer;primary_key;unique;AUTO_INCREMENT"`
	CreatedAt time.Time `gorm:"column:createdAt;type:DATETIME;not null"`
	UpdatedAt time.Time `gorm:"column:updatedAt;type:DATETIME;not null"`
}
