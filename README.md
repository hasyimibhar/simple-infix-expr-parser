This is a simple LL(1) parser for inflix arithmetic expression. It supports:

- integer literals (`[0-9]+`)
- operators `+`, `-`, `*`, `/`, and `^` (power of)
- parantheses
- operator precedence
- operator associativity (`^` is right-associative, while the rest are left-associative)

## Usage:

`expr [-dot ast.dot] <expr>`

Example:

```sh
$ go run main.go "1+2+3"
6

$ go run main.go "((-2+3*5)^3^2)/3"
3534833124
```

If `-dot` flag is provided, the abstract syntax tree is output as a dot file.
