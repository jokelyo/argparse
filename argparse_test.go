package argparse

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestFlagSimple1(t *testing.T) {
	testArgs := []string{"progname", "--flag-arg1"}

	p := NewParser("", "description")
	flag1 := p.Flag("", "flag-arg1", nil)
	flag2 := p.Flag("", "flag-arg2", nil)

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if flag1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if flag2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if *flag1 != true {
		t.Errorf("Test %s failed with flag1 being false", t.Name())
		return
	}

	if *flag2 != false {
		t.Errorf("Test %s failed with flag2 being true", t.Name())
		return
	}
}

func TestFlagSimple2(t *testing.T) {
	testArgs := []string{"progname", "--flag-arg1", "-f"}

	p := NewParser("", "description")
	flag1 := p.Flag("", "flag-arg1", nil)
	flag2 := p.Flag("", "flag-arg2", nil)
	flag3 := p.Flag("f", "flag-arg3", nil)

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if flag1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if flag2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if flag3 == nil {
		t.Errorf("Test %s failed with flag5 being nil pointer", t.Name())
		return
	}

	if *flag1 != true {
		t.Errorf("Test %s failed with flag1 being false", t.Name())
		return
	}

	if *flag2 != false {
		t.Errorf("Test %s failed with flag2 being true", t.Name())
		return
	}

	if *flag3 != true {
		t.Errorf("Test %s failed with flag3 being false", t.Name())
		return
	}
}

func TestFlagMultiShorthand1(t *testing.T) {
	testArgs := []string{"progname", "-abcd", "-e"}

	p := NewParser("", "description")
	flag1 := p.Flag("a", "aa", nil)
	flag2 := p.Flag("b", "bb", nil)
	flag3 := p.Flag("c", "cc", nil)
	flag4 := p.Flag("d", "dd", nil)
	flag5 := p.Flag("e", "ee", nil)
	flag6 := p.Flag("f", "ff", nil)

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if *flag1 != true {
		t.Errorf("Test %s failed with flag1 being false", t.Name())
	}

	if *flag2 != true {
		t.Errorf("Test %s failed with flag2 being false", t.Name())
	}

	if *flag3 != true {
		t.Errorf("Test %s failed with flag3 being false", t.Name())
	}

	if *flag4 != true {
		t.Errorf("Test %s failed with flag4 being false", t.Name())
	}

	if *flag5 != true {
		t.Errorf("Test %s failed with flag5 being false", t.Name())
	}

	if *flag6 != false {
		t.Errorf("Test %s failed with flag6 being true", t.Name())
	}
}

func TestFailDuplicate(t *testing.T) {
	testArgs := []string{"progname", "--flag-arg1", "-f"}

	p := NewParser("", "description")
	_ = p.Flag("f", "flag-arg1", nil)
	_ = p.Flag("", "flag-arg2", nil)

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed with. Duplicate flag use not detected", t.Name())
		return
	}
}

func TestFailCaseSensitive(t *testing.T) {
	testArgs := []string{"progname", "-F"}

	p := NewParser("", "description")
	_ = p.Flag("f", "", &Options{Required: true})

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed with. Sees -F as -f", t.Name())
		return
	}
}

func TestFailExcessiveArguments(t *testing.T) {
	testArgs := []string{"progname", "--flag-arg1", "whatever"}

	p := NewParser("", "description")
	_ = p.Flag("f", "flag-arg1", nil)
	_ = p.Flag("", "flag-arg2", nil)

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed with. Excessive argument not detected", t.Name())
		return
	}
}

func TestStringSimple1(t *testing.T) {
	testArgs := []string{"progname", "--flag-arg1", "test"}

	p := NewParser("", "description")
	s1 := p.String("f", "flag-arg1", nil)
	s2 := p.String("", "flag-arg2", nil)

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if s1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if s2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if *s1 != "test" {
		t.Errorf("Test %s failed. Want: [%s], got: [%s]", t.Name(), "test", *s1)
		return
	}

	if *s2 != "" {
		t.Errorf("Test %s failed. Want: [%s], got: [%s]", t.Name(), "\"\"", *s1)
		return
	}
}

func TestStringSimple2(t *testing.T) {
	testArgs := []string{"progname", "--flag-arg1", "test"}

	p := NewParser("", "description")
	s1 := p.String("f", "flag-arg1", &Options{Nargs: 3})
	s2 := p.String("", "flag-arg2", nil)

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if s1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if s2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if *s1 != "test" {
		t.Errorf("Test %s failed. Want: [%s], got: [%s]", t.Name(), "test", *s1)
		return
	}

	if *s2 != "" {
		t.Errorf("Test %s failed. Want: [%s], got: [%s]", t.Name(), "\"\"", *s1)
		return
	}
}

func TestStringsSimple1(t *testing.T) {
	testArgs := []string{"progname", "--flag-arg1", "test"}

	p := NewParser("", "description")
	s1 := p.Strings("f", "flag-arg1", nil)
	s2 := p.Strings("", "flag-arg2", nil)

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if s1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if s2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if (*s1)[0] != "test" {
		t.Errorf("Test %s failed. Want: [%s], got: [%s]", t.Name(), "test", *s1)
		return
	}

	if len(*s2) != 0 {
		t.Errorf("Test %s failed. Want: [0] length, got: [%d]", t.Name(), len(*s2))
		return
	}
}

func TestStringsNargsInvalidChar(t *testing.T) {
	testS1 := []string{"test", "test", "test"}
	testS2 := []string{"test"}
	testArgs := []string{"progname", "-f", testS1[0], testS1[1], testS1[2], "--flag-arg2", testS2[0]}

	p := NewParser("", "description")
	_ = p.Strings("f", "flag-arg1", &Options{Nargs: 3})
	_ = p.Strings("g", "flag-arg2", &Options{Nargs: "b"})

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed, expected \"invalid string value ...\" error for --flag-arg2", t.Name())
		return
	}

	if err != nil && !strings.Contains(err.Error(), "invalid string value") {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}
}

// Strings nargs N test cases:
// 	* correct number of args
// 	* incorrect number of args
func TestStringsNargsInt(t *testing.T) {
	testS1 := []string{"test", "test", "test"}
	testS2 := []string{"test"}
	testArgs := []string{"progname", "-f", testS1[0], testS1[1], testS1[2], "--flag-arg2", testS2[0], "-v"}

	p := NewParser("", "description")
	s1 := p.Strings("f", "flag-arg1", &Options{Nargs: 3})
	s2 := p.Strings("g", "flag-arg2", &Options{Nargs: 2})
	_ = p.Strings("v", "version", nil)

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed, expected \"not enough arguments ...\" error for --flag-arg2", t.Name())
		return
	}

	if err != nil && !strings.Contains(err.Error(), "not enough arguments") {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if s1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if s2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if !(len(*s1) == len(testS1) && (*s1)[0] == "test") {
		t.Errorf("Test %s s1 failed. Want: %s, got: %s", t.Name(), testS1, *s1)
		return
	}

	if len(*s2) != 0 {
		t.Errorf("Test %s s2 length failed. Want: [%d], got: [%d]", t.Name(), 0, len(*s2))
		return
	}
}

func TestStringsNargsIntWithErr(t *testing.T) {
	testS1 := []string{"test", "test", "test"}
	testS2 := []string{"test"}
	testArgs := []string{"progname", "-f", testS1[0], testS1[1], testS1[2], "--flag-arg2", testS2[0]}

	p := NewParser("", "description")
	_ = p.Strings("f", "flag-arg1", &Options{Nargs: 3})
	_ = p.Strings("g", "flag-arg2", &Options{Nargs: -1})

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed, expected \"nargs integer value ...\" error for --flag-arg2", t.Name())
		return
	}

	if err != nil && !strings.Contains(err.Error(), "nargs integer value") {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}
}

