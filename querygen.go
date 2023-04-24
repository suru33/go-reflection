package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Student struct {
	Age  int    `cn:"age"`
	Name string `cn:"name"`
}

type Employee struct {
	Id         int
	Name       string
	Age        int
	Country    string
	BaseSalary float64 `cn:"base_salary"`
}

func createQuery(i interface{}) (string, error) {
	structType := reflect.TypeOf(i)
	structValue := reflect.ValueOf(i)

	if reflect.ValueOf(i).Kind() == reflect.Struct {
		var columnNames []string
		var columnValues []string
		for i := 0; i < structValue.NumField(); i++ {
			fieldType := structType.Field(i)
			fieldValue := structValue.Field(i)

			cn := fieldType.Tag.Get("cn")
			if len(cn) == 0 {
				cn = strings.ToLower(fieldType.Name)
			}

			columnNames = append(columnNames, cn)

			switch fieldValue.Kind() {
			case reflect.Int:
				columnValues = append(columnValues, fmt.Sprintf("%d", fieldValue.Int()))
			case reflect.Float64, reflect.Float32:
				columnValues = append(columnValues, fmt.Sprintf("%f", fieldValue.Float()))
			case reflect.String:
				columnValues = append(columnValues, fmt.Sprintf("'%s'", fieldValue.String()))
			case reflect.Bool:
				columnValues = append(columnValues, strconv.FormatBool(fieldValue.Bool()))
			default:
				return "", fmt.Errorf("%s: type not supported", fieldValue.Kind())
			}
		}

		query :=
			fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
				strings.ToLower(structType.Name()),
				strings.Join(columnNames, ", "),
				strings.Join(columnValues, ", "),
			)

		return query, nil
	}

	return "", fmt.Errorf("%s is not a struct", structType.Name())
}

func main() {
	s := Student{Age: 15, Name: "Sam"}
	q1, _ := createQuery(s)
	println(q1)

	emp := Employee{
		Id:         100,
		Name:       "John",
		Age:        21,
		Country:    "FI",
		BaseSalary: 1599.65,
	}

	q2, _ := createQuery(emp)
	println(q2)

	a := 100
	if q2, err := createQuery(a); err != nil {
		panic(err)
	} else {
		println(q2)
	}

}
