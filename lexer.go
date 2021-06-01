package main

type TokenType int

const (
	TokenTypeNumber TokenType = iota
	TokenTypePlus
	TokenTypeMinus
	TokenTypeMult
	TokenTypeDiv
	TokenTypePow
	TokenTypeOpenParen
	TokenTypeCloseParen
)

type Token struct {
	Type  TokenType
	Value string
}

type TokenStream struct {
	Values []Token
}

func (ts *TokenStream) Empty() bool {
	return len(ts.Values) == 0
}

func (ts *TokenStream) Peek() Token {
	return ts.Values[0]
}

func (ts *TokenStream) Pop() Token {
	token := ts.Values[0]
	ts.Values = ts.Values[1:]
	return token
}

// Lex converts stream of characters into stream of tokens
func Lex(s string) *TokenStream {
	tokens := &TokenStream{[]Token{}}

	chars := []byte(s)
	for len(chars) > 0 {
		c := chars[0]
		chars = chars[1:]

		switch c {
		case '+':
			tokens.Values = append(tokens.Values, Token{TokenTypePlus, string(c)})
		case '-':
			tokens.Values = append(tokens.Values, Token{TokenTypeMinus, string(c)})
		case '*':
			tokens.Values = append(tokens.Values, Token{TokenTypeMult, string(c)})
		case '/':
			tokens.Values = append(tokens.Values, Token{TokenTypeDiv, string(c)})
		case '^':
			tokens.Values = append(tokens.Values, Token{TokenTypePow, string(c)})
		case '(':
			tokens.Values = append(tokens.Values, Token{TokenTypeOpenParen, string(c)})
		case ')':
			tokens.Values = append(tokens.Values, Token{TokenTypeCloseParen, string(c)})

		// Number
		default:
			val := []byte{c}

			// Consume characters until non-digit is found
			for len(chars) > 0 && chars[0] >= '0' && chars[0] <= '9' {
				val = append(val, chars[0])
				chars = chars[1:]
			}

			tokens.Values = append(tokens.Values, Token{TokenTypeNumber, string(val)})
		}
	}

	return tokens
}
