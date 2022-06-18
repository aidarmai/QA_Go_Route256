package test

import (
	"fmt"
	"github.com/ozonmp/act-device-api/greeter"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGreet(t *testing.T) {
	name := "terri"
	a := "Good morning" + " " + strings.Title(name) + "!"
	b := "Good evening" + " " + strings.Title(name) + "!"
	c := "Good night" + " " + strings.Title(name) + "!"
	d := "Hello" + " " + strings.Title(name) + "!"
	var testTableGreet = []struct {
		inputValuesName string
		inputValuesHour int
		expectedResult  string
	}{
		{inputValuesName: name, inputValuesHour: -1, expectedResult: d},
		{inputValuesName: name, inputValuesHour: 0, expectedResult: c},
		{inputValuesName: name, inputValuesHour: 1, expectedResult: c},
		{inputValuesName: name, inputValuesHour: 5, expectedResult: c},
		{inputValuesName: name, inputValuesHour: 6, expectedResult: a},
		{inputValuesName: name, inputValuesHour: 7, expectedResult: a},
		{inputValuesName: name, inputValuesHour: 11, expectedResult: a},
		{inputValuesName: name, inputValuesHour: 12, expectedResult: d},
		{inputValuesName: name, inputValuesHour: 13, expectedResult: d},
		{inputValuesName: name, inputValuesHour: 17, expectedResult: d},
		{inputValuesName: name, inputValuesHour: 18, expectedResult: b},
		{inputValuesName: name, inputValuesHour: 19, expectedResult: b},
		{inputValuesName: name, inputValuesHour: 21, expectedResult: b},
		{inputValuesName: name, inputValuesHour: 22, expectedResult: c},
		{inputValuesName: name, inputValuesHour: 23, expectedResult: c},
		{inputValuesName: name, inputValuesHour: 24, expectedResult: c},
		{inputValuesName: name, inputValuesHour: 25, expectedResult: d},
	}
	for i, pair := range testTableGreet {
		testName := fmt.Sprintf("%d: %d %s", i+1, pair.inputValuesHour, pair.expectedResult)
		v := greeter.Greet(pair.inputValuesName, pair.inputValuesHour)
		t.Run(testName, func(t *testing.T) {
			assert.Equal(t, pair.expectedResult, v, "want: %v, got: %v", pair.expectedResult, v)
		})
	}
}
