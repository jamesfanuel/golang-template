package web

type CinemaCreateRequest struct {
	CinemaCode   string `validate:"required" json:"cinema_code"`
	CinemaName   string `validate:"required" json:"cinema_name"`
	ProvinceCode string `validate:"required" json:"province_code"`
	CityCode     string `validate:"required" json:"city_code"`
	RegionCode   string `validate:"required" json:"region_code"`
	CinemaLevel  string `validate:"required" json:"cinema_level"`
}
