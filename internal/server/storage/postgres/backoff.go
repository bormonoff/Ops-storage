package postgres

import (
	"context"
	"time"

	"ops-storage/internal/server/logger"
	serror "ops-storage/internal/server/storage/error"
)

var retries = []int{1,3,5}

func (d Database) insertBackoff(query string, args ...string) error {
	var err error

	for _, val := range retries {
		time.Sleep(time.Second * time.Duration(val))

		_, err = d.db.ExecContext(context.Background(), query, args)
		if err == nil {
			return nil
		}
	}

	logger.MainLog.Warnf("Backoff insert error, err: ", err.Error())
	return serror.ErrInternal
}

func (d Database) getBackoff(query string, args ...string) (res string, err error) {
	for _, val := range retries {
		time.Sleep(time.Second * time.Duration(val))

		row := d.db.QueryRowContext(context.Background(), query, args)
		err = row.Scan(&res)
		if err == nil {
			return res, nil
		}
	}

	logger.MainLog.Warnf("Backoff get error, err: ", err.Error())
	return "", serror.ErrInternal
}

func (d Database) pingBackoff() (err error) {
	logger.MainLog.Infof("Try to establish pg connection")
	
	for _, val := range retries {
		err = d.db.PingContext(context.Background())
		if err == nil {
			return nil
		}

		time.Sleep(time.Second * time.Duration(val))
	}

	logger.MainLog.Errorf("Backoff ping error, err: ", err.Error())
	return err
}
