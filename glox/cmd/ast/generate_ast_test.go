// Should we move this under this repo's /cmd dir?
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNaiveExprPrinter(t *testing.T) {
	tests := []struct {
		name    string
		fixture string
		doTest  func(g *generator)
	}{
		{
			name:    "test header",
			fixture: "header.txt",
			doTest: func(g *generator) {
				g.writeHeader()
			},
		},
		{
			name:    "test type",
			fixture: "type.txt",
			doTest: func(g *generator) {
				g.writeTypes([]string{
					"Binary : left Expr, operator Token, right Expr",
					"Grouping : expression Expr",
				})
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// given:
			fileContent, err := ioutil.ReadFile(fmt.Sprintf("fixtures/%s", tc.fixture))
			if err != nil {
				log.Fatal(err)
			}
			expected := string(fileContent)

			// when:
			g := generator{}
			tc.doTest(&g)
			err = g.format()
			if err != nil {
				t.Errorf("%v has an unexpected err:\nerror:\n%v\n", tc.name, err)
				return
			}

			// then:
			require.Equal(t, string(g.buf.Bytes()), expected)
		})
	}
}
