package repository

import (
	"context"
	"fmt"
	"go-ms-template-service/model/domain"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CinemaRepositoryImpl struct {
	DB *gorm.DB
}

func NewCinemaRepository(db *gorm.DB) *CinemaRepositoryImpl {
	return &CinemaRepositoryImpl{DB: db}
}

func (repository *CinemaRepositoryImpl) Save(ctx context.Context, tx *gorm.DB, cinema domain.Cinema) (domain.Cinema, error) {
	result := repository.DB.Model(&cinema).Omit("deleted_at", "updated_at").Create(&cinema).Error

	if result != nil {
		// Cek jika kesalahan adalah kesalahan MySQL
		if mysqlErr, ok := result.(*mysql.MySQLError); ok {
			return cinema, mysqlErr
		}
		return cinema, result
	}

	return cinema, nil
}

func (repository *CinemaRepositoryImpl) Update(ctx context.Context, tx *gorm.DB, cinema domain.Cinema, CinemaCode string) domain.Cinema {
	tx.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&cinema, "cinema_code = ?", CinemaCode).Error
		if err != nil {
			return err
		}

		err = tx.Model(&domain.Cinema{}).Where("cinema_code = ?", CinemaCode).Updates(map[string]interface{}{
			"cinema_owner":      cinema.CinemaOwner,
			"location_code":     cinema.LocationCode,
			"province_code":     cinema.ProvinceCode,
			"city_code":         cinema.CityCode,
			"region_code":       cinema.RegionCode,
			"company_code":      cinema.CompanyCode,
			"cinema_level":      cinema.CinemaLevel,
			"oracle_code":       cinema.OracleCode,
			"is_data_migration": cinema.IsDataMigration,
			"close_flag":        cinema.CloseFlag,
			"close_start":       cinema.CloseStart,
			"close_end":         cinema.CloseEnd,
			"operator_email":    cinema.OperatorEmail,
			"updated_by":        cinema.UpdatedBy,
			"updated_host_ip":   cinema.UpdatedHostIp,
		}).Error

		return err
	})

	return cinema
}

func (repository *CinemaRepositoryImpl) Delete(ctx context.Context, tx *gorm.DB, cinema domain.Cinema) {
	dt := domain.Cinema{}
	result := repository.DB.Take(&dt, "cinema_code = ?", cinema.CinemaCode)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Print("Tidak ada record yang bisa didelete")
		}
	}

	repository.DB.Where("cinema_code = ?", cinema.CinemaCode).Update("deleted_host_ip", cinema.DeletedHostIp)
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
