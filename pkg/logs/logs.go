package logs

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	FilePath string `default:"fair.log" split_words:"true"`
}

func New() (*logrus.Logger, func() error, error) {
	conf := new(Config)
	if err := envconfig.Process("fair_log", conf); err != nil {
		return nil, nil, err
	}
	log := logrus.New()

	finalizer := func() error { return nil }
	if conf.FilePath != "-" {
		f, err := os.OpenFile(conf.FilePath, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			return nil, nil, err
		}
		finalizer = f.Close
		log.SetOutput(f)
	} else {
		log.SetOutput(os.Stdout)
	}
	log.SetFormatter(&logrus.JSONFormatter{})
	return log, finalizer, nil
}
