package main

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

// Parse converts stream of tokens into an abstract syntax tree
func Parse(ts *TokenStream) *ASTNode {
	root := new(ASTNode)
	expr(root, ts)
	root = root.Reduce()
	return root
}

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
