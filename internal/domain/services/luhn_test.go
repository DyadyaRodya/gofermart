package services

import (
	"log"
	"testing"
)

func TestLuhnValidate(t *testing.T) {
	tests := []struct {
		number   string
		expectOK bool
	}{
		{"1234567812345670", true},
		{"1111222233334444", true},
		{"1111222233334441", false},
		{"49927398716", true},
		{"49927398717", false},
		{"1234567812345678", false},
		{"79927398710", false},
		{"79927398711", false},
		{"79927398712", false},
		{"79927398713", true},
		{"79927398714", false},
		{"79927398715", false},
		{"79927398716", false},
		{"79927398717", false},
		{"79927398718", false},
		{"79927398719", false},
		{"374652346956782346957823694857692364857368475368", true},
		{"374652346956782346957823694857692364857387456834", false},
		{"8", false},
		{"0", true},
	}
	ls := &LuhnDomainService{}
	for _, test := range tests {
		t.Run(test.number, func(t *testing.T) {
			res := ls.Validate(test.number)
			if res != test.expectOK {
				log.Printf("Expected `%t` but luhn check `%t` for `%s`.", test.expectOK, res, test.number)
				t.Fail()
			}
		})
	}
}
