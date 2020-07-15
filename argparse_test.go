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

func TestInternalFunctionParse(t *testing.T) {
	// common testing data
	a := &arg{
		sname:  "f",
		lname:  "flag",
		size:   2,
		opts:   nil,
		unique: true,
	}
	args0 := []string{}
	args2 := []string{"0", "1"}
	failureMessageCommon := "[-f|--flag] followed by too many arguments"

	// Fill testing table with testing cases
	type testCase struct {
		testName, failureMessage string
		resultInterface          interface{}
	}
	var (
		resultS     string
		resultI     int
		resultF     float64
		resultFile  os.File
		resultSL    []string
		resultIL    []int
		resultFL    []float64
		resultFileL []os.File
	)
	tt := []testCase{
		testCase{"String Value", "[-f|--flag] must be followed by a string", &resultS},
		testCase{"Int Value", "[-f|--flag] must be followed by an integer", &resultI},
		testCase{"Float Value", "[-f|--flag] must be followed by a floating point number", &resultF},
		testCase{"File Value", "[-f|--flag] must be followed by a path to file", &resultFile},
		testCase{"String Values List", "[-f|--flag] must be followed by a string", &resultSL},
		testCase{"Int Values List", "[-f|--flag] must be followed by an integer", &resultIL},
		testCase{"Float Values List", "[-f|--flag] must be followed by a floating point number", &resultFL},
		testCase{"File Values List", "[-f|--flag] must be followed by a path to file", &resultFileL},
	}

	//test all cases from table of cases
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			a.result = tc.resultInterface
			if err := a.parse(args0, 1); err == nil || err.Error() != tc.failureMessage {
				t.Errorf("Test %s failed with error: \"%v\". error: %q expected", t.Name(), err, tc.failureMessage)
			}
			a.parsed = false
			if err := a.parse(args2, 1); err == nil || err.Error() != failureMessageCommon {
				t.Errorf("Test %s failed with error: \"%v\". error: %q expected", t.Name(), err, failureMessageCommon)
			}
			a.parsed = false
		})
	}
}

func TestInternalFunctionCheck(t *testing.T) {
	var resultS string
	//test string
	a := &arg{
		result: &resultS,
		sname:  "f",
		lname:  "flag",
		size:   0,
		opts:   nil,
		unique: true,
	}

	srgString := "-f"
	failureMessage := "Argument's size < 1 is not allowed"

	if _, err := a.check(srgString); err == nil || err.Error() != failureMessage {
		t.Errorf("Test %s failed with error: \"%v\". error: %q expected", t.Name(), err, failureMessage)
	}
	a.parsed = false
}

func TestFlagAddArgumentFail(t *testing.T) {
	type testCase struct {
		testName, shortArg, longArg, failureMessage string
	}
	tt := []testCase{
		testCase{testName: "Long short name", shortArg: "ff", longArg: "flag2", failureMessage: "unable to add Flag: short name must not exceed 1 character"},
		testCase{testName: "Long name not provided", shortArg: "f", longArg: "", failureMessage: "unable to add Flag: long name should be provided"},
		testCase{testName: "Long name twice", shortArg: "f", longArg: "flag1", failureMessage: "unable to add Flag: long name flag1 occurs more than once"},
		testCase{testName: "Short name twice", shortArg: "F", longArg: "flag2", failureMessage: "unable to add Flag: short name F occurs more than once"},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					rezString := fmt.Sprintf("%v", r)
					if strings.Contains(rezString, tc.failureMessage) == false {
						t.Errorf("Test %s failed with panic result: \"%v\". panic result: %q expected", t.Name(), r, tc.failureMessage)
					}
				} else {
					t.Errorf("Test %s failed with no panic, but panic expected with result: %q", t.Name(), tc.failureMessage)
				}
			}()
			p := NewParser("", "description")
			_ = p.Flag("F", "flag1", nil)
			_ = p.Flag(tc.shortArg, tc.longArg, nil)
		})
	}
}

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

func TestLongFlagEqualChar(t *testing.T) {
	testArgs := []string{"progname", "--flag1=test1", "--flag2=2", "--flag3", "test3", "--flag4=a=test4", "--flag5=a"}

	p := NewParser("", "description")
	flag1 := p.String("", "flag1", nil)
	flag2 := p.Int("", "flag2", nil)
	flag3 := p.String("", "flag3", nil)
	flag4 := p.String("", "flag4=a", nil)
	flag5 := p.Flag("", "flag5=a", nil)

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
		t.Errorf("Test %s failed with flag3 being nil pointer", t.Name())
		return
	}

	if flag4 == nil {
		t.Errorf("Test %s failed with flag4 being nil pointer", t.Name())
		return
	}

	if flag5 == nil {
		t.Errorf("Test %s failed with flag5 being nil pointer", t.Name())
		return
	}

	if *flag1 != "test1" {
		t.Errorf("Test %s failed with flag1 being false", t.Name())
		return
	}

	if *flag2 != 2 {
		t.Errorf("Test %s failed with flag2 being true", t.Name())
		return
	}

	if *flag3 != "test3" {
		t.Errorf("Test %s failed with flag3 being true", t.Name())
		return
	}

	if *flag4 != "test4" {
		t.Errorf("Test %s failed with flag3 being true", t.Name())
		return
	}

	if *flag5 != true {
		t.Errorf("Test %s failed with flag3 being true", t.Name())
		return
	}
}

func TestShortFlagEqualChar(t *testing.T) {
	testArgs := []string{"progname", "-a=test1", "-b=2", "-c", "test3"}

	p := NewParser("", "description")
	flag1 := p.String("a", "flag1", nil)
	flag2 := p.Int("b", "flag2", nil)
	flag3 := p.String("c", "flag3", nil)

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
		t.Errorf("Test %s failed with flag3 being nil pointer", t.Name())
		return
	}

	if *flag1 != "test1" {
		t.Errorf("Test %s failed with flag1 being false", t.Name())
		return
	}

	if *flag2 != 2 {
		t.Errorf("Test %s failed with flag2 being true", t.Name())
		return
	}

	if *flag3 != "test3" {
		t.Errorf("Test %s failed with flag3 being true", t.Name())
		return
	}
}

