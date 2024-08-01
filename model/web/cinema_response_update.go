package web

import "time"

type CinemaResponseUpdate struct {
	CinemaCode      string    `json:"cinema_code"`
	CinemaName      string    `json:"cinema_name"`
	CinemaOwner     string    `json:"cinema_owner"`
	LocationCode    string    `json:"location_code"`
	ProvinceCode    string    `json:"province_code"`
	CityCode        string    `json:"city_code"`
	RegionCode      string    `json:"region_code"`
	CompanyCode     string    `json:"company_code"`
	CinemaLevel     string    `json:"cinema_level"`
	OracleCode      string    `json:"oracle_code"`
	IsDataMigration string    `json:"is_data_migration"`
	CloseFlag       string    `json:"close_flag"`
	CloseStart      string    `json:"close_start"`
	CloseEnd        string    `json:"close_end"`
	OperatorEmail   string    `json:"operator_email"`
	UpdatedBy       string    `json:"updated_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedHostIp   string    `json:"updated_host_ip"`
}
