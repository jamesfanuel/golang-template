package web

type CinemaResponse struct {
	CinemaCode   string `json:"cinema_code"`
	CinemaName   string `json:"cinema_name"`
	ProvinceCode string `json:"province_code"`
	CityCode     string `json:"city_code"`
	RegionCode   string `json:"region_code"`
	CinemaLevel  string `json:"cinema_level"`
}
