package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
)

// Eval walks down the tree and evaluates each subexpression.
func Eval(n *ASTNode) int {
	switch n.Value {
	case "+":
		return Eval(n.Left) + Eval(n.Right)

	case "-":
		if n.Right != nil {
			return Eval(n.Left) - Eval(n.Right)
		} else {
			// Negation
			return 0 - Eval(n.Left)
		}

	case "*":
		return Eval(n.Left) * Eval(n.Right)

	case "/":
		return Eval(n.Left) / Eval(n.Right)

	case "^":
		return int(math.Pow(float64(Eval(n.Left)), float64(Eval(n.Right))))

	case "":
		return Eval(n.Left)

	default:
		val, _ := strconv.Atoi(n.Value)
		return val
	}
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
	fmt.Println(Eval(tree))

	if *dot != "" {
		tree.PrintDot(*dot)
	}

	os.Exit(0)
}