// if the args are given on the CLI in opposite order of how they're declared,
// the nargs operators would try to grab empty strings left behind by reduce()
func TestStringsNargsIntParseOrder(t *testing.T) {
	testArgs := []string{"progname", "-g", "test", "-f", "test"}

	p := NewParser("", "description")
	_ = p.String("f", "flag-arg1", nil)
	_ = p.Strings("g", "flag-arg2", &Options{Nargs: 2})

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed, expected \"not enough arguments ...\" error for --flag-arg2", t.Name())
		return
	}

	if err != nil && !strings.Contains(err.Error(), "not enough arguments") {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}
}

// String nargs '?' test cases:
//	* one arg
//	* no args
//	* followed by a flag
//	* at beginning of arg list
//	* at end of arg list
func TestStringNargs0or1(t *testing.T) {
	testS2 := []string{"test", "test"}
	testArgs := []string{"progname", "-f", "test1", "--flag-arg2", testS2[0], testS2[1], "-i", "-v", "-j"}

	p := NewParser("", "description")
	s1 := p.String("f", "flag-arg1", &Options{Nargs: "?"})
	s2 := p.Strings("g", "flag-arg2", &Options{Nargs: 2})
	s3 := p.String("i", "flag-arg3", &Options{Nargs: "?", Default: "test3"})
	s4 := p.String("j", "flag-arg4", &Options{Nargs: "?", Default: "test4"})
	_ = p.Flag("v", "version", nil)

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if s1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if s2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if s3 == nil {
		t.Errorf("Test %s failed with flag3 being nil pointer", t.Name())
		return
	}

	if s4 == nil {
		t.Errorf("Test %s failed with flag4 being nil pointer", t.Name())
		return
	}

	if *s1 != "test1" {
		t.Errorf("Test %s s1 failed. Want: [%s], got: [%s]", t.Name(), "test1", *s1)
		return
	}

	if *s3 != "test3" {
		t.Errorf("Test %s s3 failed. Want: [%s], got: [%s]", t.Name(), "test3", *s3)
		return
	}

	if *s4 != "test4" {
		t.Errorf("Test %s s4 failed. Want: [%s], got: [%s]", t.Name(), "test4", *s4)
		return
	}
}

// if the args are given on the CLI in opposite order of how they're declared,
// the nargs operators would try to grab empty strings left behind by reduce()
func TestStringNargs0or1ParseOrder(t *testing.T) {
	testArgs := []string{"progname", "-g", "-f", "test"}

	p := NewParser("", "description")
	s1 := p.String("f", "flag-arg1", nil)
	s2 := p.String("g", "flag-arg2", &Options{Nargs: "?", Default: "test2"})

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if s1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if s2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if *s1 != "test" {
		t.Errorf("Test %s s1 failed. Want: [test], got: [%s]", t.Name(), *s1)
		return
	}

	if *s2 != "test2" {
		t.Errorf("Test %s s2 failed. Want: [%s], got: [%s]", t.Name(), "test2", *s2)
		return
	}
}

// Strings nargs '*' test cases:
//	* no args
//	* multiple args
//	* at beginning of arg list
//	* at end of arg list
//	* multiple flags with '*'
func TestStringsNargs0orMore(t *testing.T) {
	testS1 := []string{"test"}
	testS3 := []string{"test", "test"}
	testArgs := []string{"progname", "-f", testS1[0], "-g", "-i", testS3[0], testS3[1], "-j"}

	p := NewParser("", "description")
	s1 := p.Strings("f", "flag-arg1", &Options{Nargs: "*"})
	s2 := p.Strings("g", "flag-arg2", &Options{Nargs: "*"})
	s3 := p.Strings("i", "flag-arg3", &Options{Nargs: "*"})
	s4 := p.Strings("j", "flag-arg4", &Options{Nargs: "*"})

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if s1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if s2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if s3 == nil {
		t.Errorf("Test %s failed with flag3 being nil pointer", t.Name())
		return
	}

	if s4 == nil {
		t.Errorf("Test %s failed with flag4 being nil pointer", t.Name())
		return
	}

	if !(len(*s1) == len(testS1) && (*s1)[0] == "test") {
		t.Errorf("Test %s s1 failed. Want: %s, got: %s", t.Name(), testS1, *s1)
		return
	}

	if len(*s2) != 0 {
		t.Errorf("Test %s s2 length failed. Want: [0], got: [%d]", t.Name(), len(*s2))
		return
	}

	if !(len(*s3) == len(testS3) && (*s3)[0] == "test") {
		t.Errorf("Test %s s3 failed. Want: %s, got: %s", t.Name(), testS3, *s3)
		return
	}

	if len(*s4) != 0 {
		t.Errorf("Test %s s4 length failed. Want: [0], got: [%d]", t.Name(), len(*s4))
		return
	}
}

// if the args are given on the CLI in opposite order of how they're declared,
// the nargs operators would try to grab empty strings left behind by reduce()
func TestStringsNargs0orMoreParseOrder(t *testing.T) {
	testS2 := []string{"test", "test"}
	testArgs := []string{"progname", "-g", testS2[0], testS2[1], "-f", "test"}

	p := NewParser("", "description")
	s1 := p.String("f", "flag-arg1", nil)
	s2 := p.Strings("g", "flag-arg2", &Options{Nargs: "*"})

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if s1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if s2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if !(len(*s2) == len(testS2) && (*s2)[0] == "test") {
		t.Errorf("Test %s s1 failed. Want: %s, got: %s", t.Name(), testS2, *s2)
		return
	}

	if *s1 != "test" {
		t.Errorf("Test %s s2 failed. Want: [test], got: [%s]", t.Name(), *s2)
		return
	}
}

// Strings nargs '+' test cases:
//	* one arg
//	* multiple args
//	* at beginning of arg list
//	* at end of arg list
//	* multiple flags with '+'
func TestStringsNargs1orMore(t *testing.T) {
	testS1 := []string{"test"}
	testS2 := []string{"test", "test"}
	testArgs := []string{"progname", "-f", testS1[0], "-g", testS2[0], testS2[1]}

	p := NewParser("", "description")
	s1 := p.Strings("f", "flag-arg1", &Options{Nargs: "+"})
	s2 := p.Strings("g", "flag-arg2", &Options{Nargs: "+"})

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if s1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if s2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if !(len(*s1) == len(testS1) && (*s1)[0] == "test") {
		t.Errorf("Test %s s1 failed. Want: %s, got: %s", t.Name(), testS1, *s1)
		return
	}

	if !(len(*s2) == len(testS2) && (*s2)[0] == "test") {
		t.Errorf("Test %s s2 failed. Want: %s, got: %s", t.Name(), testS2, *s2)
		return
	}
}

func TestStringsNargs1orMoreWithErr(t *testing.T) {
	testS2 := []string{"test", "test"}
	testArgs := []string{"progname", "-f", "-g", testS2[0], testS2[1]}

	p := NewParser("", "description")
	_ = p.Strings("f", "flag-arg1", &Options{Nargs: "+"})
	_ = p.Strings("i", "flag-arg2", &Options{Nargs: "+"})

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed, expected \"... at least one argument\" error for --flag-arg1", t.Name())
		return
	}

	if err != nil && !strings.Contains(err.Error(), "requires at least one argument") {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}
}