func TestFlagMultiShorthandWithParam1(t *testing.T) {
	testArgs := []string{"progname", "-ab", "10", "-c", "-de", "11", "--ee", "12"}

	testList := []int{11, 12}

	p := NewParser("", "description")
	flag1 := p.Flag("a", "aa", nil)
	int2 := p.Int("b", "bb", nil)
	flag3 := p.Flag("c", "cc", nil)
	flag4 := p.Flag("d", "dd", nil)
	intList5 := p.IntList("e", "ee", nil)
	flag6 := p.Flag("f", "ff", nil)

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if *flag1 != true {
		t.Errorf("Test %s failed with flag1 being false", t.Name())
	}

	if *int2 != 10 {
		t.Errorf("Test %s failed with *int2=%v being false", t.Name(), *int2)
	}

	if *flag3 != true {
		t.Errorf("Test %s failed with flag3 being false", t.Name())
	}

	if *flag4 != true {
		t.Errorf("Test %s failed with flag4 being false", t.Name())
	}

	if !reflect.DeepEqual(*intList5, testList) {
		t.Errorf("Test %s failed: expected [%v], got [%v]", t.Name(), testList, *intList5)
	}

	if *flag6 != false {
		t.Errorf("Test %s failed with flag6 being true", t.Name())
	}
}

func TestFlagMultiShorthandWithParamFail1(t *testing.T) {
	testArgs := []string{"progname", "-bab", "10"}

	p := NewParser("", "description")
	_ = p.Flag("a", "aa", nil)
	_ = p.Int("b", "bb", nil)

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed with no error", t.Name())
		return
	}
	errExpectation := "[-b|--bb] argument: The parameter must follow"
	if err.Error() != errExpectation {
		t.Errorf("Test %s failed. error %q getted. %q expected", t.Name(), err.Error(), errExpectation)
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

func TestFlagCounterAddArgumentFail(t *testing.T) {
	type testCase struct {
		testName, shortArg, longArg, failureMessage string
	}
	tt := []testCase{
		testCase{testName: "Long short name", shortArg: "ff", longArg: "flag2", failureMessage: "unable to add FlagCounter: short name must not exceed 1 character"},
		testCase{testName: "Long name not provided", shortArg: "f", longArg: "", failureMessage: "unable to add FlagCounter: long name should be provided"},
		testCase{testName: "Long name twice", shortArg: "f", longArg: "flag1", failureMessage: "unable to add FlagCounter: long name flag1 occurs more than once"},
		testCase{testName: "Short name twice", shortArg: "F", longArg: "flag2", failureMessage: "unable to add FlagCounter: short name F occurs more than once"},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					rezString := fmt.Sprintf("%v", r)
					if strings.Contains(rezString, tc.failureMessage) == false {
						t.Errorf("Test %s failed with panic result: \"%v\". panic result: %q expected", t.Name(), r, tc.failureMessage)
					}
				} else {
					t.Errorf("Test %s failed with no panic, but panic expected with result: %q", t.Name(), tc.failureMessage)
				}
			}()
			p := NewParser("", "description")
			_ = p.FlagCounter("F", "flag1", nil)
			_ = p.FlagCounter(tc.shortArg, tc.longArg, nil)
		})
	}
}

func TestFlagCounterSimple1(t *testing.T) {
	testArgs := []string{"progname", "--flag-arg1", "--flag-arg3", "--flag-arg3"}

	p := NewParser("", "description")
	flag1 := p.FlagCounter("", "flag-arg1", nil)
	flag2 := p.FlagCounter("", "flag-arg2", nil)
	flag3 := p.FlagCounter("", "flag-arg3", nil)

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
		t.Errorf("Test %s failed with flag3 being nil pointer", t.Name())
		return
	}

	if *flag1 != 1 {
		t.Errorf("Test %s failed with flag1 being %d", t.Name(), *flag1)
		return
	}

	if *flag2 != 0 {
		t.Errorf("Test %s failed with flag2 being %d", t.Name(), *flag2)
		return
	}

	if *flag3 != 2 {
		t.Errorf("Test %s failed with flag3 being %d", t.Name(), *flag3)
		return
	}
}

func TestFlagCounterSimple2(t *testing.T) {
	testArgs := []string{"progname", "--flag-arg1", "-f", "--flag-arg3", "-f"}

	p := NewParser("", "description")
	flag1 := p.FlagCounter("", "flag-arg1", nil)
	flag2 := p.FlagCounter("", "flag-arg2", nil)
	flag3 := p.FlagCounter("f", "flag-arg3", nil)

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
		t.Errorf("Test %s failed with flag3 being nil pointer", t.Name())
		return
	}

	if *flag1 != 1 {
		t.Errorf("Test %s failed with flag1 being %d", t.Name(), *flag1)
		return
	}

	if *flag2 != 0 {
		t.Errorf("Test %s failed with flag2 being %d", t.Name(), *flag2)
		return
	}

	if *flag3 != 3 {
		t.Errorf("Test %s failed with flag3 being %d", t.Name(), *flag3)
		return
	}
}

func TestFlagCounterMultiShorthand1(t *testing.T) {
	testArgs := []string{"progname", "-abbcbcadaa", "-e"}

	p := NewParser("", "description")
	flag1 := p.FlagCounter("a", "aa", nil)
	flag2 := p.FlagCounter("b", "bb", nil)
	flag3 := p.FlagCounter("c", "cc", nil)
	flag4 := p.FlagCounter("d", "dd", nil)
	flag5 := p.FlagCounter("e", "ee", nil)
	flag6 := p.FlagCounter("f", "ff", nil)

	err := p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
		return
	}

	if *flag1 != 4 {
		t.Errorf("Test %s failed with flag1 being %d", t.Name(), *flag1)
	}

	if *flag2 != 3 {
		t.Errorf("Test %s failed with flag2 being %d", t.Name(), *flag2)
	}

	if *flag3 != 2 {
		t.Errorf("Test %s failed with flag3 being %d", t.Name(), *flag3)
	}

	if *flag4 != 1 {
		t.Errorf("Test %s failed with flag4 being %d", t.Name(), *flag4)
	}

	if *flag5 != 1 {
		t.Errorf("Test %s failed with flag5 being %d", t.Name(), *flag5)
	}

	if *flag6 != 0 {
		t.Errorf("Test %s failed with flag6 being %d", t.Name(), *flag6)
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

	testArgs = []string{"progname", "--flag-arg2", "-ff"}

	p = NewParser("", "description")
	_ = p.Flag("f", "flag-arg1", nil)
	_ = p.Flag("", "flag-arg2", nil)

	err = p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed with. Duplicate flag use not detected", t.Name())
		return
	}

	testArgs = []string{"progname", "--flag-arg2", "-f"}

	p = NewParser("", "description")
	_ = p.Flag("f", "flag-arg1", nil)
	_ = p.Flag("", "flag-arg2", nil)

	err = p.Parse(testArgs)
	if err != nil {
		t.Errorf("Test %s failed with. Fake duplicate flag detected", t.Name())
		return
	}
}

