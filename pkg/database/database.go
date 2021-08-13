package database

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host           string `default:"localhost"`
	User           string `default:"fair"`
	Password       string `default:"fair"`
	DBName         string `envconfig:"dbname" default:"streetfair"`
	SSLMode        string `envconfig:"ssl_mode" default:"disable"`
	ConnectTimeout int    `default:5 split_words:"true"`
}

func New() (*gorm.DB, error) {
	conf := new(Config)
	if err := envconfig.Process("fair_database", conf); err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		conf.Host, conf.User, conf.Password, conf.DBName, conf.SSLMode, conf.ConnectTimeout,
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
