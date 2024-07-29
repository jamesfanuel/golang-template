package repository

import (
	"context"
	"x1-cinema/model/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CinemaRepositoryImpl struct {
	DB *gorm.DB
}

func NewCinemaRepository(db *gorm.DB) *CinemaRepositoryImpl {
	return &CinemaRepositoryImpl{DB: db}
}

func (repository *CinemaRepositoryImpl) Save(ctx context.Context, tx *gorm.DB, cinema domain.Cinema) domain.Cinema {
	repository.DB.Transaction(func(tx *gorm.DB) error {
		var cinema domain.Cinema
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&cinema, "cinema_code = ?", cinema.CinemaCode).Error
		if err != nil {
			return err
		}

		dt := domain.Cinema{
			CinemaCode:   cinema.CinemaCode,
			CinemaName:   cinema.CinemaName,
			ProvinceCode: cinema.ProvinceCode,
			CityCode:     cinema.CityCode,
			RegionCode:   cinema.RegionCode,
			CinemaLevel:  cinema.CinemaLevel,
			// RowId:        id,
		}

		return repository.DB.Save(&dt).Error
	})

	return cinema
}

// func (repository *CinemaRepositoryImpl) Update(ctx context.Context, tx *gorm.DB, cinema domain.Cinema, CinemaCode string) domain.Cinema {
// 	SQL := "update mg_cinema set cinema_name = ? where cinema_code = ?"
// 	_, err := tx.ExecContext(ctx, SQL, cinema.CinemaName, CinemaCode)

// 	helper.PanicIfError(err)

// 	return cinema
// }

func (repository *CinemaRepositoryImpl) Delete(ctx context.Context, tx *gorm.DB, cinema domain.Cinema) {
	repository.DB.Delete(&cinema)
}

func (repository *CinemaRepositoryImpl) FindByCode(ctx context.Context, tx *gorm.DB, CinemaCode string) (domain.Cinema, error) {
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
	dt := []domain.Cinema{}

	result := repository.DB.Find(&dt)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return dt, result.Error
		}
	}

	return dt, nil
}
