package web

type CinemaCreateRequest struct {
	CinemaCode string `validate:"required" json:"cinema_code"`
	CinemaName string `validate:"required" json:"cinema_name"`
}
