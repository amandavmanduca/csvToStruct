package csvToStruct

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
)

type CsvDataWithError struct {
	ErrorMessage string `json:"error_message"`
	Error        string `json:"error"`
	Tag          string `json:"tag"`
	Row          string `json:"row"`
}

type CsvColumnsToPayload[P any] interface {
	ToPayload() (P, error)
}

func CsvHandler[I CsvColumnsToPayload[P], P any](file *os.File) ([]P, []CsvDataWithError, error) {
	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("error reading file:", err)
		return []P{}, []CsvDataWithError{}, err
	}

	var columns = make([]string, 0)
	dataArray := []P{}
	rowsWithErrors := []CsvDataWithError{}
	for index, record := range records {
		// Mount CSV Columns
		if index == 0 {
			columns = record
			if len(columns) == 0 {
				rowsWithErrors = append(rowsWithErrors, CsvDataWithError{
					ErrorMessage: "Invalid columns count",
					Tag:          "CSV_INVALID_COLUMNS_COUNT",
					Row:          strings.Join(record, ", "),
				})
				break
			}
			index++
			continue
		}

		payload, err := transformCsvToPayload[I, P](columns, record)
		if err != nil {
			rowsWithErrors = append(rowsWithErrors, CsvDataWithError{
				ErrorMessage: "Error generating payload",
				Tag:          "CSV_ERROR_SETTING_STRUCT_PAYLOAD",
				Error:        err.Error(),
				Row:          strings.Join(record, ", "),
			})
			continue
		}
		dataArray = append(dataArray, payload)
		index++
	}
	fmt.Println(dataArray, rowsWithErrors)
	return dataArray, rowsWithErrors, nil
}

func transformCsvToPayload[I CsvColumnsToPayload[P], P any](columns []string, rows []string) (P, error) {
	data, err := fillStruct[I](columns, rows)
	if err != nil {
		return *new(P), err
	}
	payload, err := data.ToPayload()
	if err != nil {
		return *new(P), err
	}
	// Validate Filled Struct
	validator := validator.New()
	if err := validator.Struct(payload); err != nil {
		return payload, fmt.Errorf("INVALID_PAYLOAD_ROW_DATA: %w", err)
	}
	return payload, nil
}

func fillStruct[T any](columns []string, rows []string) (T, error) {
	var result T
	structValue := reflect.New(reflect.TypeOf(result)).Elem()
	colMap := mapColumns(columns, rows)

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := reflect.TypeOf(result).Field(i)
		if err := setFieldValue(field, fieldType, colMap); err != nil {
			return result, err
		}
	}
	return structValue.Interface().(T), nil
}

func mapColumns(columns []string, rows []string) map[string]string {
	colMap := make(map[string]string)
	for i, col := range columns {
		if i < len(rows) {
			colMap[col] = rows[i]
		}
	}
	return colMap
}

func setFieldValue(field reflect.Value, fieldType reflect.StructField, colMap map[string]string) error {
	tag := fieldType.Tag.Get("csv_column")
	val, ok := colMap[strings.Split(tag, ",")[0]]
	if !ok {
		return fmt.Errorf("column not found: %s", tag)
	}

	convertedVal, err := convertValue(val, fieldType.Type)
	if err != nil {
		return err
	}

	field.Set(reflect.ValueOf(convertedVal))
	return nil
}

func convertValue(value string, fieldType reflect.Type) (interface{}, error) {
	switch fieldType.Kind() {
	case reflect.Bool:
		return strconv.ParseBool(value)
	case reflect.Int:
		return strconv.Atoi(value)
	case reflect.Float64:
		return strconv.ParseFloat(value, 64)
	case reflect.String:
		return value, nil
	default:
		return nil, fmt.Errorf("unsupported field type: %s", fieldType.Kind())
	}
}
