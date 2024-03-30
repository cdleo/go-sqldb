package adapter

import (
	"github.com/cdleo/go-commons/sqlcommons"
)

type noopAdapter struct{}

func NewNoopAdapter() sqlcommons.SQLAdapter {
	return &noopAdapter{}
}

func (t *noopAdapter) Translate(query string) string {
	return query
}

func (s *noopAdapter) ErrorHandler(err error) error {
	return err
}