func TestFailCaseSensitive(t *testing.T) {
	testArgs := []string{"progname", "-F"}

	p := NewParser("", "description")
	_ = p.Flag("f", "flag", &Options{Required: true})

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

func TestStringAddArgumentFail(t *testing.T) {
	type testCase struct {
		testName, shortArg, longArg, failureMessage string
	}
	tt := []testCase{
		testCase{testName: "Long short name", shortArg: "ff", longArg: "flag2", failureMessage: "unable to add String: short name must not exceed 1 character"},
		testCase{testName: "Long name not provided", shortArg: "f", longArg: "", failureMessage: "unable to add String: long name should be provided"},
		testCase{testName: "Long name twice", shortArg: "f", longArg: "flag1", failureMessage: "unable to add String: long name flag1 occurs more than once"},
		testCase{testName: "Short name twice", shortArg: "F", longArg: "flag2", failureMessage: "unable to add String: short name F occurs more than once"},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					rezString := fmt.Sprintf("%v", r)
					if strings.Contains(rezString, tc.failureMessage) == false {
						t.Errorf("Test %s failed with panic result: \"%v\". panic result: %q expected", t.Name(), r, tc.failureMessage)
					}
				} else {
					t.Errorf("Test %s failed with no panic, but panic expected with result: %q", t.Name(), tc.failureMessage)
				}
			}()
			p := NewParser("", "description")
			_ = p.String("F", "flag1", nil)
			_ = p.String(tc.shortArg, tc.longArg, nil)
		})
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

func TestIntAddArgumentFail(t *testing.T) {
	type testCase struct {
		testName, shortArg, longArg, failureMessage string
	}
	tt := []testCase{
		testCase{testName: "Long short name", shortArg: "ff", longArg: "flag2", failureMessage: "unable to add Int: short name must not exceed 1 character"},
		testCase{testName: "Long name not provided", shortArg: "f", longArg: "", failureMessage: "unable to add Int: long name should be provided"},
		testCase{testName: "Long name twice", shortArg: "f", longArg: "flag1", failureMessage: "unable to add Int: long name flag1 occurs more than once"},
		testCase{testName: "Short name twice", shortArg: "F", longArg: "flag2", failureMessage: "unable to add Int: short name F occurs more than once"},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					rezString := fmt.Sprintf("%v", r)
					if strings.Contains(rezString, tc.failureMessage) == false {
						t.Errorf("Test %s failed with panic result: \"%v\". panic result: %q expected", t.Name(), r, tc.failureMessage)
					}
				} else {
					t.Errorf("Test %s failed with no panic, but panic expected with result: %q", t.Name(), tc.failureMessage)
				}
			}()
			p := NewParser("", "description")
			_ = p.Int("F", "flag1", nil)
			_ = p.Int(tc.shortArg, tc.longArg, nil)
		})
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

func TestIntFailSimple1(t *testing.T) {
	testArgs := []string{"progname", "--flag-arg1", "string"}

	p := NewParser("", "description")
	i1 := p.Int("f", "flag-arg1", nil)

	err := p.Parse(testArgs)
	errStr := "[-f|--flag-arg1] bad integer value [string]"
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

func TestEqualIntFailSimple1(t *testing.T) {
	testArgs := []string{"progname", "--flag-arg1=string"}

	p := NewParser("", "description")
	i1 := p.Int("f", "flag-arg1", nil)

	err := p.Parse(testArgs)
	errStr := "[-f|--flag-arg1] bad integer value [string]"
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

func TestEqualNoValFailSimple(t *testing.T) {
	testArgs := []string{"progname", "--flag-arg1="}

	p := NewParser("", "description")
	i1 := p.Int("f", "flag-arg1", nil)

	err := p.Parse(testArgs)
	errStr := "not enough arguments for -f|--flag-arg1"
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

func TestFileAddArgumentFail(t *testing.T) {
	type testCase struct {
		testName, shortArg, longArg, failureMessage string
	}
	tt := []testCase{
		testCase{testName: "Long short name", shortArg: "ff", longArg: "flag2", failureMessage: "unable to add File: short name must not exceed 1 character"},
		testCase{testName: "Long name not provided", shortArg: "f", longArg: "", failureMessage: "unable to add File: long name should be provided"},
		testCase{testName: "Long name twice", shortArg: "f", longArg: "flag1", failureMessage: "unable to add File: long name flag1 occurs more than once"},
		testCase{testName: "Short name twice", shortArg: "F", longArg: "flag2", failureMessage: "unable to add File: short name F occurs more than once"},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					rezString := fmt.Sprintf("%v", r)
					if strings.Contains(rezString, tc.failureMessage) == false {
						t.Errorf("Test %s failed with panic result: \"%v\". panic result: %q expected", t.Name(), r, tc.failureMessage)
					}
				} else {
					t.Errorf("Test %s failed with no panic, but panic expected with result: %q", t.Name(), tc.failureMessage)
				}
			}()
			p := NewParser("", "description")
			_ = p.File("F", "flag1", os.O_RDWR, 0666, nil)
			_ = p.File(tc.shortArg, tc.longArg, os.O_RDWR, 0666, nil)
		})
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
	if file1 == nil {
		t.Errorf("Test %s failed with file1 being nil pointer", t.Name())
		return
	}

	defer file1.Close()

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

func TestFileSimpleFail1(t *testing.T) {
	// Not existing test file location
	fpath := "./non-existent-file.tmp"
	// To be shure there is no fake file
	if _, err := os.Stat(fpath); os.IsNotExist(err) != true {
		//we could remove it, but what if it's important
		t.Errorf("Test %s failed. There is \"%s\" file in module directory, which must not exists for test purposes", t.Name(), fpath)
		return
	}

	testArgs := []string{"progname"}

	p := NewParser("", "")

	_ = p.File("f", "file", os.O_RDWR, 0666, &Options{Default: fpath})

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed. Parsing should fail.", t.Name())
	}
	err, ok := err.(*os.PathError)

	if ok == false {
		t.Errorf("Test %s failed with error: %s, that is not of *os.PathError type", t.Name(), err.Error())
	}
}

func TestFileSimpleFail2(t *testing.T) {
	// Not existing test file location
	fpath := "./non-existent-file.tmp"
	// To be shure there is no fake file
	if _, err := os.Stat(fpath); os.IsNotExist(err) != true {
		//we could remove it, but what if it's important
		t.Errorf("Test %s failed. There is \"%s\" file in module directory, which must not exists for test purposes", t.Name(), fpath)
		return
	}

	testArgs := []string{"progname", "-f", fpath}

	p := NewParser("", "")

	_ = p.File("f", "file", os.O_RDWR, 0666, nil)

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed. Parsing should fail.", t.Name())
		return
	}
	err, ok := err.(*os.PathError)

	if ok == false {
		t.Errorf("Test %s failed with error: %s, that is not of *os.PathError type", t.Name(), err.Error())
	}
}

