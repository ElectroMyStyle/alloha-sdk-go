# Alloha API SDK on Golang
This package provides a set of functions for working with the Alloha API


## Requirements
The SDK supports Go version 1.16 and above.


## Installing the SDK
Install the package using the command:
```bash
go get -u github.com/electromystyle/alloha-sdk-go
```


## Usage example


## API Methods
List of implemented API methods

| # | Method | Description                                  |
|---|--------|----------------------------------------------|
| 1 | FindByIMDbId | Finds a movie by its IMDb ID                 |
| 2 | FindByKPId | Finds a movie by its KP ID                   |
| 3 | FindByTMDbId | Finds a movie by its TMDb ID                 |
| 4 | GetListOfLatestSeries | Search and returns a list of latest series |
| 5 | SearchForOneByName | Searches and returns a single movie by name |
| 6 | SearchListByName | Searches and returns a list of movies by name |


### Testing
To start testing, you can use the command:
```bash
go test -v -timeout 30s ./...
```
If the `make` utility is available on your device, you can run the test with the command:
```bash
make test
```
or simply
```bash
make
```

## License
The Alloha SDK for Go is licensed for use under the terms and conditions of the [MIT license Agreement](https://github.com/electromystyle/alloha-sdk-go/blob/master/LICENSE).

Created by David

2024
