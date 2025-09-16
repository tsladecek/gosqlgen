package gosqlgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewValuerNumeric(t *testing.T) {
	cases := []struct {
		name        string
		minValue    float64
		maxValue    float64
		isFloat     bool
		expected    valuerNumeric
		expectedErr error
	}{
		{name: "valid min and max not set", expected: valuerNumeric{max: 32}},
		{name: "valid min not set", maxValue: 16, expected: valuerNumeric{max: 16}},
		{name: "valid min and max set", minValue: 2, maxValue: 4, expected: valuerNumeric{max: 4, min: 2}},
		{name: "valid min and max set is float", minValue: 2, maxValue: 4, isFloat: true, expected: valuerNumeric{max: 4, min: 2, isFloat: true}},
		{name: "invalid min larger than max", minValue: 2, maxValue: 1, isFloat: true, expectedErr: ErrValuerConstructor},
		{name: "invalid min set but max not set", minValue: 2, maxValue: 1, isFloat: true, expectedErr: ErrValuerConstructor},
		{name: "invalid difference between min and max less than 1 for integers", minValue: 2, maxValue: 2, isFloat: false, expectedErr: ErrValuerConstructor},
		{name: "valid difference between min and max less than 1 for floats", minValue: 2, maxValue: 3, isFloat: true, expected: valuerNumeric{min: 2, max: 3, isFloat: true}},
		{name: "invalid difference between min and max less than 0.00001 for floats", minValue: 2, maxValue: 2.000000001, isFloat: true, expectedErr: ErrValuerConstructor},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewValuerNumeric(tt.minValue, tt.maxValue, tt.isFloat)
			require.Equal(t, tt.expectedErr == nil, err == nil)

			if tt.expectedErr == nil {
				assert.Equal(t, tt.expected, v)
			} else {
				assert.ErrorIs(t, err, tt.expectedErr)
			}
		})
	}
}

func TestValuerNumericNew(t *testing.T) {
	cases := []struct {
		name   string
		valuer valuerNumeric
		prev   TestValue
	}{
		{name: "int", valuer: valuerNumeric{max: 32}, prev: TestValue{Value: 2}},
		{name: "binary integer", valuer: valuerNumeric{max: 1}, prev: TestValue{Value: 0}},
		{name: "float", valuer: valuerNumeric{min: 0.05, max: 0.1, isFloat: true}, prev: TestValue{Value: 0.75}},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			for range 10 {
				newValue, err := tt.valuer.New(tt.prev)
				require.NoError(t, err)
				assert.NotEqual(t, tt.prev, newValue)
			}
		})
	}
}
