package bindvalidate

import (
	"testing"

	"fmt"

	raven "github.com/getsentry/raven-go"
	"github.com/reactima/echo-bind-validate/validator"
)

func TestRunes(t *testing.T) {
	var tests = []struct {
		param    string
		expected bool
	}{
		{"달기&Co.", false},
		{"〩Hours", false},
	}
	for _, test := range tests {
		actual := govalidator.IsAlpha(test.param)
		if actual != test.expected {
			err := fmt.Errorf("Expected %v to be equal %v", actual, test.expected)
			res := raven.CaptureErrorAndWait(err, nil)
			fmt.Println(res)
			t.Errorf("Expected %v to be equal %v", actual, test.expected)
			t.Fail()
		}
	}
	// t.Log("ok")
}