func TestFileListSimpleFail1(t *testing.T) {
	// Test files location
	fpaths := []string{"./test1.tmp", "./non-existent-file2.tmp", "./test2.tmp"}
	// Create test files
	for i, fpath := range fpaths {
		if i == 1 {
			if _, err := os.Stat(fpath); os.IsNotExist(err) != true {
				//we could remove it, but what if it's important
				t.Errorf("Test %s failed. There is \"%s\" file in module directory, which must not exists for test purposes", t.Name(), fpath)
				return
			}
		} else {
			f, err := os.Create(fpath)
			if err != nil {
				t.Error(err)
				return
			}
			f.Close()
			defer os.Remove(fpath)
		}
	}

	testArgs := []string{"progname"}

	p := NewParser("", "")

	files := p.FileList("f", "file", os.O_RDWR, 0666, &Options{Default: fpaths})

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed. Parsing should fail.", t.Name())
	}
	if len(*files) > 0 {
		t.Errorf("Test %s failed. File list should be empty.", t.Name())
	}
}

func TestFileListSimpleFail2(t *testing.T) {
	// Test files location
	fpaths := []string{"./test1.tmp", "./non-existent-file2.tmp", "./test2.tmp"}
	// Create test files
	for i, fpath := range fpaths {
		if i == 1 {
			if _, err := os.Stat(fpath); os.IsNotExist(err) != true {
				//we could remove it, but what if it's important
				t.Errorf("Test %s failed. There is \"%s\" file in module directory, which must not exists for test purposes", t.Name(), fpath)
				return
			}
		} else {
			f, err := os.Create(fpath)
			if err != nil {
				t.Error(err)
				return
			}
			f.Close()
			defer os.Remove(fpath)
		}
	}

	testArgs := []string{"progname", "-f", fpaths[0], "--file", fpaths[1], "-f", fpaths[2]}

	p := NewParser("", "")

	files := p.FileList("f", "file", os.O_RDWR, 0666, &Options{Default: nil})

	err := p.Parse(testArgs)
	if err == nil {
		t.Errorf("Test %s failed. Parsing should fail.", t.Name())
	}
	if len(*files) > 0 {
		t.Errorf("Test %s failed. File list should be empty.", t.Name())
	}
}

func TestFileListAddArgumentFail(t *testing.T) {
	type testCase struct {
		testName, shortArg, longArg, failureMessage string
	}
	tt := []testCase{
		testCase{testName: "Long short name", shortArg: "ff", longArg: "flag2", failureMessage: "unable to add FileList: short name must not exceed 1 character"},
		testCase{testName: "Long name not provided", shortArg: "f", longArg: "", failureMessage: "unable to add FileList: long name should be provided"},
		testCase{testName: "Long name twice", shortArg: "f", longArg: "flag1", failureMessage: "unable to add FileList: long name flag1 occurs more than once"},
		testCase{testName: "Short name twice", shortArg: "F", longArg: "flag2", failureMessage: "unable to add FileList: short name F occurs more than once"},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					rezString := fmt.Sprintf("%v", r)
					if strings.Contains(rezString, tc.failureMessage) == false {
						t.Errorf("Test %s failed with panic result: \"%v\". panic result: %q expected", t.Name(), r, tc.failureMessage)
					}
				} else {
					t.Errorf("Test %s failed with no panic, but panic expected with result: %q", t.Name(), tc.failureMessage)
				}
			}()
			p := NewParser("", "description")
			_ = p.FileList("F", "flag1", os.O_RDWR, 0666, nil)
			_ = p.FileList(tc.shortArg, tc.longArg, os.O_RDWR, 0666, nil)
		})
	}
}

func TestFileListSimple1(t *testing.T) {
	// Test files location
	fpaths := []string{"./test1.tmp", "./test2.tmp"}
	// Create test files
	for _, fpath := range fpaths {
		f, err := os.Create(fpath)
		if err != nil {
			t.Error(err)
			return
		}
		f.Close()
		defer os.Remove(fpath)
	}

	testArgs := []string{"progname", "-f", fpaths[0], "--file", fpaths[1]}

	p := NewParser("", "")

	files := p.FileList("f", "file", os.O_RDWR, 0666, &Options{Default: []string{"./non-existent-file1.tmp", "./non-existent-file2.tmp"}})

	err := p.Parse(testArgs)
	switch {
	case err != nil:
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
	case files == nil:
		t.Errorf("Test %s failed with l1 being nil pointer", t.Name())
	}
	for i, file := range *files {
		defer file.Close()
		testString := "Test"
		recSlice := make([]byte, 4)
		_, err = file.WriteString(testString)
		if err != nil {
			t.Errorf("Test %s write operation with file: %s failed with error: %s", t.Name(), fpaths[i], err.Error())
			return
		}
		file.Seek(0, 0)
		n, err := file.Read(recSlice)
		if err != nil {
			t.Errorf("Test %s read operation with file: %s failed with error: %s", t.Name(), fpaths[i], err.Error())
			return
		}
		if n != 4 || string(recSlice) != testString {
			t.Errorf("Test %s failed with file: %s on read operation", t.Name(), fpaths[i])
			return
		}
	}
}

func TestFloatListAddArgumentFail(t *testing.T) {
	type testCase struct {
		testName, shortArg, longArg, failureMessage string
	}
	tt := []testCase{
		testCase{testName: "Long short name", shortArg: "ff", longArg: "flag2", failureMessage: "unable to add FloatList: short name must not exceed 1 character"},
		testCase{testName: "Long name not provided", shortArg: "f", longArg: "", failureMessage: "unable to add FloatList: long name should be provided"},
		testCase{testName: "Long name twice", shortArg: "f", longArg: "flag1", failureMessage: "unable to add FloatList: long name flag1 occurs more than once"},
		testCase{testName: "Short name twice", shortArg: "F", longArg: "flag2", failureMessage: "unable to add FloatList: short name F occurs more than once"},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					rezString := fmt.Sprintf("%v", r)
					if strings.Contains(rezString, tc.failureMessage) == false {
						t.Errorf("Test %s failed with panic result: \"%v\". panic result: %q expected", t.Name(), r, tc.failureMessage)
					}
				} else {
					t.Errorf("Test %s failed with no panic, but panic expected with result: %q", t.Name(), tc.failureMessage)
				}
			}()
			p := NewParser("", "description")
			_ = p.FloatList("F", "flag1", nil)
			_ = p.FloatList(tc.shortArg, tc.longArg, nil)
		})
	}
}