// if the args are given on the CLI in opposite order of how they're declared,
// the nargs operators would try to grab empty strings left behind by reduce()
func TestStringsNargs1orMoreParseOrder(t *testing.T) {
	testS2 := []string{"test", "test"}
	testArgs := []string{"progname", "-g", testS2[0], testS2[1], "-f", "test"}

	p := NewParser("", "description")
	s1 := p.String("f", "flag-arg1", nil)
	s2 := p.Strings("g", "flag-arg2", &Options{Nargs: "*"})

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if s1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if s2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if !(len(*s2) == len(testS2) && (*s2)[0] == "test") {
		t.Errorf("Test %s s1 failed. Want: %s, got: %s", t.Name(), testS2, *s2)
		return
	}

	if *s1 != "test" {
		t.Errorf("Test %s s2 failed. Want: [test], got: [%s]", t.Name(), *s2)
		return
	}
}

func TestIntSimple1(t *testing.T) {
	val := 5150
	testArgs := []string{"progname", "--flag-arg1", strconv.Itoa(val)}

	p := NewParser("", "description")
	i1 := p.Int("f", "flag-arg1", nil)
	i2 := p.Int("", "flag-arg2", nil)

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if i1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if i2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if *i1 != val {
		t.Errorf("Test %s failed. Want: [%d], got: [%d]", t.Name(), val, *i1)
		return
	}

	if *i2 != 0 {
		t.Errorf("Test %s failed. Want: [%d], got: [%d]", t.Name(), 0, *i1)
		return
	}
}

func TestIntSimple2(t *testing.T) {
	val := 5150
	testArgs := []string{"progname", "--flag-arg1", strconv.Itoa(val)}

	p := NewParser("", "description")
	i1 := p.Int("f", "flag-arg1", &Options{Nargs: 3})
	i2 := p.Int("", "flag-arg2", nil)

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if i1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if i2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if *i1 != val {
		t.Errorf("Test %s failed. Want: [%d], got: [%d]", t.Name(), val, *i1)
		return
	}

	if *i2 != 0 {
		t.Errorf("Test %s failed. Want: [%d], got: [%d]", t.Name(), 0, *i1)
		return
	}
}

func TestIntFailSimple1(t *testing.T) {
	testArgs := []string{"progname", "--flag-arg1", "string"}

	p := NewParser("", "description")
	i1 := p.Int("f", "flag-arg1", nil)

	err := p.Parse(testArgs)
	errStr := "[-f|--flag-arg1] bad interger value [string]"
	if err == nil || err.Error() != errStr {
		t.Errorf("Test %s expected [%s], got [%+v]", t.Name(), errStr, err)
		return
	}

	if i1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if *i1 != 0 {
		t.Errorf("Test %s failed. Want: [0], got: [%d]", t.Name(), *i1)
		return
	}
}

func TestIntsSimple1(t *testing.T) {
	testval := 5150
	testArgs := []string{"progname", "--flag-arg1", strconv.Itoa(testval)}

	p := NewParser("", "description")
	i1 := p.Ints("f", "flag-arg1", nil)
	i2 := p.Ints("", "flag-arg2", nil)

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if i1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if i2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if len(*i1) != 1 {
		t.Errorf("Test %s length failed. Want: [1], got: [%d]", t.Name(), len(*i1))
		return
	}

	if (*i1)[0] != testval {
		t.Errorf("Test %s value failed. Want: [%d], got: [%d]", t.Name(), testval, (*i1)[0])
		return
	}

	if len(*i2) != 0 {
		t.Errorf("Test %s failed. Want: [0] length, got: [%d]", t.Name(), len(*i2))
		return
	}
}

func TestIntsNargsInvalidChar(t *testing.T) {
	testArgs := []string{"progname", "-f", "1", "2", "3"}

	p := NewParser("", "description")
	_ = p.Ints("f", "flag-arg1", &Options{Nargs: "x"})

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed, expected \"invalid string value ...\" error for --flag-arg1", t.Name())
		return
	}

	if err != nil && !strings.Contains(err.Error(), "invalid string value") {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}
}

func TestIntsNargsInvalidNum(t *testing.T) {
	testArgs := []string{"progname", "-f", "1", "2", "3"}

	p := NewParser("", "description")
	_ = p.Ints("f", "flag-arg1", &Options{Nargs: -1})

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed, expected \"nargs integer value ...\" error for --flag-arg1", t.Name())
		return
	}

	if err != nil && !strings.Contains(err.Error(), "nargs integer value") {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}
}

// Strings nargs N test cases:
// 	* correct number of args
// 	* incorrect number of args
func TestIntsNargsInt(t *testing.T) {
	testI1 := []int{1, 2, 3}
	testI2 := []int{5}
	testArgs := []string{
		"progname",
		"-f",
		strconv.Itoa(testI1[0]),
		strconv.Itoa(testI1[1]),
		strconv.Itoa(testI1[2]),
		"--flag-arg2",
		strconv.Itoa(testI2[0]),
		"-v",
	}

	p := NewParser("", "description")
	i1 := p.Ints("f", "flag-arg1", &Options{Nargs: 3})
	i2 := p.Ints("g", "flag-arg2", &Options{Nargs: 2})
	_ = p.Ints("v", "version", nil)

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed, expected \"not enough arguments ...\" error for --flag-arg2", t.Name())
		return
	}

	if err != nil && !strings.Contains(err.Error(), "not enough arguments") {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if i1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if i2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if !(len(*i1) == len(testI1) && (*i1)[0] == testI1[0]) {
		t.Errorf("Test %s i1 failed. Want: %v, got: %v", t.Name(), testI1, *i1)
		return
	}

	if len(*i2) != 0 {
		t.Errorf("Test %s i2 length failed. Want: [%d], got: [%d]", t.Name(), 0, len(*i2))
		return
	}
}

func TestIntsNargsIntWithErr(t *testing.T) {
	testI1 := []int{1, 2, 3}
	testI2 := []int{5}
	testArgs := []string{
		"progname",
		"-f",
		strconv.Itoa(testI1[0]),
		strconv.Itoa(testI1[1]),
		strconv.Itoa(testI1[2]),
		"--flag-arg2",
		strconv.Itoa(testI2[0]),
	}
	p := NewParser("", "description")
	_ = p.Ints("f", "flag-arg1", &Options{Nargs: 3})
	_ = p.Ints("g", "flag-arg2", &Options{Nargs: -1})

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed, expected \"nargs integer value ...\" error for --flag-arg2", t.Name())
		return
	}

	if err != nil && !strings.Contains(err.Error(), "nargs integer value") {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}
}

// String nargs '?' test cases:
//	* one arg
//	* no args
//	* followed by a flag
//	* at beginning of arg list
//	* at end of arg list
func TestIntNargs0or1(t *testing.T) {
	testval := 5150
	testI2 := []int{5, 6}
	testArgs := []string{
		"progname",
		"-f",
		strconv.Itoa(testval),
		"--flag-arg2",
		strconv.Itoa(testI2[0]),
		strconv.Itoa(testI2[1]),
		"-i",
		"-v",
		"-j",
	}

	p := NewParser("", "description")
	i1 := p.Int("f", "flag-arg1", &Options{Nargs: "?"})
	i2 := p.Ints("g", "flag-arg2", &Options{Nargs: 2})
	i3 := p.Int("i", "flag-arg3", &Options{Nargs: "?", Default: 3})
	i4 := p.Int("j", "flag-arg4", &Options{Nargs: "?", Default: 4})
	_ = p.Flag("v", "version", nil)

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if i1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if i2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if i3 == nil {
		t.Errorf("Test %s failed with flag3 being nil pointer", t.Name())
		return
	}

	if i4 == nil {
		t.Errorf("Test %s failed with flag4 being nil pointer", t.Name())
		return
	}

	if *i1 != testval {
		t.Errorf("Test %s i1 failed. Want: [%d], got: [%d]", t.Name(), testval, *i1)
		return
	}

	if *i3 != 3 {
		t.Errorf("Test %s i3 failed. Want: [%d], got: [%d]", t.Name(), 3, *i3)
		return
	}

	if *i4 != 4 {
		t.Errorf("Test %s i4 failed. Want: [%d], got: [%d]", t.Name(), 4, *i4)
		return
	}
}

