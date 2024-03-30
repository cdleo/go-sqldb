package sqlproxy

import (
	"github.com/cdleo/go-commons/logger"
	"github.com/cdleo/go-commons/sqlcommons"
	"github.com/cdleo/go-sqldb/adapter"
)

type sqlProxyBuilder struct {
	proxy sqlProxy
}

func NewSQLProxyBuilder(connector sqlcommons.SQLConnector) *sqlProxyBuilder {
	return &sqlProxyBuilder{
		proxy: sqlProxy{
			connector:  connector,
			translator: adapter.NewNoopAdapter(),
			logger:     logger.NewNoLogLogger(),
			db:         nil,
		},
	}
}

func (s *sqlProxyBuilder) WithAdapter(translator sqlcommons.SQLAdapter) *sqlProxyBuilder {
	s.proxy.translator = translator
	return s
}

func (s *sqlProxyBuilder) WithLogger(logger logger.Logger) *sqlProxyBuilder {
	s.proxy.logger = logger
	return s
}

func (s *sqlProxyBuilder) Build() *sqlProxy {
	return &s.proxy
}
