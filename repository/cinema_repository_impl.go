package repository

import (
	"context"
	"x1-cinema/model/domain"

	"gorm.io/gorm"
)

type CinemaRepositoryImpl struct {
	DB *gorm.DB
}

func NewCinemaRepository(db *gorm.DB) *CinemaRepositoryImpl {
	return &CinemaRepositoryImpl{DB: db}
}

func (repository *CinemaRepositoryImpl) Save(ctx context.Context, tx *gorm.DB, cinema domain.Cinema) domain.Cinema {
	// SQL := "insert into mg_cinema (cinema_code, cinema_name) values (?, ?)"
	// _, err := tx.ExecContext(ctx, SQL, cinema.CinemaCode, cinema.CinemaName)

	// id := uuid.New()
	dt := domain.Cinema{
		CinemaCode:   cinema.CinemaCode,
		CinemaName:   cinema.CinemaName,
		ProvinceCode: cinema.ProvinceCode,
		CityCode:     cinema.CityCode,
		RegionCode:   cinema.RegionCode,
		CinemaLevel:  cinema.CinemaLevel,
		// RowId:        id,
	}

	repository.DB.Save(&dt)

	return cinema
}

// func (repository *CinemaRepositoryImpl) Update(ctx context.Context, tx *gorm.DB, cinema domain.Cinema, CinemaCode string) domain.Cinema {
// 	SQL := "update mg_cinema set cinema_name = ? where cinema_code = ?"
// 	_, err := tx.ExecContext(ctx, SQL, cinema.CinemaName, CinemaCode)

// 	helper.PanicIfError(err)

// 	return cinema
// }

func (repository *CinemaRepositoryImpl) Delete(ctx context.Context, tx *gorm.DB, cinema domain.Cinema) {
	// SQL := "update mg_cinema set is_active = 0 where cinema_code = ?"
	// _, err := tx.ExecContext(ctx, SQL, CinemaCode)
	// helper.PanicIfError(err)
	repository.DB.Delete(&cinema)
}

func (repository *CinemaRepositoryImpl) FindByCode(ctx context.Context, tx *gorm.DB, CinemaCode string) (domain.Cinema, error) {
	// SQL := "select cinema_code, cinema_name from mg_cinema where cinema_code = ? and is_active = 1"
	// rows, err := tx.QueryContext(ctx, SQL, CinemaCode)
	// helper.PanicIfError(err)
	// defer rows.Close()

	// cinema := domain.Cinema{}

	// if rows.Next() {
	// 	rows.Scan(&cinema.CinemaCode, &cinema.CinemaName)
	// 	helper.PanicIfError(err)

	// 	return cinema, nil
	// } else {
	// 	return cinema, errors.New("cinema is not found")
	// }
	dt := domain.Cinema{}

	result := repository.DB.Take(&dt, "cinema_code = ?", CinemaCode)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return dt, result.Error
		}
	}

	return dt, nil

}

func (repository *CinemaRepositoryImpl) FindAll(ctx context.Context, tx *gorm.DB) ([]domain.Cinema, error) {
	// SQL := "select cinema_code, cinema_name FROM mg_cinema WHERE is_active = 1"
	// rows, err := tx.QueryContext(ctx, SQL)
	// helper.PanicIfError(err)
	// defer rows.Close()

	// var cinemas []domain.Cinema
	// for rows.Next() {
	// 	cinema := domain.Cinema{}
	// 	err := rows.Scan(&cinema.CinemaCode, &cinema.CinemaName)
	// 	helper.PanicIfError(err)
	// 	cinemas = append(cinemas, cinema)
	// }

	// return cinemas

	dt := []domain.Cinema{}

	result := repository.DB.Find(&dt)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return dt, result.Error
		}
	}

	return dt, nil
}