// Strings nargs '*' test cases:
//	* no args
//	* multiple args
//	* at beginning of arg list
//	* at end of arg list
//	* multiple flags with '*'
func TestIntsNargs0orMore(t *testing.T) {
	testI1 := []int{1}
	testI3 := []int{3, 4}
	testArgs := []string{
		"progname",
		"-f",
		strconv.Itoa(testI1[0]),
		"-g",
		"-i",
		strconv.Itoa(testI3[0]),
		strconv.Itoa(testI3[1]),
		"-j",
	}

	p := NewParser("", "description")
	i1 := p.Ints("f", "flag-arg1", &Options{Nargs: "*"})
	i2 := p.Ints("g", "flag-arg2", &Options{Nargs: "*"})
	i3 := p.Ints("i", "flag-arg3", &Options{Nargs: "*"})
	i4 := p.Ints("j", "flag-arg4", &Options{Nargs: "*"})

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if i1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if i2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if i3 == nil {
		t.Errorf("Test %s failed with flag3 being nil pointer", t.Name())
		return
	}

	if i4 == nil {
		t.Errorf("Test %s failed with flag4 being nil pointer", t.Name())
		return
	}

	if !(len(*i1) == len(testI1) && (*i1)[0] == testI1[0]) {
		t.Errorf("Test %s i1 failed. Want: %d, got: %d", t.Name(), testI1, *i1)
		return
	}

	if len(*i2) != 0 {
		t.Errorf("Test %s i2 length failed. Want: [0], got: [%d]", t.Name(), len(*i2))
		return
	}

	if !(len(*i3) == len(testI3) && (*i3)[0] == testI3[0]) {
		t.Errorf("Test %s i3 failed. Want: %v, got: %v", t.Name(), testI3, *i3)
		return
	}

	if len(*i4) != 0 {
		t.Errorf("Test %s i4 length failed. Want: [0], got: [%d]", t.Name(), len(*i4))
		return
	}
}

// Strings nargs '+' test cases:
//	* one arg
//	* multiple args
//	* at beginning of arg list
//	* at end of arg list
//	* multiple flags with '+'
func TestIntsNargs1orMore(t *testing.T) {
	testI1 := []int{1}
	testI2 := []int{3, 4}
	testArgs := []string{
		"progname",
		"-f",
		strconv.Itoa(testI1[0]),
		"-g",
		strconv.Itoa(testI2[0]),
		strconv.Itoa(testI2[1]),
	}

	p := NewParser("", "description")
	i1 := p.Ints("f", "flag-arg1", &Options{Nargs: "+"})
	i2 := p.Ints("g", "flag-arg2", &Options{Nargs: "+"})

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if i1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if i2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if !(len(*i1) == len(testI1) && (*i1)[0] == testI1[0]) {
		t.Errorf("Test %s i1 failed. Want: %v, got: %v", t.Name(), testI1, *i1)
		return
	}

	if !(len(*i2) == len(testI2) && (*i2)[0] == testI2[0]) {
		t.Errorf("Test %s i2 failed. Want: %v, got: %v", t.Name(), testI2, *i2)
		return
	}
}

func TestIntsNargs1orMoreWithErr(t *testing.T) {
	testI2 := []int{3, 4}
	testArgs := []string{"progname", "-f", "-g", strconv.Itoa(testI2[0]), strconv.Itoa(testI2[1])}

	p := NewParser("", "description")
	_ = p.Ints("f", "flag-arg1", &Options{Nargs: "+"})
	_ = p.Ints("i", "flag-arg2", &Options{Nargs: "+"})

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed, expected \"... at least one argument\" error for --flag-arg1", t.Name())
		return
	}

	if err != nil && !strings.Contains(err.Error(), "requires at least one argument") {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}
}

func TestFileSimple1(t *testing.T) {
	// Test file location
	fpath := "./test.tmp"
	// Create test file
	f, err := os.Create(fpath)
	if err != nil {
		t.Error(err)
		return
	}
	f.Close()
	defer os.Remove(fpath)

	testArgs := []string{"progname", "-f", fpath}

	p := NewParser("", "")

	file1 := p.File("f", "file", os.O_RDWR, 0666, &Options{Default: "./non-existent-file.tmp"})

	err = p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}
	defer file1.Close()

	if file1 == nil {
		t.Errorf("Test %s failed with file1 being nil pointer", t.Name())
		return
	}

	testString := "Test"
	recSlice := make([]byte, 4)
	_, err = file1.WriteString(testString)
	if err != nil {
		t.Errorf("Test %s write operation failed with error: %s", t.Name(), err.Error())
		return
	}
	file1.Seek(0, 0)
	n, err := file1.Read(recSlice)
	if err != nil {
		t.Errorf("Test %s read operation failed with error: %s", t.Name(), err.Error())
		return
	}
	if n != 4 || string(recSlice) != testString {
		t.Errorf("Test %s failed on read operation", t.Name())
		return
	}
}

func TestListSimple1(t *testing.T) {
	testArgs := []string{"progname", "--flag-arg1", "test1", "--flag-arg1", "test2"}
	list1Expect := []string{"test1", "test2"}
	list2Expect := make([]string, 0)

	p := NewParser("", "description")
	l1 := p.List("f", "flag-arg1", nil)
	l2 := p.List("", "flag-arg2", nil)

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if l1 == nil {
		t.Errorf("Test %s failed with l1 being nil pointer", t.Name())
		return
	}

	if l2 == nil {
		t.Errorf("Test %s failed with l2 being nil pointer", t.Name())
		return
	}

	if !reflect.DeepEqual(*l1, list1Expect) {
		t.Errorf("Test %s failed. Want: %s, got: %s", t.Name(), list1Expect, *l1)
		return
	}

	if !reflect.DeepEqual(*l2, list2Expect) {
		t.Errorf("Test %s failed. Want: %s, got: %s", t.Name(), list2Expect, *l2)
		return
	}
}

func TestSelectorSimple1(t *testing.T) {
	flag1Expect := "test2"
	allowedValues := []string{"test1", flag1Expect}
	testArgs := []string{"progname", "--flag-arg1", flag1Expect}

	p := NewParser("", "")
	s1 := p.Selector("f", "flag-arg1", allowedValues, nil)

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if s1 == nil {
		t.Errorf("Test %s failed with s1 being nil pointer", t.Name())
		return
	}

	if *s1 != flag1Expect {
		t.Errorf("Test %s failed. Want: %s, got: %s", t.Name(), flag1Expect, *s1)
		return
	}
}

func TestSelectorFailSimple1(t *testing.T) {
	allowedValues := []string{"test1", "test2"}
	testArgs := []string{"progname", "--flag-arg1", "test3"}

	p := NewParser("", "")
	_ = p.Selector("f", "flag-arg1", allowedValues, nil)

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed. Expected error did not happen", t.Name())
		return
	}
}

