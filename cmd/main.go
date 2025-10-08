package main

import (
	"fmt"
	"github.com/suraj-9849/hindiLang.git/config/functions"
	"strings"
)

func main() {
	runTest(`ye name = 1234`)
	runTest(`
		ye count = 0
		jabtak count < 10 {
			bol(count)
			count = count + 1
		}
	`)
	runTest(`
		dohraye {
			bol("hindiScript")
			roko
		}
	`)
	runTest(`
		firseKaro add(a, b): number {
			wapas bhejo a + b
		}
	`)
	runTest(`
		ye i = 0
		jabtak i < 5 {
			i = i + 1
			aage badho
			bol(i)
		}
	`)
}

func runTest(code string) {
	code = strings.TrimSpace(code)
	tokens := functions.Lexer(code)
	fmt.Printf("Tokens: %v\n", tokens)
	ast := functions.Parser(tokens)
	fmt.Printf("AST Body Length: %d\n", len(ast.Body))
	if len(ast.Body) > 0 {
		fmt.Println("AST Nodes:")
		for i, node := range ast.Body {
			fmt.Printf("  [%d] %s\n", i, node.NodeType())
		}
	}
	fmt.Println()
}
