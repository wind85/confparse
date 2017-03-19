package confparse

import "fmt"

type ParserError struct {
	l int
	k string
	s string
	m string
}

func NewParserError(msg, sec, key string, line int) *ParserError {
	return &ParserError{m: msg, s: sec, k: key, l: line}
}

func (e *ParserError) Error() string {
	return fmt.Sprintf("%s : ,line=%d, section=%s, key=%s\n", e.m, e.l, e.s, e.k)
}

var (
	KEY_NOT_FOUND error = fmt.Errorf("key not found ")
	SEC_NOT_FOUND error = fmt.Errorf("sec not found ")
	NOT_BOOL            = fmt.Errorf("Value is not a bool ")
	NOT_INT             = fmt.Errorf("Value is not an int ")
	NOT_FLOAT           = fmt.Errorf("Value is not a float ")
	NOT_STRING          = fmt.Errorf("Value is not a string ")
)
