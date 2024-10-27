# csvToStruct

## Overview

This package is used to generate an array of Structs, based on a CSV file.


## Installing

```sh
$ go get github.com/amandavmanduca/csvToStruct@v1.0.0
```

## Usage

See complete example under [examples](examples/)


### 1. How to create the basic Structs
```go

// 1. Create a struct of the desired CSV columns
type CsvColumns struct {
	ID       string `csv_column:"ID"`
	Name     string `csv_column:"Name"`
	LastName string `csv_column:"Last Name"`
	Number   string `csv_column:"Lucky Number"`
}

// 2. Create a struct of the payload you want to obtain
type ExamplePayload struct {
	ID          string `json:"id" validate:"required"`
	FullName    string `json:"full_name,omitempty"`
	LuckyNumber int64  `json:"lucky_number,omitempty"`
}

// 3. Add a ToPayload method to the CsvColumns returning your ExamplePayload
func (csv CsvColumns) ToPayload() (ExamplePayload, error) {
	num, err := strconv.ParseInt(csv.Number, 0, 0)
	if err != nil {
		num = 0
	}
	return ExamplePayload{
		ID:          csv.ID,
		FullName:    fmt.Sprintf("%s %s", csv.Name, csv.LastName),
		LuckyNumber: num,
	}, nil
}

```


### 2. How to use the handler
```go
	file, err := os.Open("example.csv")
	if err != nil {
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	dataArray, rowsWithErrors, err := csvToStruct.CsvHandler[CsvColumns, ExamplePayload](reader)
	if err != nil {
		fmt.Println("err ", err)
	}
	fmt.Println(dataArray, rowsWithErrors)
```

### 3. Handling Response Data

1. **dataArray** contains all the correct rows of the CSV with the ExamplePayload format

2. **rowsWithErrors** is an array of CsvDataWithError, containing all the incorrect rows information

```go
type CsvDataWithError struct {
	ErrorMessage string `json:"error_message"`
	Error        string `json:"error"`
	Tag          string `json:"tag"`
	Row          string `json:"row"`
} 
```

### Testing

```sh
$ go test
```

## Contributing

I :heart: Open source!

[Follow github guides for forking a project](https://guides.github.com/activities/forking/)

[Follow github guides for contributing open source](https://guides.github.com/activities/contributing-to-open-source/#contributing)

[Squash pull request into a single commit](http://eli.thegreenplace.net/2014/02/19/squashing-github-pull-requests-into-a-single-commit/)

## License

csvToStruct is released under the [MIT license](http://opensource.org/licenses/MIT).