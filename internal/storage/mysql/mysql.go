package mysql

import (
	"database/sql"
	"exchange-rates/internal/config"
	"exchange-rates/internal/exrate"
	"exchange-rates/internal/storage"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	db *sql.DB
}

func New(cnf *config.Config) (*Storage, error) {
	const op = "storage.mysql.New"

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cnf.Db.User, cnf.Db.Pwd, cnf.Db.Host, cnf.Db.Port, cnf.Db.Name)

	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GetAllCurrencyRates() ([]exrate.Options, error) {
	const op = "storage.mysql.GetAllCurrencyRates"

	stmt, err := s.db.Prepare("" +
		"SELECT cr.currency_id, cr.date, c.code, c.name, c.scale, cr.rate " +
		"FROM currency_rate AS cr " +
		"INNER JOIN currency AS c ON cr.currency_id = c.currency_id " +
		"ORDER BY cr.date;")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	var cRates []exrate.Options

	for rows.Next() {
		var cRate exrate.Options

		sErr := rows.Scan(&cRate.CurID, &cRate.Date, &cRate.CurCode, &cRate.CurName, &cRate.CurScale, &cRate.CurRate)
		if sErr != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		cRates = append(cRates, cRate)
	}

	if cRates == nil {
		return nil, storage.ErrCurRatesNotFound
	}

	return cRates, nil
}

func (s *Storage) GetCurrencyRatesByDate(date time.Time) ([]exrate.Options, error) {
	const op = "storage.mysql.GetCurrencyRatesByDate"

	stmt, err := s.db.Prepare("" +
		"SELECT cr.currency_id, cr.date, c.code, c.name, c.scale, cr.rate " +
		"FROM currency_rate AS cr " +
		"INNER JOIN currency AS c ON cr.currency_id = c.currency_id " +
		"WHERE cr.date = ?;")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	defer stmt.Close()

	rows, err := stmt.Query(date.Format(time.DateOnly))
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	var cRates []exrate.Options

	for rows.Next() {
		var cRate exrate.Options

		sErr := rows.Scan(&cRate.CurID, &cRate.Date, &cRate.CurCode, &cRate.CurName, &cRate.CurScale, &cRate.CurRate)
		if sErr != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		cRates = append(cRates, cRate)
	}

	if cRates == nil {
		return nil, storage.ErrCurRatesNotFound
	}

	return cRates, nil
}
