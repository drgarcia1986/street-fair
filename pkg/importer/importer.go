package importer

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"

	"github.com/drgarcia1986/street-fair/pkg/fair"
	"github.com/sirupsen/logrus"
)

const (
	LONG = iota + 1
	LAT
	SETCENS
	AREAP
	CODDIST
	DISTRITO
	CODSUBPREF
	SUBPREFE
	REGIAO5
	REGIAO8
	NOME_FEIRA
	REGISTRO
	LOGRADOURO
	NUMERO
	BAIRRO
	REFERENCIA
)

type streetFairCreator interface {
	Create(m *fair.Model) (*fair.Model, error)
}

type Importer struct {
	log *logrus.Logger
	sf  streetFairCreator
}

func readFile(filePath string) ([][]string, error) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer csvFile.Close()
	lines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, errors.New("Invalid CSV file")
	}
	return lines[1:], nil
}

func (imp *Importer) parseFloat(number, fieldName, registry string) float64 {
	n, err := strconv.ParseFloat(number, 32)
	if err != nil {
		imp.log.WithField(
			"registry", registry,
		).Warningf("Invalid %s: %+v", fieldName, err)
		return 0
	}
	return n
}

func (imp *Importer) Run(filePath string) error {
	lines, err := readFile(filePath)
	if err != nil {
		return err
	}

	imp.log.WithField("count", len(lines)).Info("Starting")
	for _, line := range lines {
		m := &fair.Model{
			Longitude:      imp.parseFloat(line[LONG], "longitude", line[REGISTRO]),
			Latitude:       imp.parseFloat(line[LAT], "latitude", line[REGISTRO]),
			Setcens:        line[SETCENS],
			Areap:          line[AREAP],
			CodDistrict:    line[CODDIST],
			District:       line[DISTRITO],
			CodSubCityHall: line[CODSUBPREF],
			SubCityHall:    line[SUBPREFE],
			Region5:        line[REGIAO5],
			Region8:        line[REGIAO8],
			Name:           line[NOME_FEIRA],
			Registry:       line[REGISTRO],
			Address:        line[LOGRADOURO],
			AddressNumber:  line[NUMERO],
			Neighborhood:   line[BAIRRO],
			Landmark:       line[REFERENCIA],
		}
		if _, err = imp.sf.Create(m); err != nil {
			imp.log.WithField("registry", line[REGISTRO]).Warning("Skipped")
		}
	}
	return nil
}

func New(log *logrus.Logger, sf streetFairCreator) *Importer {
	return &Importer{
		log: log,
		sf:  sf,
	}
}
