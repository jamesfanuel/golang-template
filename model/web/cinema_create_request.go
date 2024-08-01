package web

type CinemaCreateRequest struct {
	CinemaCode      string `validate:"required" json:"cinema_code"`
	CinemaName      string `validate:"required" json:"cinema_name"`
	CinemaOwner     string `json:"cinema_owner"`
	LocationCode    string `json:"location_code"`
	ProvinceCode    string `validate:"required" json:"province_code"`
	CityCode        string `validate:"required" json:"city_code"`
	RegionCode      string `validate:"required" json:"region_code"`
	CompanyCode     string `json:"company_code"`
	CinemaLevel     string `validate:"required" json:"cinema_level"`
	OracleCode      string `json:"oracle_code"`
	IsDataMigration string `json:"is_data_migration"`
	CloseFlag       string `json:"close_flag"`
	CloseStart      string `json:"close_start"`
	CloseEnd        string `json:"close_end"`
	OperatorEmail   string `validate:"email" json:"operator_email"`
	CreatedBy       string `json:"created_by"`
	CreatedHostIp   string `json:"created_host_ip"`
}
