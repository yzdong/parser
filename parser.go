package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type BracketToken struct {
	open bool
}

func (t BracketToken) printValue() string {
	if t.open {
		return "("
	} else {
		return ")"
	}
}

type ValToken struct {
	value string
}

func (t ValToken) printValue() string {
	return t.value
}

type NodeToken struct {
	tokens []Token
}

func (t NodeToken) printValue() string {
	var str string
	ln := len(t.tokens)
	for i := ln - 1; i >= 0; i-- {
		tkn := t.tokens[i]
		switch tkn.(type) {
		case *ValToken:
			str = str + tkn.printValue() + ", "
		case *BracketToken:
			if !tkn.(*BracketToken).open {
				if len(str) > 1 && str[len(str)-1] == ' ' {
					str = str[0 : len(str)-2]
				}
			}
			str = str + tkn.printValue()
		default:
			str = str + tkn.printValue()
		}
	}
	return str
}

/* Each expression is stored in a NodeToken, which has a tokens stack. (Stack implementation is incomplete)
 */

func (t *NodeToken) push(tkn Token) {
	t.tokens = append(t.tokens, tkn)
}

func (t *NodeToken) pop() Token {
	if len(t.tokens) != 0 {
		tkn := t.tokens[len(t.tokens)-1]
		t.tokens = t.tokens[0 : len(t.tokens)-1]
		return tkn
	} else {
		return nil
	}
}

/* Token interface allows tokens with different types to be
in the value stack
*/

type Token interface {
	printValue() string
}

/* toToken and makeBracketToken are helper functions that make tokens
 */

func toToken(exp string) Token {
	return &ValToken{value: exp}
}

func makeBracketToken(open bool) Token {
	if open {
		return &BracketToken{open: true}
	} else {
		return &BracketToken{open: false}
	}
}

/* addNodeToken creates a NodeToken that contains a completed expression to be evaluated. Returns an error if expressions are malformed.
 */

func (t *NodeToken) addNodeToken() error {
	ntkn := &NodeToken{tokens: []Token{}}
	for {
		tkn := t.pop()
		if tkn == nil {
			break
		}
		switch tkn.(type) {
		case *BracketToken:
			ntkn.push(tkn)
			if tkn.(*BracketToken).open {
				t.push(ntkn) //if (, expression is complete
				return nil
			}
		default:
			ntkn.push(tkn)
			//todo: error handling
		}
	}
	return nil
}

/* tokenize takes the string expression and parses it character by character.
 */

func (t *NodeToken) tokenize(exp string) {
	var value string //token stores the characters of a value
	for _, char := range exp {
		switch {
		case char == '(':
			btkn := makeBracketToken(true)
			t.push(btkn)
		case char == ')':
			//if close bracket, push both token and bracket token
			vtkn := toToken(value)
			value = ""
			t.push(vtkn)
			btkn := makeBracketToken(false)
			t.push(btkn)
			t.addNodeToken() //if there is a closing bracket, try to create an expression
		case char == ' ':
			//if space, push token
			vtkn := toToken(value)
			value = ""
			t.push(vtkn)
		default:
			//add characters to value
			value = value + string(char)
		}
	}
}

func main() {
	fmt.Println("Type in an expression to be evaluated: ")
	t := &NodeToken{tokens: []Token{}}
	reader := bufio.NewReader(os.Stdin)
	exp, _ := reader.ReadString('\n')
	t.tokenize(strings.TrimSpace(exp)) //trim any trailing whitespace
	fmt.Println(t.printValue())
}
