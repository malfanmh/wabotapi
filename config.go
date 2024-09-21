package wabotapi

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"sync"
)

var once sync.Once

type Config struct {
	APPPort             string `envconfig:"APP_PORT" default:"8080"`
	MysqlDSN            string `envconfig:"MYSQL_DSN" default:"root:password@tcp(127.0.0.1:3306)/db"`
	WAAccessToken       string `envconfig:"WA_ACCESS_TOKEN" default:"token"`
	WABusinessAccountID string `envconfig:"WA_BUSINESS_ACCOUNT_ID"`
	WABaseURL           string `envconfig:"WA_BASE_URL"`
	WASecret            string `envconfig:"WA_SECRET"`
	FinpayMerchantID    string `envconfig:"FINPAY_MERCHANT_ID" default:""`
	FinpaySecret        string `envconfig:"FINPAY_SECRET" default:""`
	FinpayBaseURL       string `envconfig:"FINPAY_BASE_URL" default:""`
	FinpayCallbackURL   string `envconfig:"FINPAY_CALLBACK_URL" default:""`
}

func init() {
	fmt.Println("Initializing , load environment variables", godotenv.Load())
}

func New() (conf *Config) {
	once.Do(func() {
		conf = new(Config)
		envconfig.MustProcess("", conf)
	})
	return conf
}

func OpenMysqlDB(dsn string) *sqlx.DB {
	fmt.Println(dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(fmt.Errorf("failed to ping database: %s", err))
	}
	return sqlx.NewDb(db, "mysql")
}
