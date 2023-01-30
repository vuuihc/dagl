package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLexer(t *testing.T) {
	file, err := os.Open("./test.gf")
	require.NoError(t, err)
	buf, err := ioutil.ReadAll(file)
	require.NoError(t, err)
	l := newLexer(string(buf))
	for {
		token, value := l.Next()
		if token == EOF {
			break
		}
		log.Printf("%s: %v\n", token, value)
	}

}
