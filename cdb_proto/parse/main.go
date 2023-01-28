package parse

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type QueryExpression struct {
	Field    string
	Operator string
	Value    any
}

type QueryInfo struct {
	Operation    string
	SelectFields []string
	Condition    []QueryExpression
	typeInfo     reflect.Type
}

func parseWhere(queryInfo *QueryInfo, tuples []string) {
	expression := QueryExpression{
		Field:    tuples[0],
		Operator: tuples[1],
		Value:    tuples[2],
	}

	// String value
	if tuples[2][0] == '\'' {
		expression.Value = tuples[2][1 : len(tuples[2])-1]
	}

	queryInfo.Condition = append(queryInfo.Condition, expression)
}

func parseSelect(queryInfo *QueryInfo, tuples []string) {
	fmt.Printf("%v\n", tuples)

	// Check source
	if tuples[0] == "*" {
		for i := 0; i < queryInfo.typeInfo.NumField(); i++ {
			queryInfo.SelectFields = append(queryInfo.SelectFields, queryInfo.typeInfo.Field(i).Name)
		}
	}

	// Check source
	if strings.ToLower(tuples[1]) == "from" && tuples[2] == "table" {
		// ok
	}

	// Check source
	if strings.ToLower(tuples[3]) == "where" {
		parseWhere(queryInfo, tuples[4:])
	}
}

func Query[T any](query string) (QueryInfo, error) {
	queryInfo := QueryInfo{
		typeInfo: reflect.TypeOf(*new(T)),
	}

	tuples := strings.Split(query, " ")

	// Is select query
	if strings.ToLower(tuples[0]) == "select" {
		queryInfo.Operation = "select"
		parseSelect(&queryInfo, tuples[1:])
	} else {
		return queryInfo, errors.New("unknown operator " + tuples[0])
	}

	return queryInfo, nil
}