func TestCommandSimple1(t *testing.T) {
	val := 5150
	testArgsList := [][]string{
		{"progname", "cmd1", "--flag1", "--string-flag1", "test", "--int-flag1", strconv.Itoa(val)},
		{"progname", "cmd2"},
	}

	for _, testArgs := range testArgsList {
		p := NewParser("progname", "description")

		cmd1 := p.NewCommand("cmd1", "cmd1 description")
		flag1 := cmd1.Flag("f", "flag1", nil)
		string1 := cmd1.String("s", "string-flag1", nil)
		int1 := cmd1.Int("i", "int-flag1", nil)

		cmd2 := p.NewCommand("cmd2", "cmd2 description")

		p.Parse(testArgs)

		if cmd1.Happened() {
			if *flag1 != true {
				t.Errorf("Test %s failed with %s: flag1: wanted [true], got [false]", t.Name(), testArgs[1])
				return
			}
			if *string1 != "test" {
				t.Errorf("Test %s failed with %s: string1: wanted [test], got [%s]", t.Name(), testArgs[1], *string1)
				return
			}
			if *int1 != val {
				t.Errorf("Test %s failed with %s: int1: wanted [%d], got [%d]", t.Name(), testArgs[1], val, *int1)
				return
			}
		}
		if cmd2.Happened() {
			if *flag1 != false {
				t.Errorf("Test %s failed with %s: flag1: wanted [false], got [true]", t.Name(), testArgs[1])
				return
			}
			if *string1 != "" {
				t.Errorf("Test %s failed with %s: string1: wanted [], got [%s]", t.Name(), testArgs[1], *string1)
				return
			}
			if *int1 != 0 {
				t.Errorf("Test %s failed with %s: int1: wanted [0], got [%d]", t.Name(), testArgs[1], *int1)
				return
			}
		}
		if (cmd1.Happened() && cmd2.Happened()) || (!cmd1.Happened() && !cmd2.Happened()) {
			t.Errorf("Test %s failed, either cmd1 and cmd2 or neither of them Happened()", t.Name())
			return
		}
	}
}

func TestCommandMixedArgs1(t *testing.T) {
	val := 5150
	pval := 316
	testArgsList := [][]string{
		{"progname", "cmd1", "--flag1", "--string-flag1", "test", "--int-flag1", strconv.Itoa(val), "--global-flag", "--global-string", "global test string", "--global-int", strconv.Itoa(pval)},
		{"progname", "cmd2", "--global-string", "global test string", "--global-flag", "--global-int", strconv.Itoa(pval)},
	}

	for _, testArgs := range testArgsList {
		p := NewParser("progname", "description")

		cmd1 := p.NewCommand("cmd1", "cmd1 description")
		cmd1flag1 := cmd1.Flag("f", "flag1", nil)
		cmd1string1 := cmd1.String("s", "string-flag1", nil)
		cmd1int1 := cmd1.Int("i", "int-flag1", nil)

		cmd2 := p.NewCommand("cmd2", "cmd2 description")

		pflag1 := p.Flag("", "global-flag", nil)
		pstring1 := p.String("", "global-string", nil)
		pint1 := p.Int("", "global-int", nil)

		p.Parse(testArgs)

		// Check global flags
		if *pflag1 != true {
			t.Errorf("Test %s failed with %s: pflag1: wanted [true], got [false]", t.Name(), testArgs[1])
			return
		}
		if *pstring1 != "global test string" {
			t.Errorf("Test %s failed with %s: pstring1: wanted [global test string], got [%s]", t.Name(), testArgs[1], *pstring1)
			return
		}
		if *pint1 != pval {
			t.Errorf("Test %s failed with %s: pint1: wanted [%d], got [%d]", t.Name(), testArgs[1], pval, *pint1)
			return
		}

		// Check commands
		if cmd1.Happened() {
			if *cmd1flag1 != true {
				t.Errorf("Test %s failed with %s: flag1: wanted [true], got [false]", t.Name(), testArgs[1])
				return
			}
			if *cmd1string1 != "test" {
				t.Errorf("Test %s failed with %s: string1: wanted [test], got [%s]", t.Name(), testArgs[1], *cmd1string1)
				return
			}
			if *cmd1int1 != val {
				t.Errorf("Test %s failed with %s: int1: wanted [%d], got [%d]", t.Name(), testArgs[1], val, *cmd1int1)
				return
			}
		}
		if cmd2.Happened() {
			if *cmd1flag1 != false {
				t.Errorf("Test %s failed with %s: flag1: wanted [false], got [true]", t.Name(), testArgs[1])
				return
			}
			if *cmd1string1 != "" {
				t.Errorf("Test %s failed with %s: string1: wanted [], got [%s]", t.Name(), testArgs[1], *cmd1string1)
				return
			}
			if *cmd1int1 != 0 {
				t.Errorf("Test %s failed with %s: int1: wanted [0], got [%d]", t.Name(), testArgs[1], *cmd1int1)
				return
			}
		}
		if (cmd1.Happened() && cmd2.Happened()) || (!cmd1.Happened() && !cmd2.Happened()) {
			t.Errorf("Test %s failed, either cmd1 and cmd2 or neither of them Happened()", t.Name())
			return
		}
	}
}

func TestOptsRequired1(t *testing.T) {
	testArgs := []string{"progname", "--flag-arg1"}

	p := NewParser("", "description")
	_ = p.Flag("", "flag-arg1", nil)
	_ = p.String("", "flag-arg2", &Options{Required: true})

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed to detect required argument", t.Name())
		return
	}
}

func TestOptsRequired2(t *testing.T) {
	testArgs := []string{"progname", "--flag-arg1"}

	p := NewParser("", "description")
	_ = p.Flag("", "flag-arg1", nil)
	_ = p.Int("", "int-arg1", &Options{Required: true})

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed to detect required argument", t.Name())
		return
	}
}

var stropts = &Options{Validate: func(args []string) error {
	if len(args) > 0 {
		if args[0] != "pass" {
			return errors.New("failure")
		}
	}
	return nil
},
}

func TestOptsValidatePass1(t *testing.T) {
	testArgsList := [][]string{
		{"progname", "--string-flag1", "pass"},
		{"progname", "--string-flag1", "fail"},
	}

	for _, testArgs := range testArgsList {
		p := NewParser("progname", "")

		string1 := p.String("", "string-flag1", stropts)

		err := p.Parse(testArgs)

		if testArgs[2] == "pass" {
			if err != nil {
				t.Errorf("Test %s failed on %s with err: %s", t.Name(), testArgs[2], err.Error())
				return
			}

			if *string1 != "pass" {
				t.Errorf("Test %s failed on %s; string1 expected [%s], got [%s]", t.Name(), testArgs[2], testArgs[2], *string1)
				return
			}
		} else {
			if err == nil {
				t.Errorf("Test %s failed to validate argument (should return error)", t.Name())
				return
			}
		}
	}
}

func TestOptsValidatePass2(t *testing.T) {
	val1 := 5150
	val2 := 316

	var intopts = &Options{Validate: func(args []string) error {
		if len(args) > 0 {
			myval, err := strconv.Atoi(args[0])
			if err != nil {
				return errors.New("conversion failure")
			} else if myval != val1 {
				return errors.New("failure")
			}
		}
		return nil
	},
	}

	testArgsList := [][]string{
		{"progname", "--int-flag1", strconv.Itoa(val1)},
		{"progname", "--int-flag1", strconv.Itoa(val2)},
	}

	for _, testArgs := range testArgsList {
		p := NewParser("progname", "")

		int1 := p.Int("", "int-flag1", intopts)

		err := p.Parse(testArgs)

		if testArgs[2] == strconv.Itoa(val1) {
			if err != nil {
				t.Errorf("Test %s failed on %s with err: %s", t.Name(), testArgs[2], err.Error())
				return
			}

			if *int1 != val1 {
				t.Errorf("Test %s failed on %s; int1 expected [%d], got [%d]", t.Name(), testArgs[2], val1, *int1)
				return
			}
		} else {
			if err == nil {
				t.Errorf("Test %s failed to validate argument (should return error)", t.Name())
				return
			}
		}
	}
}

