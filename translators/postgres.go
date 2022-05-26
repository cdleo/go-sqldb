package translators

import (
	"regexp"
	"strings"

	"github.com/cdleo/go-commons/sqlcommons"
	"github.com/cdleo/go-sqldb"
)

type postgresTranslator struct {
	paramRegExp *regexp.Regexp
	fromEngine  string
}

func NewPostgresTranslator(fromEngine string) sqlcommons.SQLSyntaxTranslator {
	return &postgresTranslator{
		regexp.MustCompile(":[1-9]"),
		fromEngine,
	}
}

func (s *postgresTranslator) Translate(query string) string {

	if s.fromEngine == sqldb.Oracle_Engine {
		return s.paramRegExp.ReplaceAllStringFunc(query, func(m string) string {
			return strings.Replace(m, ":", "$", 1)
		})
	} else {
		return query
	}

}
