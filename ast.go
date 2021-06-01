package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

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
