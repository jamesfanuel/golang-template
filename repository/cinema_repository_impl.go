package repository

import (
	"context"
	"database/sql"
	"errors"
	"x1-cinema/helper"
	"x1-cinema/model/domain"
)

type CinemaRepositoryImpl struct {
}

func NewCinemaRepository() CinemaRepository {
	return &CinemaRepositoryImpl{}
}

func (repository *CinemaRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, cinema domain.Cinema) domain.Cinema {
	SQL := "insert into mg_cinema (cinema_code, cinema_name) values (?, ?)"
	_, err := tx.ExecContext(ctx, SQL, cinema.CinemaCode, cinema.CinemaName)

	helper.PanicIfError(err)

	return cinema
}

func (repository *CinemaRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, cinema domain.Cinema) domain.Cinema {
	SQL := "update mg_cinema set cinema_name = ? where cinema_code = ?"
	_, err := tx.ExecContext(ctx, SQL, cinema.CinemaName, cinema.CinemaCode)

	helper.PanicIfError(err)

	return cinema
}

func (repository *CinemaRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, cinema domain.Cinema) {
	SQL := "delete from mg_cinema where cinema_code = ?"
	_, err := tx.ExecContext(ctx, SQL, cinema.CinemaCode)
	helper.PanicIfError(err)
}

func (repository *CinemaRepositoryImpl) FindByCode(ctx context.Context, tx *sql.Tx, CinemaCode string) (domain.Cinema, error) {
	SQL := "select cinema_code, cinema_name from mg_cinema where cinema_code = ?"
	rows, err := tx.QueryContext(ctx, SQL, CinemaCode)
	helper.PanicIfError(err)
	defer rows.Close()

	cinema := domain.Cinema{}

	if rows.Next() {
		rows.Scan(&cinema.CinemaCode, &cinema.CinemaName)
		helper.PanicIfError(err)

		return cinema, nil
	} else {
		return cinema, errors.New("Cinema is not found")
	}
}

func (repository *CinemaRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.Cinema {
	SQL := "select cinema_code, cinema_name FROM mg_cinema"
	rows, err := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(err)
	defer rows.Close()

	var cinemas []domain.Cinema
	for rows.Next() {
		cinema := domain.Cinema{}
		err := rows.Scan(&cinema.CinemaCode, &cinema.CinemaName)
		helper.PanicIfError(err)
		cinemas = append(cinemas, cinema)
	}

	return cinemas
}
