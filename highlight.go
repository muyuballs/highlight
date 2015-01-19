package highlight

import (
	"bytes"
	"encoding/binary"
	"text/scanner"
	"unicode"
	"unicode/utf8"
)

const (
	KEYWORD = iota
	TYPE
	PLAIN
	NUMBER
	STRING
	COMMENT
	WHITESPACE
	PUNCTUATION
)

type Styntax struct {
	Type  int
	Start int32
	End   int32
	Line  int32
	Text  string
}

func (s *Styntax) Mini() []byte {
	b_buf := bytes.NewBuffer([]byte{})
	b_buf.WriteByte(byte(s.Type))
	binary.Write(b_buf, binary.BigEndian, s.Start)
	binary.Write(b_buf, binary.BigEndian, s.End)
	binary.Write(b_buf, binary.BigEndian, s.Line)
	return b_buf.Bytes()
}

func HighLight(data []byte) []Styntax {
	taxes := make([]Styntax, 0)
	s := NewScanner(data)
	tok := s.Scan()
	startPos := 0
	for tok != scanner.EOF {
		tokText := s.TokenText()
		p := s.Pos()

		tax := Styntax{
			Type:  tokenKind(tok, tokText),
			End:   int32(p.Offset),
			Line:  int32(p.Line),
			Text:  tokText,
			Start: int32(startPos),
		}
		startPos = p.Offset
		taxes = append(taxes, tax)
		tok = s.Scan()
	}
	return taxes
}

func tokenKind(tok rune, tokText string) int {
	switch tok {
	case scanner.Ident:
		if _, isKW := keywords[tokText]; isKW {
			return KEYWORD
		}
		if r, _ := utf8.DecodeRuneInString(tokText); unicode.IsUpper(r) {
			return TYPE
		}
		return PLAIN
	case scanner.Float, scanner.Int:
		return NUMBER
	case scanner.Char, scanner.String, scanner.RawString:
		return STRING
	case scanner.Comment:
		return COMMENT
	}
	if unicode.IsSpace(tok) {
		return WHITESPACE
	}
	return PUNCTUATION
}

func NewScanner(src []byte) *scanner.Scanner {
	var s scanner.Scanner
	s.Init(bytes.NewReader(src))
	s.Error = func(_ *scanner.Scanner, _ string) {}
	s.Whitespace = 0
	s.Mode = s.Mode ^ scanner.SkipComments
	return &s
}
