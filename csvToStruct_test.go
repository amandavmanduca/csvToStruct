package csvToStruct

import (
	"os"
	"reflect"
	"testing"
)

type ExamplePayload struct {
	ID   int    `csv_column:"id"`
	Name string `csv_column:"name"`
}

func (e ExamplePayload) ToPayload() (ExamplePayload, error) {
	return e, nil
}

func TestCsvHandler(t *testing.T) {
	file, err := os.Open("test.csv")
	if err != nil {
		return
	}
	defer file.Close()
	data, errors, err := CsvHandler[ExamplePayload, ExamplePayload](file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(errors) != 0 {
		t.Fatalf("unexpected errors: %v", errors)
	}

	if len(data) != 2 {
		t.Fatalf("expected 2 records, got %d", len(data))
	}

	if data[0].ID != 1 || data[0].Name != "John Doe" {
		t.Errorf("unexpected data: %+v", data[0])
	}

	if data[1].ID != 2 || data[1].Name != "Jane Smith" {
		t.Errorf("unexpected data: %+v", data[1])
	}
}

func TestTransformCsvToPayload(t *testing.T) {
	columns := []string{"id", "name"}
	rows := []string{"1", "John Doe"}

	payload, err := transformCsvToPayload[ExamplePayload, ExamplePayload](columns, rows)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if payload.ID != 1 || payload.Name != "John Doe" {
		t.Errorf("unexpected payload: %+v", payload)
	}
}

func TestFillStruct(t *testing.T) {
	type TestStruct struct {
		ID   int    `csv_column:"id"`
		Name string `csv_column:"name"`
	}

	columns := []string{"id", "name"}
	rows := []string{"1", "John Doe"}

	result, err := fillStruct[TestStruct](columns, rows)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != 1 || result.Name != "John Doe" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestMapColumns(t *testing.T) {
	columns := []string{"id", "name"}
	rows := []string{"1", "John Doe"}

	colMap := mapColumns(columns, rows)
	if colMap["id"] != "1" || colMap["name"] != "John Doe" {
		t.Errorf("unexpected column map: %+v", colMap)
	}
}

func TestSetFieldValue(t *testing.T) {
	type TestStruct struct {
		ID   int    `csv_column:"id"`
		Name string `csv_column:"name"`
	}

	columns := []string{"id", "name"}
	rows := []string{"1", "John Doe"}
	colMap := mapColumns(columns, rows)

	var structValue TestStruct
	val := reflect.ValueOf(&structValue).Elem()
	typ := reflect.TypeOf(structValue)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		if err := setFieldValue(field, fieldType, colMap); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if structValue.ID != 1 || structValue.Name != "John Doe" {
		t.Errorf("unexpected struct value: %+v", structValue)
	}
}

func TestConvertValue(t *testing.T) {
	tests := []struct {
		value     string
		fieldType reflect.Type
		expected  interface{}
	}{
		{"1", reflect.TypeOf(1), 1},
		{"true", reflect.TypeOf(true), true},
		{"3.14", reflect.TypeOf(3.14), 3.14},
		{"hello", reflect.TypeOf(""), "hello"},
	}

	for _, test := range tests {
		result, err := convertValue(test.value, test.fieldType)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != test.expected {
			t.Errorf("expected %v, got %v", test.expected, result)
		}
	}
}
