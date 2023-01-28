package parse

import (
	"errors"
	"reflect"
	"strings"
)

type QueryExpression struct {
	//Field    string
	//Value    any
	Operator string
	Left     any
	Right    any
}

type QueryInfo struct {
	Operation    string
	SelectFields []string
	Condition    QueryExpression
	TypeInfo     reflect.Type
}

type TokenType struct {
	Token string
	Type  string
}

func hasPrecedence(op1, op2 string) bool {
	if op2 == "(" || op2 == ")" {
		return false
	}
	if (op1 == "AND") && (op2 == "==") {
		return false
	}
	if (op1 == "*" || op1 == "/") && (op2 == "+" || op2 == "-") {
		return false
	} else {
		return true
	}
}

var cnt = 0

func applyOp(op string, a any, b any) any {
	if op == "AND" {
		cnt += 1
		return QueryExpression{Operator: "AND", Left: a, Right: b}
	}
	if op == "==" {
		cnt += 1
		return QueryExpression{Operator: "==", Left: a, Right: b}
	}
	return ""
}

func pop[T any](s *[]T) T {
	/*if len(*s) == 0 {
		return ""
	}*/
	v := (*s)[len(*s)-1]
	*s = (*s)[0 : len(*s)-1]
	return v
}

func parseWhere(queryInfo *QueryInfo, tokens []TokenType) {
	values := make([]any, 0)
	ops := make([]string, 0)

	for i := 0; i < len(tokens); i++ {
		if tokens[i].Token == "(" {
			ops = append(ops, tokens[i].Token)
		} else if tokens[i].Token == ")" {
			for {

				op := pop(&ops)
				v1 := pop(&values)
				v2 := pop(&values)
				values = append(values, applyOp(op, v2, v1))
				if ops[len(ops)-1] == "(" {
					break
				}
			}
			pop(&ops)
		} else if tokens[i].Token == "==" || tokens[i].Token == "AND" {
			for {
				if len(ops) == 0 {
					break
				}
				if !hasPrecedence(tokens[i].Token, ops[len(ops)-1]) {
					break
				}

				op := pop(&ops)
				v1 := pop(&values)
				v2 := pop(&values)
				values = append(values, applyOp(op, v2, v1))
			}
			ops = append(ops, tokens[i].Token)
		} else {
			values = append(values, tokens[i].Token)
		}
	}

	for {
		if len(ops) == 0 {
			break
		}

		op := pop(&ops)
		v1 := pop(&values)
		v2 := pop(&values)

		values = append(values, applyOp(op, v2, v1))
	}

	// cmhp_print.Print(pop(&values))
	//fmt.Printf("%v\n", ops)
	// fmt.Printf("%v\n", pop(&values))
	//fmt.Printf("%v\n", queryInfo.Condition)
	queryInfo.Condition = pop(&values).(QueryExpression)
}

func parseSelect(queryInfo *QueryInfo, tokens []TokenType) {
	// fmt.Printf("%v\n", tuples)

	// Check source
	if tokens[0].Token == "*" {
		for i := 0; i < queryInfo.TypeInfo.NumField(); i++ {
			queryInfo.SelectFields = append(queryInfo.SelectFields, queryInfo.TypeInfo.Field(i).Name)
		}
	}

	// Check source
	if strings.ToLower(tokens[1].Token) == "from" && tokens[2].Token == "table" {
		// ok
	}

	// Check source
	if strings.ToLower(tokens[3].Token) == "where" {
		parseWhere(queryInfo, tokens[4:])
	}
}

func tokenizer(str string) []TokenType {
	out := make([]TokenType, 0)
	tempStr := ""
	tempNumber := ""
	mode := ""
	previousMode := ""
	isQuoteMode := false
	tempQuote := ""
	str += " "
	tempOp := ""

	for i := 0; i < len(str); i++ {
		if isQuoteMode {
			if str[i] == '\'' {
				isQuoteMode = false
				mode = ""
				continue
			}
			tempQuote += string(str[i])
			continue
		}

		if str[i] == '\'' {
			mode = "quote"
			isQuoteMode = true
		} else if str[i] == ' ' {
			mode = "space"
		} else if str[i] >= '0' && str[i] <= '9' {
			tempNumber += string(str[i])
			mode = "number"
		} else if str[i] == '(' {
			mode = "("
		} else if str[i] == ')' {
			mode = ")"
		} else if str[i] == '*' {
			mode = "*"
		} else if str[i] == '=' {
			tempOp += "="
			mode = "="
		} else {
			tempStr += string(str[i])
			mode = "string"
		}

		if mode != previousMode {
			if previousMode == "=" {
				out = append(out, TokenType{Token: tempOp, Type: "op"})
				tempOp = ""
			} else if previousMode == ")" {
				out = append(out, TokenType{Token: ")", Type: "op"})
			} else if previousMode == "(" {
				out = append(out, TokenType{Token: "(", Type: "op"})
			} else if previousMode == "*" {
				out = append(out, TokenType{Token: "*", Type: "op"})
			} else if previousMode == "space" {

			} else if previousMode == "string" {
				out = append(out, TokenType{Token: tempStr, Type: "var"})
				tempStr = ""
			} else if previousMode == "number" {
				out = append(out, TokenType{Token: tempNumber, Type: "number"})
				tempNumber = ""
			} else if previousMode == "quote" {
				out = append(out, TokenType{Token: tempQuote, Type: "string"})
				tempQuote = ""
			}

			previousMode = mode
		}

	}
	return out
}

func Query[T any](query string) (QueryInfo, error) {
	queryInfo := QueryInfo{
		TypeInfo: reflect.TypeOf(*new(T)),
	}

	tokens := tokenizer(query)
	// cmhp_print.Print(tokens)

	if strings.ToLower(tokens[0].Token) == "select" {
		queryInfo.Operation = "select"
		parseSelect(&queryInfo, tokens[1:])
	} else {
		return queryInfo, errors.New("unknown operator " + tokens[0].Token)
	}

	return queryInfo, nil
}
