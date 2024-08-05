package domain

import (
	"time"
)

type Cinema struct {
	CinemaCode      string `gorm:"primary_key;<-:create"`
	CinemaName      string `gorm:"<-:create"`
	CinemaOwner     string
	LocationCode    string
	ProvinceCode    string
	CityCode        string
	RegionCode      string
	CompanyCode     string
	CinemaLevel     string
	OracleCode      string
	IsDataMigration string
	CloseFlag       string
	// RowId        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	CloseStart    string
	CloseEnd      string
	OperatorEmail string `gorm:"unique"`
	CreatedBy     string
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime;<-:create"`
	CreatedHostIp string
	UpdatedBy     string
	UpdatedAt     time.Time
	UpdatedHostIp string
	DeletedBy     string
	DeletedAt     time.Time
	DeletedHostIp string
}

func (c *Cinema) TableName() string {
	return "mg_cinema"
}
