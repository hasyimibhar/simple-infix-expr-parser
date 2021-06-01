package main

import (
	"flag"
	"fmt"
	"os"
)

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
