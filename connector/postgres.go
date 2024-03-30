package connector

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"strings"

	"github.com/cdleo/go-commons/logger"
	"github.com/cdleo/go-commons/sqlcommons"
	pgx "github.com/jackc/pgx/v4"
	stdlib "github.com/jackc/pgx/v4/stdlib"
)

type pgSqlConn struct {
	host      string
	port      int
	user      string
	password  string
	database  string
	sslMode   string
	TLSConfig *tls.Config
}

const postgresProxyName = "pgx-proxy"

func NewPostgreSqlConnector(host string, port int, user string, password string, database string) sqlcommons.SQLConnector {

	return &pgSqlConn{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		database: database,
		sslMode:  "disable",
	}
}

func (s *pgSqlConn) WithTLS(sslMode string, allowInsecure bool, serverName string, serverCertificate string, clientCertificate string, clientKey string) error {

	config := &tls.Config{
		InsecureSkipVerify: allowInsecure,
		ServerName:         serverName,
	}

	if serverCertificate != "" {
		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM([]byte(serverCertificate))
		if !ok {
			return fmt.Errorf("unable to append Certs from PEM")
		}
		config.RootCAs = caCertPool
	}

	if clientCertificate != "" && clientKey != "" {
		keypair, err := tls.X509KeyPair([]byte(clientCertificate), []byte(clientKey))
		if err != nil {
			return fmt.Errorf("unable to create keypair of client [%v]", err)
		}
		config.Certificates = []tls.Certificate{keypair}
	}

	s.TLSConfig = config
	s.sslMode = sslMode
	return nil
}

func (s *pgSqlConn) Open(logger logger.Logger, translator sqlcommons.SQLAdapter) (*sql.DB, error) {

	registerProxy(postgresProxyName, logger, translator, stdlib.GetDefaultDriver())

	psqlConn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v", s.host, s.port, s.user, s.password, s.database, s.sslMode)

	config, err := pgx.ParseConfig(psqlConn)
	if err != nil {
		return nil, err
	}
	config.TLSConfig = s.TLSConfig

	dbURI := stdlib.RegisterConnConfig(config)
	dbPool, err := sql.Open(postgresProxyName, dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	return dbPool, nil
}

func (s *pgSqlConn) GetNextSequenceQuery(sequenceName string) string {
	return fmt.Sprintf("SELECT nextval('%s')", strings.ToLower(sequenceName))
}
