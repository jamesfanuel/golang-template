package domain

import (
	"time"

	"gorm.io/gorm"
)

// type Cinema struct {
// 	CinemaCode string `json:"cinema_code"`
// 	CinemaName string `json:"cinema_name"`
// }

type Cinema struct {
	CinemaCode   string `gorm:"primary_key"`
	CinemaName   string
	ProvinceCode string
	CityCode     string
	RegionCode   string
	CinemaLevel  string
	// RowId        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (c *Cinema) TableName() string {
	return "mg_cinema"
}
