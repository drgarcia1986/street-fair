# Street Fair API
A simple street fair API

## Stack
This project was made with <3 using:
* Golang 1.16
* Postgres 13.3
* SQLite (for local unit tests)

## Logs
Every log is saved on file called `fair.log`, to change that use the environment variable `FAIR_LOG_FILE_PATH`.
If you want to send logs to stdout, change environment variable `FAIR_LOG_FILE_PATH` to `-`.

## Database
To configure the database access use the follow environment variable:

| Env | Default value |
| --- | --- |
| FAIR_DATABASE_HOST | localhost |
| FAIR_DATABASE_USER | fair |
| FAIR_DATABASE_PASSWORD | fair |
| FAIR_DATABASE_DBNAME | streetfair |
| FAIR_DATABASE_SSL_MODE | disable |
| FAIR_DATABASE_CONNECTION_TIMEOUT | 5 |

## Import data
To starts a fresh database with some data provided by the Prefeitura de São Paulo, you can run the command
`make import FILE_PATH="path of csv file"` (the default value for argument `FILE_PATH` is `./DEINFO_AB_FEIRASLIVRES_2014.csv`).
This command compile and run the importer assuming the default database connection parameters, to change that, take a look
at [Database](#Database).

## API
To starts a new instance of the API, you can run the command `make run`.
This command compile and run the API server assuming the default database connection parameters and default port (8000), to change
the database connection parameters take a look at [database](#Database), to change the API port
use the command line argument `port`, f.ex: `make run PORT=8080`.

**IMPORTANT**: This is a REST API.

### Examples
#### Create a new Street Fair
**POST /**
```
$ curl -i -d '{"longitude":-46450424,"latitude":-23602582,"setcens":"355030833000022","areap":"3550308005274","cod_district":"32","district":"IGUATEMI","cod_sub_city_hall":"30","sub_city_hall":"SAO MATEUS","region_5":"Leste","region_8":"Leste 2","name":"JD.BOA ESPERANCA","registry":"5171-3","address":"RUA IGUPIARA","address_number":"S/N","neighborhood":"JD BOA ESPERANCA","landmark":""}' http://localhost:8000/

HTTP/1.1 201 Created
Content-Type: application/json
Date: Fri, 13 Aug 2021 18:51:28 GMT
Content-Length: 375

{"longitude":-46450424,"latitude":-23602582,"setcens":"355030833000022","areap":"3550308005274","cod_district":"32","district":"IGUATEMI","cod_sub_city_hall":"30","sub_city_hall":"SAO MATEUS","region_5":"Leste","region_8":"Leste 2","name":"JD.BOA ESPERANCA","registry":"5171-3","address":"RUA IGUPIARA","address_number":"S/N","neighborhood":"JD BOA ESPERANCA","landmark":""}
```

#### Retrieve a Street Fair
**GET /{registry}/**
```
$ curl -i http://localhost:8000/5171-3/

HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 13 Aug 2021 18:56:45 GMT
Content-Length: 375

{"longitude":-46450424,"latitude":-23602582,"setcens":"355030833000022","areap":"3550308005274","cod_district":"32","district":"IGUATEMI","cod_sub_city_hall":"30","sub_city_hall":"SAO MATEUS","region_5":"Leste","region_8":"Leste 2","name":"JD.BOA ESPERANCA","registry":"5171-3","address":"RUA IGUPIARA","address_number":"S/N","neighborhood":"JD BOA ESPERANCA","landmark":""}
```

#### Update a Street Fair
**PUT /{registry}/**
```
$ curl -i -X PUT -d '{"longitude":-46450424,"latitude":-23602582,"setcens":"355030833000022","areap":"3550308005274","cod_district":"32","district":"IGUATEMI","cod_sub_city_hall":"30","sub_city_hall":"SAO MATEUS","region_5":"Leste","region_8":"Leste 2","name":"JD.BOA ESPERANCA II","registry":"5171-3","address":"RUA IGUPIARA","address_number":"S/N","neighborhood":"JD BOA ESPERANCA","landmark":""}' http://localhost:8000/5171-3/

HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 13 Aug 2021 19:00:20 GMT
Content-Length: 378

{"longitude":-46450424,"latitude":-23602582,"setcens":"355030833000022","areap":"3550308005274","cod_district":"32","district":"IGUATEMI","cod_sub_city_hall":"30","sub_city_hall":"SAO MATEUS","region_5":"Leste","region_8":"Leste 2","name":"JD.BOA ESPERANCA II","registry":"5171-3","address":"RUA IGUPIARA","address_number":"S/N","neighborhood":"JD BOA ESPERANCA","landmark":""}
```

#### Delete a Street Fair
**DELETE /{registry}/**
```
$ curl -i -X DELETE http://localhost:8000/5171-3/

HTTP/1.1 204 No Content
Date: Fri, 13 Aug 2021 19:02:19 GMT
```

#### Retrieve all Street Fairs
**GET /**
```
$ curl -i http://localhost:8000/

HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 13 Aug 2021 19:04:16 GMT
Content-Length: 813

[{"longitude":-46550164,"latitude":-23558732,"setcens":"355030885000091","areap":"3550308005040","cod_district":"87","district":"VILA FORMOSA","cod_sub_city_hall":"26","sub_city_hall":"ARICANDUVA-FORMOSA-CARRAO","region_5":"Leste","region_8":"Leste 1","name":"VILA FORMOSA","registry":"4041-0","address":"RUA MARAGOJIPE","address_number":"S/N","neighborhood":"VL FORMOSA","landmark":"TV RUA PRETORIA"},{"longitude":-46574716,"latitude":-23584852,"setcens":"355030893000035","areap":"3550308005042","cod_district":"95","district":"VILA PRUDENTE","cod_sub_city_hall":"29","sub_city_hall":"VILA PRUDENTE","region_5":"Leste","region_8":"Leste 1","name":"PRACA SANTA HELENA","registry":"4045-2","address":"RUA JOSE DOS REIS","address_number":"909.000000","neighborhood":"VL ZELINA","landmark":"RUA OLIVEIRA GOUVEIA"}]
```

For this endpoint you can use the following filters:
* _district_
* _region5_
* _name_
* _neighborhood_

e.g.:

```
$ curl -i http://localhost:8000/\?district\=MORUMBI

HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 13 Aug 2021 20:56:05 GMT
Content-Length: 1593

[{"longitude":-46705028,"latitude":-23610576,"setcens":"355030854000048","areap":"3550308005104","cod_district":"55","district":"MORUMBI","cod_sub_city_hall":"10","sub_city_hall":"BUTANTA","region_5":"Oeste","region_8":"Oeste","name":"FEIRAO DA ECONOMIA REAL PARQUE","registry":"5143-8","address":"AV BARAO DE MONTE MOR","address_number":"S/N","neighborhood":"REAL PQ MORUMBI","landmark":""},{"longitude":-46705652,"latitude":-23579220,"setcens":"355030854000027","areap":"3550308005104","cod_district":"55","district":"MORUMBI","cod_sub_city_hall":"10","sub_city_hall":"BUTANTA","region_5":"Oeste","region_8":"Oeste","name":"BIBI","registry":"4012-6","address":"PC ROBERTO GOMES PEDROSA","address_number":"520.000000","neighborhood":"ITAIM BIBI","landmark":"PC ROBERTO GOMES PEDROSA"},{"longitude":-46720092,"latitude":-23599440,"setcens":"355030854000042","areap":"3550308005104","cod_district":"55","district":"MORUMBI","cod_sub_city_hall":"10","sub_city_hall":"BUTANTA","region_5":"Oeste","region_8":"Oeste","name":"CAXINGUI","registry":"3038-4","address":"PC ROBERTO GOMES PEDROSA","address_number":"","neighborhood":"ESTADIO DO MORUMBI","landmark":"AO LADO PC ROBERTO G PEDROSA"},{"longitude":-46705164,"latitude":-23610496,"setcens":"355030854000038","areap":"3550308005104","cod_district":"55","district":"MORUMBI","cod_sub_city_hall":"10","sub_city_hall":"BUTANTA","region_5":"Oeste","region_8":"Oeste","name":"REAL PARQUE","registry":"1089-8","address":"RUA BARAO DE MONTE MOR","address_number":"166.000000","neighborhood":"PAINEIRAS DO MORUMBI","landmark":"RUA BARAO DE C.GERAIS"}]
```

## Docker
For docker users a simple `docker-compose up` starts a fresh database (with all street fairs already imported) and an API instance.