func TestFloatListSimple1(t *testing.T) {
	testArgs := []string{"progname", "--flag-arg1", "12", "--flag-arg1", "-10.1", "--flag-arg1", "+10"}
	list1Expect := []float64{12, -10.1, 10}
	list2Expect := make([]float64, 0)

	p := NewParser("", "description")
	l1 := p.FloatList("f", "flag-arg1", nil)
	l2 := p.FloatList("", "flag-arg2", nil)

	err := p.Parse(testArgs)
	switch {
	case err != nil:
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
	case l1 == nil:
		t.Errorf("Test %s failed with l1 being nil pointer", t.Name())
	case l2 == nil:
		t.Errorf("Test %s failed with l2 being nil pointer", t.Name())
	case !reflect.DeepEqual(*l1, list1Expect):
		t.Errorf("Test %s failed. Want: %f, got: %f", t.Name(), list1Expect, *l1)
	case !reflect.DeepEqual(*l2, list2Expect):
		t.Errorf("Test %s failed. Want: %f, got: %f", t.Name(), list2Expect, *l2)
	}
}

func TestFloatListTypeFail(t *testing.T) {
	testArgs := []string{"progname", "--flag-arg1", "12", "--flag-arg1", "10,1"}

	p := NewParser("", "description")
	p.FloatList("f", "flag-arg1", nil)

	err := p.Parse(testArgs)
	failureText := "[-f|--flag-arg1] bad floating point value [10,1]"
	if err == nil || err.Error() != failureText {
		t.Errorf("Test %s failed: expected error: [%s], got error: [%+v]", t.Name(), failureText, err)
	}
}

func TestIntListAddArgumentFail(t *testing.T) {
	type testCase struct {
		testName, shortArg, longArg, failureMessage string
	}
	tt := []testCase{
		testCase{testName: "Long short name", shortArg: "ff", longArg: "flag2", failureMessage: "unable to add IntList: short name must not exceed 1 character"},
		testCase{testName: "Long name not provided", shortArg: "f", longArg: "", failureMessage: "unable to add IntList: long name should be provided"},
		testCase{testName: "Long name twice", shortArg: "f", longArg: "flag1", failureMessage: "unable to add IntList: long name flag1 occurs more than once"},
		testCase{testName: "Short name twice", shortArg: "F", longArg: "flag2", failureMessage: "unable to add IntList: short name F occurs more than once"},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					rezString := fmt.Sprintf("%v", r)
					if strings.Contains(rezString, tc.failureMessage) == false {
						t.Errorf("Test %s failed with panic result: \"%v\". panic result: %q expected", t.Name(), r, tc.failureMessage)
					}
				} else {
					t.Errorf("Test %s failed with no panic, but panic expected with result: %q", t.Name(), tc.failureMessage)
				}
			}()
			p := NewParser("", "description")
			_ = p.IntList("F", "flag1", nil)
			_ = p.IntList(tc.shortArg, tc.longArg, nil)
		})
	}
}

func TestIntListSimple1(t *testing.T) {
	testArgs := []string{"progname", "--flag-arg1", "12", "--flag-arg1", "-10", "--flag-arg1", "+10"}
	list1Expect := []int{12, -10, 10}
	list2Expect := make([]int, 0)

	p := NewParser("", "description")
	l1 := p.IntList("f", "flag-arg1", nil)
	l2 := p.IntList("", "flag-arg2", nil)

	err := p.Parse(testArgs)
	switch {
	case err != nil:
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
	case l1 == nil:
		t.Errorf("Test %s failed with l1 being nil pointer", t.Name())
	case l2 == nil:
		t.Errorf("Test %s failed with l2 being nil pointer", t.Name())
	case !reflect.DeepEqual(*l1, list1Expect):
		t.Errorf("Test %s failed. Want: %d, got: %d", t.Name(), list1Expect, *l1)
	case !reflect.DeepEqual(*l2, list2Expect):
		t.Errorf("Test %s failed. Want: %d, got: %d", t.Name(), list2Expect, *l2)
	}
}

func TestIntListTypeFail(t *testing.T) {
	testArgs := []string{"progname", "--flag-arg1", "12", "--flag-arg1", "=10"}

	p := NewParser("", "description")
	p.IntList("f", "flag-arg1", nil)

	err := p.Parse(testArgs)
	failureText := "[-f|--flag-arg1] bad integer value [=10]"
	if err == nil || err.Error() != failureText {
		t.Errorf("Test %s failed: expected error: [%s], got error: [%+v]", t.Name(), failureText, err)
	}
}

func TestStringListAddArgumentFail(t *testing.T) {
	type testCase struct {
		testName, shortArg, longArg, failureMessage string
	}
	tt := []testCase{
		testCase{testName: "Long short name", shortArg: "ff", longArg: "flag2", failureMessage: "unable to add StringList: short name must not exceed 1 character"},
		testCase{testName: "Long name not provided", shortArg: "f", longArg: "", failureMessage: "unable to add StringList: long name should be provided"},
		testCase{testName: "Long name twice", shortArg: "f", longArg: "flag1", failureMessage: "unable to add StringList: long name flag1 occurs more than once"},
		testCase{testName: "Short name twice", shortArg: "F", longArg: "flag2", failureMessage: "unable to add StringList: short name F occurs more than once"},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					rezString := fmt.Sprintf("%v", r)
					if strings.Contains(rezString, tc.failureMessage) == false {
						t.Errorf("Test %s failed with panic result: \"%v\". panic result: %q expected", t.Name(), r, tc.failureMessage)
					}
				} else {
					t.Errorf("Test %s failed with no panic, but panic expected with result: %q", t.Name(), tc.failureMessage)
				}
			}()
			p := NewParser("", "description")
			_ = p.StringList("F", "flag1", nil)
			_ = p.StringList(tc.shortArg, tc.longArg, nil)
		})
	}
}

