package lox

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/modulitos/glox/pkg/interpreter"
	"github.com/stretchr/testify/assert"
)

func TestInterpreterIntegration(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected string
		regex    *regexp.Regexp
	}{
		{
			name: "nested globals",
			source: `
var a = "global a";
var b = "global b";
var c = "global c";
{
  var a = "outer a";
  var b = "outer b";
  {
    var a = "inner a";
    print a;
    print b;
    print c;
  }
  print a;
  print b;
  print c;
}
print a;
print b;
print c;
		`,
			expected: `inner a
outer b
global c
outer a
outer b
global c
global a
global b
global c
`,
		},
		{
			name: "if stmt",
			source: `
if (false)
  print "ok";
else
  print "not ok";
		`,
			expected: "not ok\n",
		},
		{
			name:     "or stmt",
			source:   `print (false or "qwer");`,
			expected: "qwer\n",
		},
		{
			name:     "and stmt",
			source:   `print ("qwer" and "foo");`,
			expected: "foo\n",
		},
		{
			name: "while stmt",
			source: `
var x = 1;
while (x < 3) {
	print "x: " + x;
	x = x + 1;
}
			`,
			expected: "x: 1\nx: 2\n",
		},
		{
			name: "c-style for loops",
			source: `
var x = 0;
var temp;
for (var y = 1; y < 10; y = temp + y) {
	print y;
	temp = x;
	x = y;
}
			`,
			expected: "1\n1\n2\n3\n5\n8\n",
		},
		{
			name:   "native function",
			source: `print clock();`,
			regex:  regexp.MustCompile(`^\d+\.\d+\n$`),
		},
		{
			name: "user-defined function",
			source: `
fun sayHi(first, last) {
  print "Hi, " + first + " " + last + "!";
}

sayHi("Dear", "Reader");
`,
			expected: "Hi, Dear Reader!\n",
		},
		{
			name: "return statement",
			source: `
fun fibx(n) {
	if (n <= 1) return n;
	return fibx(n - 2) + fibx(n - 1);
}

for (var i = 0; i < 12; i = i + 1) {
  print fibx(i);
}
`,
			expected: "0\n1\n1\n2\n3\n5\n8\n13\n21\n34\n55\n89\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Given:
			buf := new(bytes.Buffer)
			interpreter := interpreter.NewInterpreter(buf)

			// When:
			err := run([]byte(tc.source), interpreter)
			if err != nil {
				t.Errorf("%v has an unexpected err:\nerror:\n%v\n", tc.name, err)

				return
			}

			b, err := ioutil.ReadAll(buf)
			if err != nil {
				t.Errorf("%v has an unexpected err:\nerror:\n%v\n", tc.name, err)
				return
			}
			actual := string(b)

			// Then:
			if len(tc.expected) != 0 {
				assert.Equal(t, tc.expected, actual)
			} else {
				assert.Regexp(t, tc.regex, actual)
			}
		})
	}
}