var pUsage = `usage: verylongprogname <Command> [-h|--help] [-s|--verylongstring-flag1
                        "<value>"] [-i|--integer-flag1 <integer>]

                        prog description

Commands:

  veryverylongcmd1  cmd1 description
  cmd2              cmd2 description

Arguments:

  -h  --help                  Print help information
  -s  --verylongstring-flag1  string1 description
  -i  --integer-flag1         integer1 description

`

var cmd1Usage = `usage: verylongprogname veryverylongcmd1 [-f|--verylongflag1]
                        -a|--verylongflagA [-h|--help]
                        [-s|--verylongstring-flag1 "<value>"]
                        [-i|--integer-flag1 <integer>]

                        cmd1 description

Arguments:

  -f  --verylongflag1         flag1 description
  -a  --verylongflagA         flag1 description
  -h  --help                  Print help information
  -s  --verylongstring-flag1  string1 description
  -i  --integer-flag1         integer1 description

`

var cmd2Usage = `usage: verylongprogname cmd2 [-h|--help] [-s|--verylongstring-flag1 "<value>"]
                        [-i|--integer-flag1 <integer>]

                        cmd2 description

Arguments:

  -h  --help                  Print help information
  -s  --verylongstring-flag1  string1 description
  -i  --integer-flag1         integer1 description

`

func TestUsageSimple1(t *testing.T) {
	p := NewParser("verylongprogname", "prog description")

	cmd1 := p.NewCommand("veryverylongcmd1", "cmd1 description")
	_ = cmd1.Flag("f", "verylongflag1", &Options{Help: "flag1 description"})
	_ = cmd1.Flag("a", "verylongflagA", &Options{Required: true, Help: "flag1 description"})
	_ = p.String("s", "verylongstring-flag1", &Options{Help: "string1 description"})
	_ = p.Int("i", "integer-flag1", &Options{Help: "integer1 description"})

	cmd2 := p.NewCommand("cmd2", "cmd2 description")

	p.Parse(os.Args)

	if pUsage != p.Usage(nil) {
		t.Errorf("%s", p.Usage(nil))
	}
	if cmd1Usage != cmd1.Usage(nil) {
		t.Errorf("%s", cmd1.Usage(nil))
	}
	if cmd2Usage != cmd2.Usage(nil) {
		t.Errorf("%s", cmd2.Usage(nil))
	}
}

func TestUsageHidden1(t *testing.T) {
	p := NewParser("verylongprogname", "prog description")

	cmd1 := p.NewCommand("veryverylongcmd1", "cmd1 description")
	_ = cmd1.Flag("f", "verylongflag1", &Options{Help: "flag1 description"})
	_ = cmd1.Flag("a", "verylongflagA", &Options{Required: true, Help: "flag1 description"})
	_ = p.String("s", "verylongstring-flag1", &Options{Help: "string1 description"})
	_ = p.Int("i", "integer-flag1", &Options{Help: "integer1 description"})
	_ = p.Int("i2", "integer-flag2", &Options{Help: DisableDescription})

	_ = p.NewCommand("cmd2", "cmd2 description")

	cmd3 := p.NewCommand("cmd3", DisableDescription)
	_ = cmd3.Flag("f", "verylongflag1", &Options{Help: "flag1 description"})
	_ = cmd3.Flag("a", "verylongflagA", &Options{Required: true, Help: "flag1 description"})

	p.Parse(os.Args)

	if pUsage != p.Usage(nil) {
		t.Errorf("%s", p.Usage(nil))
	}
	if cmd1Usage != cmd1.Usage(nil) {
		t.Errorf("%s", cmd1.Usage(nil))
	}
}

func TestStringMissingArgFail(t *testing.T) {
	testArgs := []string{"progname", "-s"}

	p := NewParser("progname", "Prog description")

	_ = p.String("s", "string", &Options{Required: true, Help: "A test string"})

	err := p.Parse(testArgs)

	if err != nil {
		// Test should pass on failure
		if err.Error() != "not enough arguments for -s|--string" {
			t.Errorf("Test %s failed: expected error [%s], got error [%s]", t.Name(), "not enough arguments for -s|--string", err.Error())
		}
	}
}

func TestIntMissingArgFail(t *testing.T) {
	testArgs := []string{"progname", "-i"}

	p := NewParser("progname", "Prog description")

	_ = p.Int("i", "integer", &Options{Required: true, Help: "A test integer"})

	err := p.Parse(testArgs)

	if err != nil {
		// Test should pass on failure
		errStr := "not enough arguments for -i|--integer"
		if err.Error() != errStr {
			t.Errorf("Test %s failed: expected error [%s], got error [%s]", t.Name(), errStr, err.Error())
		}
	}
}

func TestFlagDefaultValuePass(t *testing.T) {
	testArgs := []string{"progname"}

	p := NewParser("progname", "Prog description")

	f := p.Flag("f", "flag", &Options{Default: true})

	err := p.Parse(testArgs)

	// Should fail on failure
	if err != nil {
		t.Error(err.Error())
	}

	// Should fail if not true
	if *f != true {
		t.Errorf("expected [true], got [%t]", *f)
	}
}

func TestFlagDefaultValueFail(t *testing.T) {
	testArgs := []string{"progname"}

	p := NewParser("progname", "Prog description")

	_ = p.Flag("f", "flag", &Options{Default: "string"})

	err := p.Parse(testArgs)

	// Should pass on failure
	if err == nil || err.Error() != "cannot use default type [string] as type [bool]" {
		t.Errorf("Test %s failed: expected error [%s], got error [%+v]", t.Name(), "cannot use default type [string] as type [bool]", err)
	}
}

func TestStringDefaultValuePass(t *testing.T) {
	testArgs := []string{"progname"}
	testString := "test string"

	p := NewParser("progname", "Prog description")

	s := p.String("s", "string", &Options{Default: testString})

	err := p.Parse(testArgs)

	// Should fail on failure
	if err != nil {
		t.Error(err.Error())
	}

	// Should fail if not true
	if *s != testString {
		t.Errorf("expected [string], got [%T]", *s)
	}
}

func TestStringDefaultValueFail(t *testing.T) {
	testArgs := []string{"progname"}

	p := NewParser("progname", "Prog description")

	_ = p.String("s", "string", &Options{Default: false})

	err := p.Parse(testArgs)

	// Should pass on failure
	if err == nil || err.Error() != "cannot use default type [bool] as type [string]" {
		t.Errorf("Test %s failed: expected error [%s], got error [%+v]", t.Name(), "cannot use default type [bool] as type [string]", err)
	}
}

func TestIntDefaultValuePass(t *testing.T) {
	testArgs := []string{"progname"}
	testVal := 5150

	p := NewParser("progname", "Prog description")

	i := p.Int("i", "integer", &Options{Default: testVal})

	err := p.Parse(testArgs)

	// Should fail on failure
	if err != nil {
		t.Error(err.Error())
	}

	// Should fail if not true
	if *i != testVal {
		t.Errorf("expected [%d], got [%d]", testVal, *i)
	}
}

func TestIntDefaultValueFail(t *testing.T) {
	testArgs := []string{"progname"}

	p := NewParser("progname", "Prog description")

	_ = p.Int("i", "integer", &Options{Default: "fail"})

	err := p.Parse(testArgs)

	// Should pass on failure
	if err == nil || err.Error() != "cannot use default type [string] as type [int]" {
		t.Errorf("Test %s failed: expected error [%s], got error [%+v]", t.Name(), "cannot use default type [bool] as type [string]", err)
	}
}