func TestStringListSimple1(t *testing.T) {
	testArgs := []string{"progname", "--flag-arg1", "test1", "--flag-arg1", "test2"}
	list1Expect := []string{"test1", "test2"}
	list2Expect := make([]string, 0)

	p := NewParser("", "description")
	l1 := p.StringList("f", "flag-arg1", nil)
	l2 := p.StringList("", "flag-arg2", nil)

	err := p.Parse(testArgs)
	switch {
	case err != nil:
		t.Errorf("Test %s failed with error: %s", t.Name(), err.Error())
	case l1 == nil:
		t.Errorf("Test %s failed with l1 being nil pointer", t.Name())
	case l2 == nil:
		t.Errorf("Test %s failed with l2 being nil pointer", t.Name())
	case !reflect.DeepEqual(*l1, list1Expect):
		t.Errorf("Test %s failed. Want: %s, got: %s", t.Name(), list1Expect, *l1)
	case !reflect.DeepEqual(*l2, list2Expect):
		t.Errorf("Test %s failed. Want: %s, got: %s", t.Name(), list2Expect, *l2)
	}
}

func TestListAddArgumentFail(t *testing.T) {
	type testCase struct {
		testName, shortArg, longArg, failureMessage string
	}
	tt := []testCase{
		testCase{testName: "Long short name", shortArg: "ff", longArg: "flag2", failureMessage: "unable to add StringList: short name must not exceed 1 character"},
		testCase{testName: "Long name not provided", shortArg: "f", longArg: "", failureMessage: "unable to add StringList: long name should be provided"},
		testCase{testName: "Long name twice", shortArg: "f", longArg: "flag1", failureMessage: "unable to add StringList: long name flag1 occurs more than once"},
		testCase{testName: "Short name twice", shortArg: "F", longArg: "flag2", failureMessage: "unable to add StringList: short name F occurs more than once"},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					rezString := fmt.Sprintf("%v", r)
					if strings.Contains(rezString, tc.failureMessage) == false {
						t.Errorf("Test %s failed with panic result: \"%v\". panic result: %q expected", t.Name(), r, tc.failureMessage)
					}
				} else {
					t.Errorf("Test %s failed with no panic, but panic expected with result: %q", t.Name(), tc.failureMessage)
				}
			}()
			p := NewParser("", "description")
			_ = p.StringList("F", "flag1", nil)
			_ = p.StringList(tc.shortArg, tc.longArg, nil)
		})
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

func TestSelectorAddArgumentFail(t *testing.T) {
	type testCase struct {
		testName, shortArg, longArg, failureMessage string
	}
	tt := []testCase{
		testCase{testName: "Long short name", shortArg: "ff", longArg: "flag2", failureMessage: "unable to add Selector: short name must not exceed 1 character"},
		testCase{testName: "Long name not provided", shortArg: "f", longArg: "", failureMessage: "unable to add Selector: long name should be provided"},
		testCase{testName: "Long name twice", shortArg: "f", longArg: "flag1", failureMessage: "unable to add Selector: long name flag1 occurs more than once"},
		testCase{testName: "Short name twice", shortArg: "F", longArg: "flag2", failureMessage: "unable to add Selector: short name F occurs more than once"},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					rezString := fmt.Sprintf("%v", r)
					if strings.Contains(rezString, tc.failureMessage) == false {
						t.Errorf("Test %s failed with panic result: \"%v\". panic result: %q expected", t.Name(), r, tc.failureMessage)
					}
				} else {
					t.Errorf("Test %s failed with no panic, but panic expected with result: %q", t.Name(), tc.failureMessage)
				}
			}()
			p := NewParser("", "description")
			allowedValues := []string{"test1", "test2"}
			_ = p.Selector("F", "flag1", allowedValues, nil)
			_ = p.Selector(tc.shortArg, tc.longArg, allowedValues, nil)
		})
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
		type commandCase struct {
			cmd        *Command
			cmd1flag   bool
			cmd1string string
			cmd1int    int
		}
		ct := []commandCase{
			commandCase{
				cmd:        cmd1,
				cmd1flag:   true,
				cmd1string: "test",
				cmd1int:    val,
			},
			commandCase{
				cmd:        cmd2,
				cmd1flag:   false,
				cmd1string: "",
				cmd1int:    0,
			},
		}

		for _, cc := range ct {
			if cc.cmd.Happened() {
				if *cmd1flag1 != cc.cmd1flag {
					t.Errorf("Test %s failed with %s: flag1: wanted [%t], got [%t]", t.Name(), testArgs[1], cc.cmd1flag, *cmd1flag1)
					return
				}
				if *cmd1string1 != cc.cmd1string {
					t.Errorf("Test %s failed with %s: string1: wanted [%s], got [%s]", t.Name(), testArgs[1], cc.cmd1string, *cmd1string1)
					return
				}
				if *cmd1int1 != cc.cmd1int {
					t.Errorf("Test %s failed with %s: int1: wanted [%d], got [%d]", t.Name(), testArgs[1], cc.cmd1int, *cmd1int1)
					return
				}
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
		t.Errorf("pUsage: get:\n%s\nexpect:\n%s", p.Usage(nil), pUsage)
	}
	if cmd1Usage != cmd1.Usage(nil) {
		t.Errorf("cmd1Usage: get:\n%s\nexpect:\n%s", cmd1.Usage(nil), cmd1Usage)
	}
	if cmd2Usage != cmd2.Usage(nil) {
		t.Errorf("cmd2Usage: get:\n%s\nexpect:\n%s", cmd2.Usage(nil), cmd2Usage)
	}
}

