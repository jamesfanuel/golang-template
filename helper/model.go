package helper

import (
	"x1-cinema/model/domain"
	"x1-cinema/model/web"
)

func ToCinemaResponse(cinema domain.Cinema) web.CinemaResponseFind {
	return web.CinemaResponseFind{
		CinemaCode:      cinema.CinemaCode,
		CinemaName:      cinema.CinemaName,
		CinemaOwner:     cinema.CinemaOwner,
		LocationCode:    cinema.LocationCode,
		ProvinceCode:    cinema.ProvinceCode,
		CityCode:        cinema.CityCode,
		RegionCode:      cinema.RegionCode,
		CompanyCode:     cinema.CompanyCode,
		CinemaLevel:     cinema.CinemaLevel,
		OracleCode:      cinema.OracleCode,
		IsDataMigration: cinema.IsDataMigration,
		CloseFlag:       cinema.CloseFlag,
		CloseStart:      cinema.CloseStart,
		CloseEnd:        cinema.CloseEnd,
		OperatorEmail:   cinema.OperatorEmail,
		CreatedBy:       cinema.CreatedBy,
		CreatedAt:       cinema.CreatedAt,
		CreatedHostIp:   cinema.CreatedHostIp,
		UpdatedBy:       cinema.UpdatedBy,
		UpdatedAt:       cinema.UpdatedAt,
		UpdatedHostIp:   cinema.UpdatedHostIp,
	}
}

func ToCinemaCreateResponse(cinema domain.Cinema) web.CinemaCreateResponse {
	return web.CinemaCreateResponse{
		CinemaCode:      cinema.CinemaCode,
		CinemaName:      cinema.CinemaName,
		CinemaOwner:     cinema.CinemaOwner,
		LocationCode:    cinema.LocationCode,
		ProvinceCode:    cinema.ProvinceCode,
		CityCode:        cinema.CityCode,
		RegionCode:      cinema.RegionCode,
		CompanyCode:     cinema.CompanyCode,
		CinemaLevel:     cinema.CinemaLevel,
		OracleCode:      cinema.OracleCode,
		IsDataMigration: cinema.IsDataMigration,
		CloseFlag:       cinema.CloseFlag,
		CloseStart:      cinema.CloseStart,
		CloseEnd:        cinema.CloseEnd,
		OperatorEmail:   cinema.OperatorEmail,
		CreatedBy:       cinema.CreatedBy,
		CreatedAt:       cinema.CreatedAt,
		CreatedHostIp:   cinema.CreatedHostIp,
	}
}

func ToCinemaUpdateResponse(cinema domain.Cinema) web.CinemaUpdateResponse {
	return web.CinemaUpdateResponse{
		CinemaCode:      cinema.CinemaCode,
		CinemaName:      cinema.CinemaName,
		CinemaOwner:     cinema.CinemaOwner,
		LocationCode:    cinema.LocationCode,
		ProvinceCode:    cinema.ProvinceCode,
		CityCode:        cinema.CityCode,
		RegionCode:      cinema.RegionCode,
		CompanyCode:     cinema.CompanyCode,
		CinemaLevel:     cinema.CinemaLevel,
		OracleCode:      cinema.OracleCode,
		IsDataMigration: cinema.IsDataMigration,
		CloseFlag:       cinema.CloseFlag,
		CloseStart:      cinema.CloseStart,
		CloseEnd:        cinema.CloseEnd,
		OperatorEmail:   cinema.OperatorEmail,
		UpdatedBy:       cinema.UpdatedBy,
		UpdatedAt:       cinema.UpdatedAt,
		UpdatedHostIp:   cinema.UpdatedHostIp,
	}
}
