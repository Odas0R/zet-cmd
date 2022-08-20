package assert

import (
	"fmt"
	"testing"

	"github.com/odas0r/zet/cmd/columnize"
)

func Equal(t *testing.T, got interface{}, expected interface{}, message string) {
	assert(t, got, expected, message, true)
}

func NotEqual(t *testing.T, got interface{}, expected interface{}, message string) {
	assert(t, got, expected, message, false)
}

func assert(t *testing.T, got interface{}, expected interface{}, message string, expectation bool) {
	errorMessage := []string{
    "\n",
		fmt.Sprintf("message: | \"%s\"", message),
		fmt.Sprintf("expected: | %s", expected),
		fmt.Sprintf("got: | %s", got),
	}

	formattedErrorMessage := columnize.SimpleFormat(errorMessage)

	if (expected == got) != expectation {
		t.Errorf(formattedErrorMessage)
	}
}