func TestFileDefaultValuePass(t *testing.T) {
	// Test file location
	fpath := "./test.tmp"
	// Create test file
	f, err := os.Create(fpath)
	if err != nil {
		t.Error(err)
		return
	}
	f.Close()
	defer os.Remove(fpath)

	testArgs := []string{"progname"}

	p := NewParser("", "")

	file1 := p.File("f", "file", os.O_RDWR, 0666, &Options{Default: fpath})

	err = p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}
	defer file1.Close()
}

func TestFileDefaultValueFail(t *testing.T) {
	// Test file location
	fpath := "./test.tmp"
	// Create test file
	f, err := os.Create(fpath)
	if err != nil {
		t.Error(err)
		return
	}
	f.Close()
	defer os.Remove(fpath)

	testArgs := []string{"progname"}

	p := NewParser("", "")

	file1 := p.File("f", "file", os.O_RDWR, 0666, &Options{Default: true})

	err = p.Parse(testArgs)
	if err == nil || err.Error() != "cannot use default type [bool] as type [string]" {
		t.Errorf("Test %s failed: expected error [%s], got error [%+v]", t.Name(), "cannot use default type [bool] as type [string]", err)
	}
	defer file1.Close()
}

func TestListDefaultValuePass(t *testing.T) {
	testArgs := []string{"progname"}
	testList := []string{"test", "list"}

	p := NewParser("progname", "Prog description")

	s := p.List("s", "string", &Options{Default: testList})

	err := p.Parse(testArgs)

	// Should fail on failure
	if err != nil {
		t.Error(err.Error())
	}

	// Should fail if not true
	if !reflect.DeepEqual(*s, testList) {
		t.Errorf("expected [%v], got [%v]", testList, *s)
	}
}

func TestListDefaultValueFail(t *testing.T) {
	testArgs := []string{"progname"}

	p := NewParser("progname", "Prog description")

	_ = p.List("s", "string", &Options{Default: false})

	err := p.Parse(testArgs)

	// Should pass on failure
	if err == nil || err.Error() != "cannot use default type [bool] as type [[]string]" {
		t.Errorf("Test %s failed: expected error [%s], got error [%+v]", t.Name(), "cannot use default type [bool] as type [[]string]", err)
	}
}

func TestSelectorDefaultValuePass(t *testing.T) {
	testArgs := []string{"progname"}
	testString := "test list"

	p := NewParser("progname", "Prog description")

	s := p.Selector("s", "string", []string{"opt1", "opt2"}, &Options{Default: testString})

	err := p.Parse(testArgs)

	// Should fail on failure
	if err != nil {
		t.Error(err.Error())
	}

	// Should fail if not true
	if *s != testString {
		t.Errorf("expected [%v], got [%v]", testString, *s)
	}
}

func TestSelectorDefaultValueFail(t *testing.T) {
	testArgs := []string{"progname"}

	p := NewParser("progname", "Prog description")

	_ = p.Selector("s", "string", []string{"opt1", "opt2"}, &Options{Default: false})

	err := p.Parse(testArgs)

	// Should pass on failure
	if err == nil || err.Error() != "cannot use default type [bool] as type [string]" {
		t.Errorf("Test %s failed: expected error [%s], got error [%+v]", t.Name(), "cannot use default type [bool] as type [string]", err)
	}
}

func TestFloatSimple1(t *testing.T) {
	pi := "3.1415"
	piVal := 3.1415
	testArgs := []string{"progname", "--float1", pi}

	p := NewParser("", "description")
	f1 := p.Float("f", "float1", nil)
	f2 := p.Float("", "float2", nil)

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if f1 == nil {
		t.Errorf("Test %s failed with float1 being nil pointer", t.Name())
		return
	}

	if f2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if *f1 != piVal {
		t.Errorf("Test %s failed. Want: [%f], got: [%f]", t.Name(), piVal, *f1)
		return
	}

	if *f2 != 0 {
		t.Errorf("Test %s failed. Want: [%f], got: [%f]", t.Name(), 0.0, *f2)
		return
	}
}

func TestFloatFail1(t *testing.T) {
	badArg := "stringNotANumber"
	testArgs := []string{"progname", "--float1", badArg}

	p := NewParser("", "description")
	f1 := p.Float("f", "float1", nil)

	err := p.Parse(testArgs)
	errStr := "[-f|--float1] bad floating point value [stringNotANumber]"
	if err == nil || err.Error() != errStr {
		t.Errorf("Test %s expected [%s], got [%+v]", t.Name(), errStr, err)
		return
	}

	if f1 == nil {
		t.Errorf("Test %s failed with float1 being nil pointer", t.Name())
		return
	}

	if *f1 != 0 {
		t.Errorf("Test %s failed. Want: [0], got: [%f]", t.Name(), *f1)
		return
	}
}

func TestFloatsSimple1(t *testing.T) {
	testval := 5.5
	testArgs := []string{"progname", "--flag-arg1", fmt.Sprintf("%.1f", testval)}

	p := NewParser("", "description")
	i1 := p.Floats("f", "flag-arg1", nil)
	i2 := p.Floats("", "flag-arg2", nil)

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if i1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if i2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if len(*i1) != 1 {
		t.Errorf("Test %s length failed. Want: [1], got: [%d]", t.Name(), len(*i1))
		return
	}

	if (*i1)[0] != testval {
		t.Errorf("Test %s value failed. Want: [%f], got: [%f]", t.Name(), testval, (*i1)[0])
		return
	}

	if len(*i2) != 0 {
		t.Errorf("Test %s failed. Want: [0] length, got: [%d]", t.Name(), len(*i2))
		return
	}
}

func TestFloatsNargsInvalidChar(t *testing.T) {
	testArgs := []string{"progname", "-f", "1", "2", "3"}

	p := NewParser("", "description")
	_ = p.Floats("f", "flag-arg1", &Options{Nargs: "x"})

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed, expected \"invalid string value ...\" error for --flag-arg1", t.Name())
		return
	}

	if err != nil && !strings.Contains(err.Error(), "invalid string value") {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}
}

func TestFloatsNargsInvalidNum(t *testing.T) {
	testArgs := []string{"progname", "-f", "1", "2", "3"}

	p := NewParser("", "description")
	_ = p.Floats("f", "flag-arg1", &Options{Nargs: -1})

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed, expected \"nargs integer value ...\" error for --flag-arg1", t.Name())
		return
	}

	if err != nil && !strings.Contains(err.Error(), "nargs integer value") {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}
}

// Strings nargs N test cases:
// 	* correct number of args
// 	* incorrect number of args
func TestFloatsNargsInt(t *testing.T) {
	testF1 := []float64{1, 2, 3}
	testF2 := []float64{5}
	testArgs := []string{
		"progname",
		"-f",
		fmt.Sprintf("%f", testF1[0]),
		fmt.Sprintf("%f", testF1[1]),
		fmt.Sprintf("%f", testF1[2]),
		"--flag-arg2",
		fmt.Sprintf("%f", testF2[0]),
		"-v",
	}

	p := NewParser("", "description")
	f1 := p.Floats("f", "flag-arg1", &Options{Nargs: 3})
	f2 := p.Floats("g", "flag-arg2", &Options{Nargs: 2})
	_ = p.Floats("v", "version", nil)

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed, expected \"not enough arguments ...\" error for --flag-arg2", t.Name())
		return
	}

	if err != nil && !strings.Contains(err.Error(), "not enough arguments") {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if f1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if f2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if !(len(*f1) == len(testF1) && (*f1)[0] == testF1[0]) {
		t.Errorf("Test %s i1 failed. Want: %v, got: %v", t.Name(), testF1, *f1)
		return
	}

	if len(*f2) != 0 {
		t.Errorf("Test %s i2 length failed. Want: [%d], got: [%d]", t.Name(), 0, len(*f2))
		return
	}
}

