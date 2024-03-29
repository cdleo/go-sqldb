package sqldb

import (
	"context"
	"database/sql"
	"time"

	"github.com/cdleo/go-commons/logger"
	//"github.com/cdleo/go-commons/sql"
)

type sqlDB struct {
	engineAdapter SQLEngineAdapter
	translator    SQLSyntaxTranslator
	logger        logger.Logger
	db            *sql.DB
}

func NewSQLDB(adapter SQLEngineAdapter, translator SQLSyntaxTranslator, logger logger.Logger) *sqlDB {
	return &sqlDB{
		adapter,
		translator,
		logger,
		nil,
	}
}

func (s *sqlDB) Open() (*sql.DB, error) {

	if db, err := s.engineAdapter.Open(s.logger, s.translator); err != nil {
		return nil, s.translator.ErrorHandler(err)
	} else {
		s.db = db
	}

	return s.db, nil
}

func (s *sqlDB) IsOpen() error {
	if s.db == nil {
		return DBNotInitialized
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

func (s *sqlDB) Close() error {

	if s.db == nil {
		return DBNotInitialized
	}

	err := s.db.Close()
	s.db = nil
	return err
}
