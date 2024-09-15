package main

import (
	"csvToStruct"
	"fmt"
	"os"
	"strconv"
)

type CsvColumns struct {
	ID       string `csv_column:"ID"`
	Name     string `csv_column:"Name"`
	LastName string `csv_column:"Last Name"`
	Number   string `csv_column:"Lucky Number"`
}

type ExamplePayload struct {
	ID          string `json:"id" validate:"required"`
	FullName    string `json:"full_name,omitempty"`
	LuckyNumber int64  `json:"lucky_number,omitempty"`
}

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

func main() {
	file, err := os.Open("example.csv")
	if err != nil {
		return
	}
	defer file.Close()
	dataArray, rowsWithErrors, err := csvToStruct.CsvHandler[CsvColumns, ExamplePayload](file)
	if err != nil {
		fmt.Println("err ", err)
	}
	fmt.Println(dataArray, rowsWithErrors)
}
