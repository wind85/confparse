package confparse

import (
	"fmt"
	"io"
	"strconv"
)

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
	KEY_NOT_FOUND error = fmt.Errorf("key not found\n")
	SEC_NOT_FOUND error = fmt.Errorf("sec not found\n")
	NOT_BOOL            = fmt.Errorf("Value is not a bool\n")
	NOT_INT             = fmt.Errorf("Value is not an int\n")
	NOT_FLOAT           = fmt.Errorf("Value is not a float\n")
	NOT_STRING          = fmt.Errorf("Value is not a string\n")
)

type Parser struct {
	s   *Lexer
	buf struct {
		tok    Token
		values []string
		n      int
	}
}

func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewLexer(r)}
}

func (p *Parser) scan() (item *itemType) {
	// If we have a token on the buffer, then return it.
	if p.buf.values == nil {
		p.buf.values = make([]string, 0)
	}

	if p.buf.n != 0 {
		p.buf.n = 0
		item.Tok = p.buf.tok
		item.Values = append(item.Values, p.buf.values...)
	}

	// Otherwise read the next token from the scanner.
	item = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok = item.Tok
	p.buf.values = append(p.buf.values, item.Values...)

	return
}

func (p *Parser) unscan() { p.buf.n = 1 }

//parser does not take into consideration whitespaces ever
func (p *Parser) Parse() (item *itemType) {
	item = p.scan()
	if item.Tok == WHITESPACE {
		item = p.scan()
	}
	return
}

type Config struct {
	C map[string]map[string]string
}

func NewConfig() *Config {
	conf := &Config{C: make(map[string]map[string]string, 0)}
	conf.C["default"] = make(map[string]string, 0)
	conf.C["default"]["version"] = "0.1"
	return conf

}

func (c *Config) getValue(section, key string) (string, error) {
	sec, ok := c.C[section]
	if !ok {
		return "", SEC_NOT_FOUND
	}
	val, ok := sec[key]
	if !ok {
		return "", KEY_NOT_FOUND
	}

	return val, nil

}

type IniParser struct {
	p *Parser
	c *Config
}

func NewIniParser(r io.Reader) *IniParser {
	return &IniParser{p: NewParser(r), c: NewConfig()}
}

func (i *IniParser) Parse() {
	var lastsection string

	for {
		item := i.p.Parse()

		switch {
		case item.Tok == EOF:
			return
		case item.Tok == KEY_VALUE:
			i.c.C[lastsection][item.Values[0]] = item.Values[1]
		case item.Tok == SECTION:
			lastsection = item.Values[0]
			i.c.C[item.Values[0]] = make(map[string]string, 0)

		}
	}
}

func (i *IniParser) GetBool(section, key string) (bool, error) {
	value, er := i.c.getValue(section, key)
	if er != nil {
		return false, er
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		return false, NewParserError(err.Error(), section, key, i.errorLine(key))
	}

	return b, nil

}

func (i *IniParser) GetInt(section, key string) (int64, error) {
	value, err := i.c.getValue(section, key)
	if err != nil {
		return -1, err
	}
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return -1, NewParserError(err.Error(), section, key, i.errorLine(key))
	}

	return n, nil

}

func (i *IniParser) GetFloat(section, key string) (float64, error) {
	value, err := i.c.getValue(section, key)
	if err != nil {
		return -0.1, err
	}
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return -1, NewParserError(err.Error(), section, key, i.errorLine(key))

	}

	return f, nil

}

func (i *IniParser) GetString(section, key string) (string, error) {
	value, err := i.c.getValue(section, key)
	if err != nil {
		return " ", NewParserError(err.Error(), section, key, i.errorLine(key))
	}
	return value, nil
}

func (i *IniParser) errorLine(word string) int {
	lineno, _ := i.p.s.findLine(word)
	return lineno

}
