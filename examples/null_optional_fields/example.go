package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/amandavmanduca/csvToStruct"
)

type ExampleCsv struct {
	ID     string `csv_column:"ID,required"`
	Name   string `csv_column:"Name,omitempty"`
	Age    string `csv_column:"Age,omitempty"`
	Gender string `csv_column:"Gender,omitempty"`
	Date   string `csv_column:"Date,omitempty"`
	Active string `csv_column:"Active,omitempty"`
}

type ExamplePayload struct {
	ID              string     `json:"id" validate:"required"`
	Name            string     `json:"name,omitempty"`
	Age             int64      `json:"age,omitempty"`
	Gender          string     `json:"gender,omitempty"`
	Date            *time.Time `json:"date,omitempty"`
	Active          *bool      `json:"forum_enabled,omitempty"`
	SetFieldsToNull []string   `json:"set_fields_to_null,omitempty"`
}

func (csv ExampleCsv) ToPayload() (ExamplePayload, error) {
	nullFields := []string{}
	date, setDateToNull, err := parseDatePtr(csv.Date)
	if err != nil {
		return ExamplePayload{}, err
	}
	if setDateToNull {
		nullFields = append(nullFields, "Date")
	}
	age, err := strconv.ParseInt(csv.Age, 0, 0)
	if err != nil {
		age = 0
	}
	return ExamplePayload{
		ID:              csv.ID,
		Name:            csv.Name,
		Active:          parseBoolPtr(csv.Active),
		Age:             age,
		Date:            date,
		Gender:          csv.Gender,
		SetFieldsToNull: nullFields,
	}, nil
}

func parseBoolPtr(s string) *bool {
	if strings.TrimSpace(strings.ToLower(s)) == "true" {
		b := true
		return &b
	} else if strings.TrimSpace(strings.ToLower(s)) == "false" {
		b := false
		return &b
	}
	return nil
}

func parseDatePtr(s string) (*time.Time, bool, error) {
	if s == "" {
		return nil, false, nil
	}
	if strings.TrimSpace(strings.ToLower(s)) == "null" {
		return nil, true, nil
	}
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return nil, false, err
	}
	return &t, false, nil
}

func parseFloatPtr(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0, err
	}
	return f, nil
}

func main() {
	file, err := os.Open("example.csv")
	if err != nil {
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	dataArray, rowsWithErrors, err := csvToStruct.CsvHandler[ExampleCsv, ExamplePayload](reader)
	if err != nil {
		fmt.Println("err ", err)
	}
	fmt.Println(dataArray, rowsWithErrors)
}
