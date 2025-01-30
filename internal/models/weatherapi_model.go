package models

import "fmt"

type DBContractWeatherapi interface {
	CheckUserAPIKey(apiKey string) (bool, error)
}

type WeatherapiModel struct {
	db DBContractWeatherapi
}

func NewWeatherapiModel(db DBContractWeatherapi) *WeatherapiModel {
	return &WeatherapiModel{db: db}
}

func (msql *MySQL) CheckUserAPIKey(apiKey string) (bool, error) {
	stmt := `SELECT COUNT(*) FROM api_keys WHERE api_key=?`

	var count int

	err := msql.DB.QueryRow(stmt, apiKey).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to scan count of api key in the database: %w", err)
	}

	if count > 0 {
		return true, nil
	}
	return false, ErrAPIKeyNotFound
}