func TestUsageHidden1(t *testing.T) {
	p := NewParser("verylongprogname", "prog description")

	cmd1 := p.NewCommand("veryverylongcmd1", "cmd1 description")
	_ = cmd1.Flag("f", "verylongflag1", &Options{Help: "flag1 description"})
	_ = cmd1.Flag("a", "verylongflagA", &Options{Required: true, Help: "flag1 description"})
	_ = p.String("s", "verylongstring-flag1", &Options{Help: "string1 description"})
	_ = p.Int("i", "integer-flag1", &Options{Help: "integer1 description"})
	_ = p.Int("I", "integer-flag2", &Options{Help: DisableDescription})

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

func TestUsageSubCommand(t *testing.T) {
	expected := `[sub]Command required
usage: zooprog <Command> [-h|--help]

               Program that walks us through the zoo

Commands:

  dog  We are going to see dog

Arguments:

  -h  --help  Print help information

`

	parser := NewParser("zooprog", "Program that walks us through the zoo")

	// dog command
	parser.
		NewCommand("dog", "We are going to see dog"). // adds command to parser
		NewCommand("speak", "Make the dog speak")     // adds subcommand to previous command

	err := newSubCommandError(&parser.Command)
	actual := parser.Usage(err)
	if expected != actual {
		t.Errorf("Expectations unmet. expected: %s, actual: %s", expected, actual)
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
	if err == nil || err.Error() != "cannot use default type [string] as value of pointer with type [*bool]" {
		t.Errorf("Test %s failed: expected error [%s], got error [%+v]", t.Name(), "cannot use default type [string] as value of pointer with type [*bool]", err)
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
	if err == nil || err.Error() != "cannot use default type [bool] as value of pointer with type [*string]" {
		t.Errorf("Test %s failed: expected error [%s], got error [%+v]", t.Name(), "cannot use default type [bool] as value of pointer with type [*string]", err)
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
	if err == nil || err.Error() != "cannot use default type [string] as value of pointer with type [*int]" {
		t.Errorf("Test %s failed: expected error [%s], got error [%+v]", t.Name(), "cannot use default type [bool] as value of pointer with type [*string]", err)
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
	if err == nil || err.Error() != "cannot use default type [bool] as value of pointer with type [*string]" {
		t.Errorf("Test %s failed: expected error [%s], got error [%+v]", t.Name(), "cannot use default type [bool] as value of pointer with type [*string]", err)
	}
	defer file1.Close()
}

func TestFileListDefaultValuePass(t *testing.T) {
	testArgs := []string{"progname"}
	// Test files location
	fpaths := []string{"./test1.tmp", "./test2.tmp"}
	// Create test files
	for _, fpath := range fpaths {
		f, err := os.Create(fpath)
		if err != nil {
			t.Error(err)
			return
		}
		f.Close()
		defer os.Remove(fpath)
	}

	p := NewParser("progname", "Prog description")

	files := p.FileList("f", "float", os.O_RDWR, 0666, &Options{Default: fpaths})

	err := p.Parse(testArgs)

	if err != nil {
		t.Error(err.Error())
	}
	for i, file := range *files {
		defer file.Close()
		testString := "Test"
		recSlice := make([]byte, 4)
		_, err = file.WriteString(testString)
		if err != nil {
			t.Errorf("Test %s write operation with file: %s failed with error: %s", t.Name(), fpaths[i], err.Error())
			return
		}
		file.Seek(0, 0)
		n, err := file.Read(recSlice)
		if err != nil {
			t.Errorf("Test %s read operation with file: %s failed with error: %s", t.Name(), fpaths[i], err.Error())
			return
		}
		if n != 4 || string(recSlice) != testString {
			t.Errorf("Test %s failed with file: %s on read operation", t.Name(), fpaths[i])
			return
		}
	}

}

func TestFloatListDefaultValuePass(t *testing.T) {
	testArgs := []string{"progname"}
	testList := []float64{12.0, -10}

	p := NewParser("progname", "Prog description")

	s := p.FloatList("f", "float", &Options{Default: testList})

	err := p.Parse(testArgs)

	switch {
	// Should fail on failure
	case err != nil:
		t.Error(err.Error())
	// Should fail if not true
	case !reflect.DeepEqual(*s, testList):
		t.Errorf("expected [%v], got [%v]", testList, *s)
	}
}

func TestIntListDefaultValuePass(t *testing.T) {
	testArgs := []string{"progname"}
	testList := []int{12, -10}

	p := NewParser("progname", "Prog description")

	s := p.IntList("i", "int", &Options{Default: testList})

	err := p.Parse(testArgs)

	switch {
	// Should fail on failure
	case err != nil:
		t.Error(err.Error())
	// Should fail if not true
	case !reflect.DeepEqual(*s, testList):
		t.Errorf("expected [%v], got [%v]", testList, *s)
	}
}

func TestStringListDefaultValuePass(t *testing.T) {
	testArgs := []string{"progname"}
	testList := []string{"test", "list"}

	p := NewParser("progname", "Prog description")

	s := p.StringList("s", "string", &Options{Default: testList})

	err := p.Parse(testArgs)

	switch {
	// Should fail on failure
	case err != nil:
		t.Error(err.Error())
	// Should fail if not true
	case !reflect.DeepEqual(*s, testList):
		t.Errorf("expected [%v], got [%v]", testList, *s)
	}
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

func TestFileListDefaultValueFail(t *testing.T) {
	testArgs := []string{"progname"}

	p := NewParser("progname", "Prog description")

	_ = p.FileList("f", "float", os.O_RDWR, 0666, &Options{Default: false})

	err := p.Parse(testArgs)

	// Should pass on failure
	failureMessage := "cannot use default type [bool] as value of pointer with type [*[]string]"
	if err == nil || err.Error() != failureMessage {
		t.Errorf("Test %s failed: expected error [%s], got error [%+v]", t.Name(), failureMessage, err)
	}
}

func TestFloatListDefaultValueFail(t *testing.T) {
	testArgs := []string{"progname"}

	p := NewParser("progname", "Prog description")

	_ = p.FloatList("f", "float", &Options{Default: false})

	err := p.Parse(testArgs)

	// Should pass on failure
	failureMessage := "cannot use default type [bool] as value of pointer with type [*[]float64]"
	if err == nil || err.Error() != failureMessage {
		t.Errorf("Test %s failed: expected error [%s], got error [%+v]", t.Name(), failureMessage, err)
	}
}

func TestIntListDefaultValueFail(t *testing.T) {
	testArgs := []string{"progname"}

	p := NewParser("progname", "Prog description")

	_ = p.IntList("i", "int", &Options{Default: false})

	err := p.Parse(testArgs)

	// Should pass on failure
	failureMessage := "cannot use default type [bool] as value of pointer with type [*[]int]"
	if err == nil || err.Error() != failureMessage {
		t.Errorf("Test %s failed: expected error [%s], got error [%+v]", t.Name(), failureMessage, err)
	}
}

func TestStringListDefaultValueFail(t *testing.T) {
	testArgs := []string{"progname"}

	p := NewParser("progname", "Prog description")

	_ = p.StringList("s", "string", &Options{Default: false})

	err := p.Parse(testArgs)

	// Should pass on failure
	failureMessage := "cannot use default type [bool] as value of pointer with type [*[]string]"
	if err == nil || err.Error() != failureMessage {
		t.Errorf("Test %s failed: expected error [%s], got error [%+v]", t.Name(), failureMessage, err)
	}
}

func TestListDefaultValueFail(t *testing.T) {
	testArgs := []string{"progname"}

	p := NewParser("progname", "Prog description")

	_ = p.List("s", "string", &Options{Default: false})

	err := p.Parse(testArgs)

	// Should pass on failure
	if err == nil || err.Error() != "cannot use default type [bool] as value of pointer with type [*[]string]" {
		t.Errorf("Test %s failed: expected error [%s], got error [%+v]", t.Name(), "cannot use default type [bool] as value of pointer with type [*[]string]", err)
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
	if err == nil || err.Error() != "cannot use default type [bool] as value of pointer with type [*string]" {
		t.Errorf("Test %s failed: expected error [%s], got error [%+v]", t.Name(), "cannot use default type [bool] as value of pointer with type [*string]", err)
	}
}

func TestFloatAddArgumentFail(t *testing.T) {
	type testCase struct {
		testName, shortArg, longArg, failureMessage string
	}
	tt := []testCase{
		testCase{testName: "Long short name", shortArg: "ff", longArg: "flag2", failureMessage: "unable to add Float: short name must not exceed 1 character"},
		testCase{testName: "Long name not provided", shortArg: "f", longArg: "", failureMessage: "unable to add Float: long name should be provided"},
		testCase{testName: "Long name twice", shortArg: "f", longArg: "flag1", failureMessage: "unable to add Float: long name flag1 occurs more than once"},
		testCase{testName: "Short name twice", shortArg: "F", longArg: "flag2", failureMessage: "unable to add Float: short name F occurs more than once"},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					rezString := fmt.Sprintf("%v", r)
					if strings.Contains(rezString, tc.failureMessage) == false {
						t.Errorf("Test %s failed with panic result: \"%v\". panic result: %q expected", t.Name(), r, tc.failureMessage)
					}
				} else {
					t.Errorf("Test %s failed with no panic, but panic expected with result: %q", t.Name(), tc.failureMessage)
				}
			}()
			p := NewParser("", "description")
			_ = p.Float("F", "flag1", nil)
			_ = p.Float(tc.shortArg, tc.longArg, nil)
		})
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

func TestParserHelpFuncDefault(t *testing.T) {
	parser := NewParser("parser", "")
	if parser.HelpFunc == nil || parser.Help(nil) != parser.Usage(nil) {
		t.Errorf("HelpFunc should default to Usage function")
	}
}

func TestCommandHelpFuncDefault(t *testing.T) {
	parser := NewParser("parser", "")
	command := parser.NewCommand("command", "")
	if command.HelpFunc != nil {
		t.Errorf("HelpFunc should default to Usage function")
	}
}

func TestCommandHelpFuncDefaultToParent(t *testing.T) {
	parser := NewParser("parser", "")
	command := parser.NewCommand("command", "")

	parser.HelpFunc = func(c *Command, msg interface{}) string {
		return "testing"
	}

	if command.Help(nil) == command.Usage(nil) || command.Help(nil) != parser.Help(nil) {
		t.Errorf("command HelpFunc should default to parent function")
	}
}

func TestParserExitOnHelpTrue(t *testing.T) {
	exited := false
	exit = func(n int) {
		exited = true
	}

	parser := NewParser("parser", "")

	print = func(...interface{}) (int, error) {
		return 0, nil
	}

	if err := parser.Parse([]string{"parser", "-h"}); err == nil {
		if !exited {
			t.Errorf("Parsing help should have invoked os.Exit")
		}
	} else {
		t.Error(err)
	}
}

func TestParserExitOnHelpFalse(t *testing.T) {
	exited := false
	exit = func(n int) {
		exited = true
	}

	parser := NewParser("parser", "")
	parser.ExitOnHelp(false)

	print = func(...interface{}) (int, error) {
		return 0, nil
	}

	if err := parser.Parse([]string{"parser", "-h"}); exited {
		t.Errorf("Parsing help should not have invoked os.Exit")
	} else if err != nil {
		t.Error(err)
	}
}

func TestParserDisableHelp(t *testing.T) {
	parser := NewParser("parser", "")
	parser.DisableHelp()
	if len(parser.args) > 0 {
		t.Errorf("Parser should not have any arguments")
	}

	print = func(...interface{}) (int, error) {
		return 0, nil
	}

	if err := parser.Parse([]string{"parser", "-h"}); err == nil {
		t.Errorf("Parsing should fail, help argument shouldn't exist")
	}
}

func TestParserSetHelp(t *testing.T) {
	sname, lname := "x", "xyz"
	parser := NewParser("parser", "")
	parser.SetHelp(sname, lname)
	if len(parser.args) != 1 {
		t.Errorf("Parser should have one argument:\n%s", parser.Help(nil))
	}
	arg := parser.args[0]
	if _, ok := arg.result.(*help); !ok {
		t.Errorf("Argument should be %T, is %T", help{}, arg.result)
	}
	if arg.sname != sname {
		t.Errorf("Argument short name should be %s, is %s", sname, arg.sname)
	}
	if arg.lname != lname {
		t.Errorf("Argument long name should be %s, is %s", lname, arg.lname)
	}
}

func TestCommandExitOnHelpTrue(t *testing.T) {
	exited := false
	exit = func(n int) {
		exited = true
	}

	parser := NewParser("parser", "")
	parser.NewCommand("command", "")

	print = func(...interface{}) (int, error) {
		return 0, nil
	}

	if err := parser.Parse([]string{"parser", "command", "-h"}); exited {
		if err != nil {
			t.Error(err)
		}
	} else {
		t.Errorf("Parsing help should have invoked os.Exit")
	}
}

func TestCommandExitOnHelpFalse(t *testing.T) {
	exited := false
	exit = func(n int) {
		exited = true
	}

	parser := NewParser("parser", "")
	parser.NewCommand("command", "")
	parser.ExitOnHelp(false)

	print = func(...interface{}) (int, error) {
		return 0, nil
	}

	if err := parser.Parse([]string{"parser", "command", "-h"}); exited {
		t.Error("Parsing help should not have exited")
	} else if err != nil {
		t.Error(err)
	}
}

func TestCommandDisableHelp(t *testing.T) {
	parser := NewParser("parser", "")
	parser.NewCommand("command", "")
	parser.DisableHelp()
	if len(parser.args) > 0 {
		t.Errorf("Parser should not have any arguments")
	}

	print = func(...interface{}) (int, error) {
		return 0, nil
	}

	if err := parser.Parse([]string{"parser", "command", "-h"}); err == nil {
		t.Errorf("Parsing should fail, help argument shouldn't exist")
	}
}

func TestCommandHelpInheritance(t *testing.T) {
	parser := NewParser("parser", "")
	command := parser.NewCommand("command", "")
	parser.ExitOnHelp(false)

	if command.exitOnHelp != false {
		t.Errorf("Command should inherit exitOnHelp from parent, even after creation")
	}
}

func TestCommandHelpSetSnameOnly(t *testing.T) {
	parser := NewParser("parser", "")
	parser.SetHelp("q", "")

	arg := parser.args[0]

	_, ok := arg.result.(*help)
	if !ok {
		t.Error("Argument should be of help type")
	}

	if arg.sname != "h" || arg.lname != "help" {
		t.Error("Help arugment names should have defaulted")
	}
}
