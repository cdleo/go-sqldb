package sqldb

import (
	"github.com/cdleo/go-commons/logger"
	"github.com/cdleo/go-commons/sqlcommons"
	"github.com/cdleo/go-sqldb/adapter"
)

type SQLProxyBuilder struct {
	proxy SQLProxy
}

func NewSQLProxyBuilder(connector sqlcommons.SQLConnector) *SQLProxyBuilder {
	return &SQLProxyBuilder{
		proxy: SQLProxy{
			connector:  connector,
			translator: adapter.NewNoopAdapter(),
			logger:     logger.NewNoLogLogger(),
			db:         nil,
		},
	}
}

func (s *SQLProxyBuilder) WithAdapter(translator sqlcommons.SQLAdapter) *SQLProxyBuilder {
	s.proxy.translator = translator
	return s
}

func (s *SQLProxyBuilder) WithLogger(logger logger.Logger) *SQLProxyBuilder {
	s.proxy.logger = logger
	return s
}

func (s *SQLProxyBuilder) Build() *SQLProxy {
	return &s.proxy
}
