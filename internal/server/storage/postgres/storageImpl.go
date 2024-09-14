package postgres

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"strconv"

	"ops-storage/internal/server/logger"
	serror "ops-storage/internal/server/storage/error"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var insertGauge = "INSERT INTO gauge (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = $2"
var insertCounter = "INSERT INTO counter (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = counter.value + $2"

func (d Database) Insert(counterType string, name string, val string) error {
	var query string
	if counterType == "gauge" {
		_, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return serror.ErrIvalidMetric
		}
		query = insertGauge
	} else {
		_, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return serror.ErrIvalidMetric
		}
		query = insertCounter
	}

	_, err := d.db.ExecContext(context.Background(), query, name, val)
	if err != nil {
		var netErr *net.OpError
		if errors.As(err, &netErr) {
			return d.insertBackoff(query, name, val)
		}
		logger.MainLog.Errorf("Unexpected insert error, err: ", err.Error())
		return serror.ErrInternal
	}

	return nil
}

var getGauge = "SELECT value FROM gauge WHERE name = $1"
var getCounter = "SELECT value FROM counter WHERE name = $1"

func (d Database) Get(counterType string, name string) (string, error) {
	var query string
	if counterType == "gauge" {
		query = getGauge
	} else {
		query = getCounter
	}

	row := d.db.QueryRowContext(context.Background(), query, name)

	var res string
	err := row.Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", serror.ErrNotFound
		} 
		var netErr *net.OpError
		if errors.As(err, &netErr) {
			return d.getBackoff(query, name)
		}
		logger.MainLog.Errorf("Unexpected get error, err: ", err.Error())
		return "", serror.ErrInternal
	}
	return res, err
}

var getAllGauge = "SELECT name, value FROM gauge"
var getAllCounter = "SELECT name, value FROM counter"

type metric struct {
	name string
	val  string
}

func (d Database) GetAll() (*map[string]string, error) {
	res := make(map[string]string)

	err := d.fillSpecType(getAllGauge, &res)
	if err != nil {
		return nil, err
	}

	err = d.fillSpecType(getAllCounter, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (d Database) fillSpecType(query string, res *map[string]string) error {
	rows, err := d.db.QueryContext(context.Background(), query)
	if err != nil || rows.Err() != nil {
		logger.MainLog.Errorf("Unexpected getAll error, err: ", err.Error())
		return serror.ErrInternal
	}
	defer rows.Close()

	var m metric
	for rows.Next() {
		err = rows.Scan(&m.name, &m.val)
		if err != nil {
			logger.MainLog.Errorf("Parse getAll error, err: ", err.Error())
			return serror.ErrInternal
		}

		(*res)[m.name] = m.val
	}

	return nil
}

func (d Database) IsStorageAlive() bool {
	err := d.db.Ping()
	if err != nil {
		logger.MainLog.Warnf("Ping db error, Err: ", err.Error())
		return false
	}
	return true
}
