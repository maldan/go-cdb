package parse

import (
	"errors"
	"fmt"
	"github.com/maldan/go-cdb/cdb_proto/core"
	"reflect"
	"strings"
)

type QueryExpression struct {
	//Field    string
	//Value any
	//Type  string
	//LeftOffset  int
	//RightOffset int
	Operator string
	Left     string
	Right    string

	//Left     *QueryExpression
	//Right    *QueryExpression
}

type QueryInfo struct {
	Operation      string
	SelectFields   []string
	WhereCondition []TokenType
	TypeInfo       reflect.Type
}

type TokenType struct {
	Token       string
	Type        uint8
	Value       any
	B           []byte
	TableOffset int
}

func precedence(op TokenType) int {
	switch op.Token {
	case "-":
	case "+":
		return 11
	case "*":
	case "/":
		return 12
	case "==":
		return 8
	case "&&":
		return 4
	default:
		return -1
	}
	return -1
}

func reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func infixToPostfix(tokens []TokenType) []TokenType {
	postfix := make([]TokenType, 0)
	stack := make([]TokenType, 0)

	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == core.TokenOp {
			for {
				if len(stack) == 0 {
					break
				}
				if precedence(top(&stack)) < precedence(tokens[i]) {
					break
				}

				postfix = append(postfix, top(&stack))
				pop(&stack)
			}
			stack = append(stack, tokens[i])
		} else {
			postfix = append(postfix, tokens[i])
		}
	}

	for {
		if len(stack) == 0 {
			break
		}
		postfix = append(postfix, top(&stack))
		pop(&stack)
	}

	fmt.Printf("%v\n", postfix)
	return postfix
}

func pop[T any](s *[]T) T {
	v := (*s)[len(*s)-1]
	*s = (*s)[0 : len(*s)-1]
	return v
}

func top[T any](s *[]T) T {
	return (*s)[len(*s)-1]
}

func parseWhere(queryInfo *QueryInfo, tokens []TokenType) {
	// Change tokes
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Token == "AND" {
			tokens[i].Token = "&&"
			tokens[i].Type = core.TokenOp
		}
	}

	queryInfo.WhereCondition = infixToPostfix(tokens)
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
				out = append(out, TokenType{Token: tempOp, Type: core.TokenOp})
				tempOp = ""
			} else if previousMode == ")" {
				out = append(out, TokenType{Token: ")", Type: core.TokenOp})
			} else if previousMode == "(" {
				out = append(out, TokenType{Token: "(", Type: core.TokenOp})
			} else if previousMode == "*" {
				out = append(out, TokenType{Token: "*", Type: core.TokenOp})
			} else if previousMode == "space" {

			} else if previousMode == "string" {
				out = append(out, TokenType{Token: tempStr, Type: core.TokenIdentifier})
				tempStr = ""
			} else if previousMode == "number" {
				out = append(out, TokenType{Token: tempNumber, Type: core.TokenNumber})
				tempNumber = ""
			} else if previousMode == "quote" {
				out = append(out, TokenType{Token: tempQuote, B: []byte(tempQuote), Type: core.TokenString})
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
