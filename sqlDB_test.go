package sqldb_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/cdleo/go-commons/sqlcommons"
	"github.com/cdleo/go-sqldb"
	enginesMocks "github.com/cdleo/go-sqldb/engines/mocks"
	translatorsMocks "github.com/cdleo/go-sqldb/translators/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type Customers struct {
	Id         int             `db:"id"`
	Name       string          `db:"name"`
	Updatetime time.Time       `db:"updatetime"`
	Age        sql.NullFloat64 `db:"age"`
	Group      int             `db:"cust_group"`
	Dummy      string          `db:"not_existing_field"`
}

func Test_sqlConn_Init(t *testing.T) {
	// Setup
	controller := gomock.NewController(t)

	adapter := enginesMocks.NewMockSQLEngineAdapter(controller)
	adapter.EXPECT().Open().Return(nil, fmt.Errorf("Can't connect")).Times(1)
	adapter.EXPECT().ErrorHandler(gomock.Any()).Return(sqlcommons.ConnectionFailed).Times(1)
	translator := translatorsMocks.NewMockSQLSyntaxTranslator(controller)

	sqlConn := sqldb.NewSQLDB(adapter, translator)
	require.Error(t, sqlConn.Open())
}

/*
func TestSQL(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Storage Tests Suite")
}

var _ = Describe("Testing: storage", func() {

	Describe("Testing: Database", func() {
		var (
			config         DatabaseConfig
			loggerInstance Logger
		)

		loggerInstance = NewLogger()

		JustBeforeEach(func() {
			config = NewDatabaseConfig()
		})

		Context("When the SqlConn is bad configured", func() {

			It("Can throw an error if connection fails", func() {
				config.DriverName = "postgresql"
				config.Host = "127.0.0.1"
				config.Port = 1234
				config.Database = "public"
				config.User = "user"
				config.Password = "pass"
				sqlConn := engines.NewPostgreSqlConn("127.0.0.1", 1234, "user", "pass", "public")
				err := sqlConn.Open()
				Expect(err).Should(HaveOccurred())
			})

		})

		// Follow tests NEEDS to run over the Veritran network
		Context("When the ORACLE SqlConn is well configured", func() {

			It("Can open a SID connection successfully", func() {
				config.DriverName = "oracle"
				config.Host = "ar-oradb-03.veritran.local"
				config.Port = 1521
				config.Database = "SID=VTDBDEV12"
				config.User = "VTDB_VT5D"
				config.Password = "veritran"

				sqlConn := vtLibSQLImpl.NewSqlClient(config, loggerInstance)
				err := sqlConn.Open()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("Can open a SERVICE NAME connection successfully", func() {
				config.DriverName = "oracle"
				config.Host = "10.241.0.76"
				config.Port = 1521
				config.Database = "SERVICE_NAME=VTDB"
				config.User = "VTDB"
				config.Password = "veritran"

				sqlConn := vtLibSQLImpl.NewSqlClient(config, loggerInstance)
				err := sqlConn.Open()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("Can scan data to struct", func() {
				config.DriverName = "oracle"
				config.Host = "10.241.0.76"
				config.Port = 1521
				config.Database = "SERVICE_NAME=VTDB"
				config.User = "VTDB"
				config.Password = "veritran"
				sqlConn := vtLibSQLImpl.NewSqlClient(config, loggerInstance)
				err := sqlConn.Open()
				Expect(err).ShouldNot(HaveOccurred())

				rows, err2 := sqlConn.QueryX("SELECT * FROM CUSTOMERS c")
				defer rows.Close()
				Expect(err2).ShouldNot(HaveOccurred())

				Expect(rows.Next()).Should(BeTrue())
				var c Customers
				err3 := rows.StructScan(&c)
				Expect(err3).ShouldNot(HaveOccurred())
				Expect(c.Id).ShouldNot(BeZero())
				Expect(c.Name).ShouldNot(BeEmpty())
				Expect(c.Group).ShouldNot(BeZero())
				Expect(c.Dummy).Should(BeEmpty())
			})
		})

		Context("When the SqlConn is well configured", func() {

			JustBeforeEach(func() {
				config.DriverName = "postgresql"
				config.Host = "ar-pgsql-02.veritran.local"
				config.Port = 5432
				config.Database = "vtdb"
				config.User = "vtadmin"
				config.Password = "01SbU8DYr3GG"
			})

			It("Can open the connection successfully", func() {
				sqlConn := vtLibSQLImpl.NewSqlClient(config, loggerInstance)
				err := sqlConn.Open()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("Can drop tables (if exists)", func() {
				sqlConn := vtLibSQLImpl.NewSqlClient(config, loggerInstance)
				err := sqlConn.Open()
				Expect(err).ShouldNot(HaveOccurred())

				_, err = sqlConn.Exec(`DROP TABLE IF EXISTS public.customers_groups CASCADE`)
				Expect(err).ShouldNot(HaveOccurred())

				_, err = sqlConn.Exec(`DROP TABLE IF EXISTS public.customers CASCADE`)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("Can create tables, indexes and constraints ", func() {
				sqlConn := vtLibSQLImpl.NewSqlClient(config, loggerInstance)
				err := sqlConn.Open()
				Expect(err).ShouldNot(HaveOccurred())

				_, err = sqlConn.Exec(`CREATE TABLE IF NOT EXISTS public.customers_groups (
					id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
					groupname varchar NOT NULL,
					CONSTRAINT customers_groups_pk PRIMARY KEY (id))`)
				Expect(err).ShouldNot(HaveOccurred())

				_, err = sqlConn.Exec(`CREATE UNIQUE INDEX customers_groups_id_idx ON public.customers_groups USING btree (id)`)
				Expect(err).ShouldNot(HaveOccurred())

				_, err = sqlConn.Exec(`CREATE TABLE IF NOT EXISTS public.customers (
					id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
					name varchar(10) NOT NULL,
					updatetime timestamp(0) NULL DEFAULT CURRENT_TIMESTAMP,
					age numeric(3,1) NULL,
					cust_group int4 NOT NULL,
					CONSTRAINT customers_un UNIQUE (name))`)
				Expect(err).ShouldNot(HaveOccurred())

				_, err = sqlConn.Exec(`ALTER TABLE public.customers
					ADD CONSTRAINT customers_fk FOREIGN KEY (cust_group) REFERENCES public.customers_groups(id)`)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("Can store data", func() {
				sqlConn := vtLibSQLImpl.NewSqlClient(config, loggerInstance)
				err := sqlConn.Open()
				Expect(err).ShouldNot(HaveOccurred())

				_, err = sqlConn.Exec(`INSERT INTO public.customers_groups (groupname) VALUES('General');`)
				Expect(err).ShouldNot(HaveOccurred())

				statement, err := sqlConn.Prepare("INSERT INTO public.customers (name, updatetime, age, cust_group)VALUES(:1, :2, :3, :4)")
				Expect(err).ShouldNot(HaveOccurred())
				_, err = statement.Exec("Juan", clock.NewSystemClock().CurrentInstant(), nil, 1)
				Expect(err).ShouldNot(HaveOccurred())
				_, err = statement.Exec("Pedro", clock.NewSystemClock().CurrentInstantInUTC(), nil, 1)
				Expect(err).ShouldNot(HaveOccurred())
				_, err = statement.Exec("Pablo", clock.NewSystemClock().CurrentInstant(), 99, 1)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("Can return data", func() {
				sqlConn := vtLibSQLImpl.NewSqlClient(config, loggerInstance)
				err := sqlConn.Open()
				Expect(err).ShouldNot(HaveOccurred())

				rows, err2 := sqlConn.Query("SELECT c.name FROM public.customers c")
				defer rows.Close()

				Expect(err2).ShouldNot(HaveOccurred())
				Expect(rows.Next()).Should(BeTrue())
			})

			It("Can throw an invalid table error", func() {
				sqlConn := vtLibSQLImpl.NewSqlClient(config, loggerInstance)
				err := sqlConn.Open()
				Expect(err).ShouldNot(HaveOccurred())

				_, err = sqlConn.Query("SELECT c.name FROM public.customxers c")
				Expect(err).Should(HaveOccurred())
			})

			It("Can throw a cannot insert null error", func() {
				sqlConn := vtLibSQLImpl.NewSqlClient(config, loggerInstance)
				err := sqlConn.Open()
				Expect(err).ShouldNot(HaveOccurred())

				_, err = sqlConn.Exec("INSERT INTO public.customers (name, updatetime) VALUES (:1,:2)", nil, time.Now())
				Expect(err.Cause()).Should(Equal(vtLibSQL.CannotSetNullColumn))
			})

			It("Can throw a cannot update null error", func() {
				sqlConn := vtLibSQLImpl.NewSqlClient(config, loggerInstance)
				err := sqlConn.Open()
				Expect(err).ShouldNot(HaveOccurred())

				_, err = sqlConn.Exec("UPDATE public.customers c SET name = :1 WHERE c.name = :2", nil, "Pablo")
				Expect(err.Cause()).Should(Equal(vtLibSQL.CannotSetNullColumn))
			})

			It("Can throw a too large value error", func() {
				sqlConn := vtLibSQLImpl.NewSqlClient(config, loggerInstance)
				err := sqlConn.Open()
				Expect(err).ShouldNot(HaveOccurred())

				_, err = sqlConn.Exec("INSERT INTO public.customers (name, updatetime) VALUES (:1,:2)", "1234567891011", time.Now())
				Expect(err.Cause()).Should(Equal(vtLibSQL.ValueTooLargeForColumn))
			})

			It("Can throw a subqury returns more than one row", func() {
				sqlConn := vtLibSQLImpl.NewSqlClient(config, loggerInstance)
				err := sqlConn.Open()
				Expect(err).ShouldNot(HaveOccurred())

				_, err = sqlConn.Query("SELECT name FROM public.customers c WHERE c.id = (SELECT d.id FROM public.customers d)")
				Expect(err.Cause()).Should(Equal(vtLibSQL.SubqueryReturnsMoreThanOneRow))
			})

			It("Can throw an invalid numeric value error", func() {
				sqlConn := vtLibSQLImpl.NewSqlClient(config, loggerInstance)
				err := sqlConn.Open()
				Expect(err).ShouldNot(HaveOccurred())

				_, err = sqlConn.Exec("UPDATE public.customers c SET age = :1 WHERE c.name = :2", "1a2", "Pablo")
				Expect(err.Cause()).Should(Equal(vtLibSQL.InvalidNumericValue))
			})

			It("Can throw a value larger than precision error", func() {
				sqlConn := vtLibSQLImpl.NewSqlClient(config, loggerInstance)
				err := sqlConn.Open()
				Expect(err).ShouldNot(HaveOccurred())

				_, err = sqlConn.Exec("UPDATE public.customers c SET age = :1 WHERE c.name = :2", 949.0044, "Pablo")
				Expect(err.Cause()).Should(Equal(vtLibSQL.ValueLargerThanPrecision))
			})

			It("Can throw a value larger than precision error", func() {
				sqlConn := vtLibSQLImpl.NewSqlClient(config, loggerInstance)
				err := sqlConn.Open()
				Expect(err).ShouldNot(HaveOccurred())

				_, err = sqlConn.Exec("DELETE from public.customers_groups g WHERE g.id = :1", 1)
				Expect(err.Cause()).Should(Equal(vtLibSQL.IntegrityConstraintViolation))
			})

			It("Can scan data to struct", func() {
				sqlConn := vtLibSQLImpl.NewSqlClient(config, loggerInstance)
				err := sqlConn.Open()
				Expect(err).ShouldNot(HaveOccurred())

				rows, err2 := sqlConn.QueryX("SELECT * FROM public.customers c")
				defer rows.Close()
				Expect(err2).ShouldNot(HaveOccurred())

				Expect(rows.Next()).Should(BeTrue())
				var c Customers
				err3 := rows.StructScan(&c)
				Expect(err3).ShouldNot(HaveOccurred())
				Expect(c.Id).ShouldNot(BeZero())
				Expect(c.Name).ShouldNot(BeEmpty())
				Expect(c.Group).ShouldNot(BeZero())
				Expect(c.Dummy).Should(BeEmpty())
			})

		})

	})
})
*/
