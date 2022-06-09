package translators

import "github.com/cdleo/go-sqldb"

type noopTranslator struct{}

func NewNoopTranslator() sqldb.SQLSyntaxTranslator {
	return &noopTranslator{}
}

func (t *noopTranslator) Translate(query string) string {
	return query
}
