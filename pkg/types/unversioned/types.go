package unversioned

import (
	"fmt"
	"strings"
	"time"
)

type TestResult struct {
	Name     string
	Pass     bool
	Stdout   string   `json:",omitempty"`
	Stderr   string   `json:",omitempty"`
	Errors   []string `json:",omitempty"`
	Duration time.Duration
}

func (t *TestResult) String() string {
	strRepr := fmt.Sprintf("\nTest Name:%s", t.Name)
	testStatus := "Fail"
	if t.IsPass() {
		testStatus = "Pass"
	}
	strRepr += fmt.Sprintf("\nTest Status:%s", testStatus)
	if t.Stdout != "" {
		strRepr += fmt.Sprintf("\nStdout:%s", t.Stdout)
	}
	if t.Stderr != "" {
		strRepr += fmt.Sprintf("\nStderr:%s", t.Stderr)
	}
	strRepr += fmt.Sprintf("\nErrors:%s\n", strings.Join(t.Errors, ","))
	strRepr += fmt.Sprintf("\nDuration:%s\n", t.Duration.String())
	return strRepr
}

func (t *TestResult) Error(s string) {
	t.Errors = append(t.Errors, s)
}

func (t *TestResult) Errorf(s string, args ...interface{}) {
	t.Errors = append(t.Errors, fmt.Sprintf(s, args...))
}

func (t *TestResult) Fail() {
	t.Pass = false
}

func (t *TestResult) IsPass() bool {
	return t.Pass
}

type SummaryObject struct {
	Pass     int
	Fail     int
	Total    int
	Duration time.Duration
	Results  []*TestResult `json:",omitempty"`
}
