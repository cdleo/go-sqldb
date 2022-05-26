package translators

import (
	"github.com/cdleo/go-commons/sqlcommons"
)

type noopTranslator struct{}

func NewNoopTranslator() sqlcommons.SQLSyntaxTranslator {
	return &noopTranslator{}
}

func (t *noopTranslator) Translate(query string) string {
	return query
}
