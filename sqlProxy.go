package sqldb

import (
	"context"
	"database/sql"
	"time"

	"github.com/cdleo/go-commons/logger"
	"github.com/cdleo/go-commons/sqlcommons"
)

type SQLProxy struct {
	connector  sqlcommons.SQLConnector
	translator sqlcommons.SQLAdapter
	logger     logger.Logger
	db         *sql.DB
}

func (s *SQLProxy) Open() (*sql.DB, error) {
	if db, err := s.connector.Open(s.logger, s.translator); err != nil {
		return nil, s.translator.ErrorHandler(err)
	} else {
		s.db = db
	}

	return s.db, nil
}

func (s *SQLProxy) IsOpen() error {
	if s.db == nil {
		return sqlcommons.DBNotInitialized
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if stdErr := s.db.PingContext(ctx); stdErr != nil {
		s.Close()
		if _, err := s.Open(); err != nil {
			return s.translator.ErrorHandler(err)
		} else {
			ctxReconnect, cancelReconnect := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancelReconnect()
			return s.db.PingContext(ctxReconnect)
		}
	}
	return nil
}

func (s *SQLProxy) Close() error {

	if s.db == nil {
		return sqlcommons.DBNotInitialized
	}

	err := s.db.Close()
	s.db = nil
	return err
}

func (s *SQLProxy) GetNextSequenceValue(ctx context.Context, sequenceName string) (int64, error) {
	if err := s.IsOpen(); err != nil {
		return 0, sqlcommons.ConnectionClosed
	}

	query := s.connector.GetNextSequenceQuery(sequenceName)
	row := s.db.QueryRowContext(ctx, query)
	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, sqlcommons.NextValueFailed
	}
	return id, nil
}
