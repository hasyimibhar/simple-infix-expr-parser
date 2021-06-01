package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
)

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

type ASTNode struct {
	Value string
	Left  *ASTNode
	Right *ASTNode
}

// Reduce simplifies the tree by removing redundant
// empty nodes
func (n *ASTNode) Reduce() *ASTNode {
	if n.Left == nil && n.Right == nil {
		return n
	}
	if n.Value == "" {
		return n.Left.Reduce()
	}

	if n.Left != nil {
		n.Left = n.Left.Reduce()
	}
	if n.Right != nil {
		n.Right = n.Right.Reduce()
	}
	return n
}

func (n *ASTNode) PrintDot(file string) {
	f, _ := os.Create(file)
	defer f.Close()

	f.WriteString("digraph G {\n")
	f.WriteString("\tnode [shape=circle]\n")
	f.WriteString("\n")

	i := 1
	n.print(f, &i)
	f.WriteString("}\n")
}

func (n *ASTNode) print(f *os.File, i *int) {
	id := *i
	f.WriteString(fmt.Sprintf("\tnode_%d [label=\" %s \"]\n", id, n.Value))

	var leftId, rightId int
	if n.Left != nil {
		*i += 1
		leftId = *i
		f.WriteString(fmt.Sprintf("\tnode_%d -> node_%d\n", id, leftId))
		n.Left.print(f, i)
	}

	if n.Right != nil {
		*i += 1
		rightId = *i
		f.WriteString(fmt.Sprintf("\tnode_%d -> node_%d\n", id, rightId))
		n.Right.print(f, i)
	}

	if n.Left != nil && n.Right != nil {
		f.WriteString(fmt.Sprintf("\t{ rank=same; node_%d -> node_%d [style=invis] }\n", leftId, rightId))
	}
}

func (n *ASTNode) Eval() int {
	switch n.Value {
	case "+":
		return n.Left.Eval() + n.Right.Eval()

	case "-":
		if n.Right != nil {
			return n.Left.Eval() - n.Right.Eval()
		} else {
			// Negation
			return 0 - n.Left.Eval()
		}

	case "*":
		return n.Left.Eval() * n.Right.Eval()

	case "/":
		return n.Left.Eval() / n.Right.Eval()

	case "^":
		return int(math.Pow(float64(n.Left.Eval()), float64(n.Right.Eval())))

	case "":
		return n.Left.Eval()

	default:
		val, _ := strconv.Atoi(n.Value)
		return val
	}
}

/**

The code below implements the following LL(1) grammar:

expr ::= addSubExpr

addSubExpr ::= mulDivExpr addSubExprTail
addSubExprTail ::= ('+' | '-') mulDivExpr addSubExprTail
						 		 | empty

mulDivExpr ::= powExpr mulDivExprTail
mulDivExprTail ::= ('*' | '/') powExpr mulDivExprTail
						 		 | empty

powExpr ::= parenExpr powExprTail
				  | '-' parenExpr
powExprTail ::= '^' parenExpr powExprTail
						  | empty

parenExpr ::= '(' expr ')'
            | NUMBER

NUMBER ::= ['0'..'9']+

**/

func expr(parent *ASTNode, ts *TokenStream) {
	addSubExpr(parent, ts)
}

func addSubExpr(parent *ASTNode, ts *TokenStream) {
	node := new(ASTNode)
	mulDivExpr(node, ts)

	parent.Left = addSubExprTail(node, ts)
}

func addSubExprTail(parent *ASTNode, ts *TokenStream) *ASTNode {
	if ts.Empty() {
		return parent
	}
	if ts.Peek().Type != TokenTypePlus && ts.Peek().Type != TokenTypeMinus {
		return parent
	}

	node := new(ASTNode)
	node.Value = ts.Pop().Value
	node.Left = parent

	right := new(ASTNode)
	mulDivExpr(right, ts)
	node.Right = right

	return addSubExprTail(node, ts)
}

func mulDivExpr(parent *ASTNode, ts *TokenStream) {
	node := new(ASTNode)
	powExpr(node, ts)

	parent.Left = mulDivExprTail(node, ts)
}

func mulDivExprTail(parent *ASTNode, ts *TokenStream) *ASTNode {
	if ts.Empty() {
		return parent
	}
	if ts.Peek().Type != TokenTypeMult && ts.Peek().Type != TokenTypeDiv {
		return parent
	}

	node := new(ASTNode)
	node.Value = ts.Pop().Value
	node.Left = parent

	right := new(ASTNode)
	powExpr(right, ts)
	node.Right = right

	return mulDivExprTail(node, ts)
}

// powExpr and powExprTail build the AST such
// that it grows rightwards to enforce right-associativity
func powExpr(parent *ASTNode, ts *TokenStream) {
	if ts.Peek().Type == TokenTypeMinus {
		// Handle negation
		node := new(ASTNode)
		node.Value = ts.Pop().Value

		left := new(ASTNode)
		parenExpr(left, ts)
		node.Left = left

		parent.Left = node
	} else {
		left := new(ASTNode)
		parenExpr(left, ts)

		right := new(ASTNode)
		right = powExprTail(right, ts)

		if right != nil {
			node := new(ASTNode)
			node.Value = "^"
			node.Left = left
			node.Right = right

			parent.Left = node
		} else {
			parent.Left = left
		}
	}
}

func powExprTail(parent *ASTNode, ts *TokenStream) *ASTNode {
	if ts.Empty() {
		return nil
	}
	if ts.Peek().Type != TokenTypePow {
		return nil
	}

	ts.Pop()

	left := new(ASTNode)
	parenExpr(left, ts)

	right := new(ASTNode)
	right = powExprTail(right, ts)

	if right != nil {
		node := new(ASTNode)
		node.Value = "^"
		node.Left = left
		node.Right = right
		return node
	} else {
		return left
	}
}

func parenExpr(parent *ASTNode, ts *TokenStream) {
	if ts.Peek().Type == TokenTypeOpenParen {
		ts.Pop()
		expr(parent, ts)
		ts.Pop()
	} else {
		parent.Value = ts.Pop().Value
	}
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

// Parse converts stream of tokens into an abstract syntax tree
func Parse(ts *TokenStream) *ASTNode {
	root := new(ASTNode)
	expr(root, ts)
	root = root.Reduce()
	return root
}

func main() {
	var dot = flag.String("dot", "", "output AST as dot file")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("usage: expr [-dot ast.dot] <expr>")
		fmt.Println("example: expr \"1+2+3\"")
		os.Exit(1)
	}

	expr := flag.Arg(0)

	tree := Parse(Lex(expr))
	fmt.Println(tree.Eval())

	if *dot != "" {
		tree.PrintDot(*dot)
	}

	os.Exit(0)
}