func TestFloatsNargsIntWithErr(t *testing.T) {
	testF1 := []float64{1, 2, 3}
	testF2 := []float64{5}
	testArgs := []string{
		"progname",
		"-f",
		fmt.Sprintf("%f", testF1[0]),
		fmt.Sprintf("%f", testF1[1]),
		fmt.Sprintf("%f", testF1[2]),
		"--flag-arg2",
		fmt.Sprintf("%f", testF2[0]),
	}
	p := NewParser("", "description")
	_ = p.Floats("f", "flag-arg1", &Options{Nargs: 3})
	_ = p.Floats("g", "flag-arg2", &Options{Nargs: -1})

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed, expected \"nargs integer value ...\" error for --flag-arg2", t.Name())
		return
	}

	if err != nil && !strings.Contains(err.Error(), "nargs integer value") {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}
}

// String nargs '?' test cases:
//	* one arg
//	* no args
//	* followed by a flag
//	* at beginning of arg list
//	* at end of arg list
func TestFloatNargs0or1(t *testing.T) {
	testval := 5.5
	testF2 := []float64{5, 6}
	testArgs := []string{
		"progname",
		"-f",
		fmt.Sprintf("%f", testval),
		"--flag-arg2",
		fmt.Sprintf("%f", testF2[0]),
		fmt.Sprintf("%f", testF2[1]),
		"-i",
		"-v",
		"-j",
	}

	p := NewParser("", "description")
	f1 := p.Float("f", "flag-arg1", &Options{Nargs: "?"})
	f2 := p.Floats("g", "flag-arg2", &Options{Nargs: 2})
	f3 := p.Float("i", "flag-arg3", &Options{Nargs: "?", Default: 3.0})
	f4 := p.Float("j", "flag-arg4", &Options{Nargs: "?", Default: 4.0})
	_ = p.Flag("v", "version", nil)

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if f1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if f2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if f3 == nil {
		t.Errorf("Test %s failed with flag3 being nil pointer", t.Name())
		return
	}

	if f4 == nil {
		t.Errorf("Test %s failed with flag4 being nil pointer", t.Name())
		return
	}

	if *f1 != testval {
		t.Errorf("Test %s i1 failed. Want: [%f], got: [%f]", t.Name(), testval, *f1)
		return
	}

	if *f3 != 3.0 {
		t.Errorf("Test %s i3 failed. Want: [%f], got: [%f]", t.Name(), 3.0, *f3)
		return
	}

	if *f4 != 4.0 {
		t.Errorf("Test %s i4 failed. Want: [%f], got: [%f]", t.Name(), 4.0, *f4)
		return
	}
}

// Strings nargs '*' test cases:
//	* no args
//	* multiple args
//	* at beginning of arg list
//	* at end of arg list
//	* multiple flags with '*'
func TestFloatsNargs0orMore(t *testing.T) {
	testF1 := []float64{1}
	testF3 := []float64{3, 4}
	testArgs := []string{
		"progname",
		"-f",
		fmt.Sprintf("%f", testF1[0]),
		"-g",
		"-i",
		fmt.Sprintf("%f", testF3[0]),
		fmt.Sprintf("%f", testF3[1]),
		"-j",
	}

	p := NewParser("", "description")
	f1 := p.Floats("f", "flag-arg1", &Options{Nargs: "*"})
	f2 := p.Floats("g", "flag-arg2", &Options{Nargs: "*"})
	f3 := p.Floats("i", "flag-arg3", &Options{Nargs: "*"})
	f4 := p.Floats("j", "flag-arg4", &Options{Nargs: "*"})

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if f1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if f2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if f3 == nil {
		t.Errorf("Test %s failed with flag3 being nil pointer", t.Name())
		return
	}

	if f4 == nil {
		t.Errorf("Test %s failed with flag4 being nil pointer", t.Name())
		return
	}

	if !(len(*f1) == len(testF1) && (*f1)[0] == testF1[0]) {
		t.Errorf("Test %s i1 failed. Want: %f, got: %f", t.Name(), testF1, *f1)
		return
	}

	if len(*f2) != 0 {
		t.Errorf("Test %s i2 length failed. Want: [0], got: [%d]", t.Name(), len(*f2))
		return
	}

	if !(len(*f3) == len(testF3) && (*f3)[0] == testF3[0]) {
		t.Errorf("Test %s i3 failed. Want: %v, got: %v", t.Name(), testF3, *f3)
		return
	}

	if len(*f4) != 0 {
		t.Errorf("Test %s i4 length failed. Want: [0], got: [%d]", t.Name(), len(*f4))
		return
	}
}

// Strings nargs '+' test cases:
//	* one arg
//	* multiple args
//	* at beginning of arg list
//	* at end of arg list
//	* multiple flags with '+'
func TestFloatsNargs1orMore(t *testing.T) {
	testF1 := []float64{1}
	testF2 := []float64{3, 4}
	testArgs := []string{
		"progname",
		"-f",
		fmt.Sprintf("%f", testF1[0]),
		"-g",
		fmt.Sprintf("%f", testF2[0]),
		fmt.Sprintf("%f", testF2[1]),
	}

	p := NewParser("", "description")
	f1 := p.Floats("f", "flag-arg1", &Options{Nargs: "+"})
	f2 := p.Floats("g", "flag-arg2", &Options{Nargs: "+"})

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if f1 == nil {
		t.Errorf("Test %s failed with flag1 being nil pointer", t.Name())
		return
	}

	if f2 == nil {
		t.Errorf("Test %s failed with flag2 being nil pointer", t.Name())
		return
	}

	if !(len(*f1) == len(testF1) && (*f1)[0] == testF1[0]) {
		t.Errorf("Test %s i1 failed. Want: %v, got: %v", t.Name(), testF1, *f1)
		return
	}

	if !(len(*f2) == len(testF2) && (*f2)[0] == testF2[0]) {
		t.Errorf("Test %s i2 failed. Want: %v, got: %v", t.Name(), testF2, *f2)
		return
	}
}

func TestFloatsNargs1orMoreWithErr(t *testing.T) {
	testArgs := []string{"progname", "-f", "-g", "3.0", "4.0"}

	p := NewParser("", "description")
	_ = p.Floats("f", "flag-arg1", &Options{Nargs: "+"})
	_ = p.Floats("i", "flag-arg2", &Options{Nargs: "+"})

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed, expected \"... at least one argument\" error for --flag-arg1", t.Name())
		return
	}

	if err != nil && !strings.Contains(err.Error(), "requires at least one argument") {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}
}

var pUsageString = `test string
usage: prog [-h|--help]

            program description

Arguments:

  -h  --help  Print help information

`

func TestUsageString(t *testing.T) {
	p := NewParser("prog", "program description")

	p.Parse(os.Args)

	usage := p.Usage("test string")

	if usage != pUsageString {
		t.Errorf("%s", usage)
	}
}

type s string

func (s s) String() string {
	return string(s)
}

var pUsageStringer = `stringer message
usage: prog [-h|--help]

            program description

Arguments:

  -h  --help  Print help information

`

func TestUsageStringer(t *testing.T) {
	p := NewParser("prog", "program description")

	p.Parse(os.Args)

	var msg s = "stringer message"

	usage := p.Usage(msg)

	if usage != pUsageStringer {
		t.Errorf("%s", usage)
	}
}
