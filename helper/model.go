package helper

import (
	"x1-cinema/model/domain"
	"x1-cinema/model/web"
)

func ToCinemaResponse(cinema domain.Cinema) web.CinemaResponse {
	return web.CinemaResponse{
		CinemaCode:   cinema.CinemaCode,
		CinemaName:   cinema.CinemaName,
		ProvinceCode: cinema.ProvinceCode,
		CityCode:     cinema.CityCode,
		RegionCode:   cinema.RegionCode,
		CinemaLevel:  cinema.CinemaLevel,
	}
}
