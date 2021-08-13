package main

import (
	"flag"

	"github.com/drgarcia1986/street-fair/pkg/api"
	"github.com/drgarcia1986/street-fair/pkg/database"
	"github.com/drgarcia1986/street-fair/pkg/fair"
	"github.com/drgarcia1986/street-fair/pkg/logs"
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
		log.Fatalf("Creating new Street Fair instance: %+v", err)
	}

	port := flag.Int("port", 8000, "The port to bind")
	flag.Parse()

	httpSvc := fair.NewHTTPService(sf)
	server := api.NewServer(*port, log)
	httpSvc.RegisterHandlers(server.Router)

	if err := server.Run(); err != nil {
		log.Errorf("Running Server: %v", err)
	}
}
