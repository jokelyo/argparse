package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/jokelyo/argparse"
)

func main() {
	// Create new parser object
	parser := argparse.NewParser("print", "Prints provided string to stdout")
	// Create string flags
	s := parser.String("s", "string", &argparse.Options{Help: "Silently ignores nargs N value", Nargs: 3})
	s2 := parser.String("", "string2", &argparse.Options{Help: "Requires 0 or 1 arguments, Default value not set", Nargs: "?"})
	s3 := parser.Strings("", "strings", &argparse.Options{Help: "One or more arguments", Nargs: "+"})
	s4 := parser.Strings("", "strings2", &argparse.Options{Help: "Zero or more arguments", Nargs: "*"})
	s5 := parser.Strings("", "strings3", &argparse.Options{Help: "Requires 3 arguments", Nargs: 3})
	// Create int flags
	i := parser.Int("i", "int", &argparse.Options{Help: "Silently ignores nargs N value", Nargs: 3})
	i2 := parser.Int("", "int2", &argparse.Options{Help: "Requires 0 or 1 arguments, Default value set", Nargs: "?", Default: 5})
	i3 := parser.Ints("", "ints", &argparse.Options{Help: "One or more arguments", Nargs: "+"})
	i4 := parser.Ints("", "ints2", &argparse.Options{Help: "Zero or more arguments", Nargs: "*"})
	i5 := parser.Ints("", "ints3", &argparse.Options{Help: "Requires 2 arguments", Nargs: 2})
	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		return
	}
	// Finally print the collected values
	args := map[string]interface{}{
		"--string":   *s,
		"--string2":  *s2,
		"--strings":  *s3,
		"--strings2": *s4,
		"--strings3": *s5,
		"--int":      *i,
		"--int2":     *i2,
		"--ints":     *i3,
		"--ints2":    *i4,
		"--ints3":    *i5,
	}
	keys := make([]string, 0, len(args))
	for k := range args {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s: %v\n", k, args[k])
	}
}
