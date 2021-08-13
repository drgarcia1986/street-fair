package main

import (
	"flag"

	"github.com/drgarcia1986/street-fair/pkg/database"
	"github.com/drgarcia1986/street-fair/pkg/fair"
	"github.com/drgarcia1986/street-fair/pkg/importer"
	"github.com/drgarcia1986/street-fair/pkg/logs"
	"github.com/sirupsen/logrus"
)

func main() {
	log, loggerFinalizer, err := logs.New()
	if err != nil {
		panic(err)
	}
	defer loggerFinalizer()

	db, err := database.New()
	if err != nil {
		log.Fatalf("Connecting to database: %+v", err)
	}

	sf, err := fair.New(db, log)
	if err != nil {
		log.Fatalf("Error loading the StreetFair module: %+v", err)
	}
	imp := importer.New(log, sf)

	filePath := flag.String("path", "./DEINFO_AB_FEIRASLIVRES_2014.csv", "The path of file with street fairs data")
	flag.Parse()

	if err := imp.Run(*filePath); err != nil {
		log.WithFields(logrus.Fields{
			"file": filePath,
		}).Fatalf("Error importing file: %+v", err)
	}
	log.Info("Finished")
}
