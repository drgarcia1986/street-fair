PORT=8000
FILE_PATH=./DEINFO_AB_FEIRASLIVRES_2014.csv

test:
	@go test ./... -cover

build-api:
	@go build -o street-fair cmd/api/main.go

build-importer:
	@go build -o importer cmd/importer/main.go

import:
	@go run cmd/importer/main.go -path ${FILE_PATH}

run:
	@go run cmd/api/main.go -port ${PORT}
