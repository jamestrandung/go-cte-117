package cte

import (
	"strings"
)

type validator interface {
	validate() error
}

type componentStack []string

func (s componentStack) push(componentName string) componentStack {
	return append(s, componentName)
}

func (s componentStack) pop() componentStack {
	return s[0 : len(s)-1]
}

func (s componentStack) clone() componentStack {
	result := make([]string, 0, len(s))

	for _, c := range s {
		result = append(result, c)
	}

	return result
}

func (s componentStack) String() string {
	return strings.Join(s, " >> ")
}
