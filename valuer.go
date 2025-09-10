package gosqlgen

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

var (
	ErrValuer            = errors.New("failed to infer new value")
	ErrStringKind        = errors.New("unrecognized string kind")
	ErrPrevType          = errors.New("type of previous value does not match valuer")
	ErrValuerConstructor = errors.New("failed to construct valuer")
	ErrValueFormat       = errors.New("failed to format value")
	ErrValueType         = errors.New("invalid type")
)

type valuerNumeric struct {
	max     float64
	min     float64
	isFloat bool
}

func NewValuerNumeric(minValue, maxValue float64, isFloat bool) (valuerNumeric, error) {
	if minValue > maxValue {
		return valuerNumeric{}, fmt.Errorf("%w: min value is greater than max value", ErrValuerConstructor)
	}

	if maxValue == 0 && minValue == 0 {
		maxValue = 32
	} else if minValue == maxValue {
		return valuerNumeric{}, fmt.Errorf("%w: min can not equal max value", ErrValuerConstructor)
	}

	if maxValue-minValue < 1 && !isFloat {
		return valuerNumeric{}, fmt.Errorf("%w: difference between min and max must be greater than one for integers", ErrValuerConstructor)
	}

	// float32 has 6-9 significant bit precision
	minFloatDiff := 0.00001
	if maxValue-minValue < minFloatDiff && isFloat {
		return valuerNumeric{}, fmt.Errorf("%w: difference between min and max must be greater %f", ErrValuerConstructor, minFloatDiff)
	}

	return valuerNumeric{max: maxValue, min: minValue, isFloat: isFloat}, nil
}

func (v valuerNumeric) otherInt(prev int) (int, error) {
	maxInt := int(v.max)
	minInt := int(v.min)

	if prev == maxInt {
		return prev - 1, nil
	}

	if prev == minInt {
		return prev + 1, nil
	}

	return minInt, nil
}

func (v valuerNumeric) otherFloat(prev float64) (float64, error) {
	step := (v.max - v.min) / 10

	if prev == v.max {
		return prev - step, nil
	}

	if prev == v.min {
		return prev + step, nil
	}
	return v.min, nil
}

func (v valuerNumeric) New(prev TestValue) (TestValue, error) {
	if !v.isFloat {
		p, ok := prev.Value.(int)
		if !ok {
			return TestValue{}, fmt.Errorf("%w: previous value not an int", ErrValuer)
		}
		iv, err := v.otherInt(p)
		if err != nil {
			return TestValue{}, fmt.Errorf("%w: when generating new integer value", ErrValuer)
		}

		return TestValue{Value: iv}, nil
	}

	p, ok := prev.Value.(float64)
	if !ok {
		return TestValue{}, fmt.Errorf("%w: previous value not a float", ErrValuer)
	}
	fv, err := v.otherFloat(p)
	if err != nil {
		return TestValue{}, fmt.Errorf("%w: when generating new integer value", ErrValuer)
	}
	return TestValue{Value: fv}, nil
}

func (v valuerNumeric) Zero() TestValue {
	if v.isFloat {
		return TestValue{v.min}
	}
	return TestValue{int(v.min)}
}

type stringKind string

const (
	stringKindBasic stringKind = "basic"
	stringKindJSON  stringKind = "json"
	stringKindUUID  stringKind = "UUID"
)

type valuerString struct {
	length   int
	kind     stringKind
	charSet  []rune
	valueSet []string
}

func NewValuerString(length int, kind stringKind, charSet []rune, valueSet []string) (valuerString, error) {
	if !slices.Contains([]stringKind{stringKindBasic, stringKindJSON, stringKindUUID}, kind) {
		return valuerString{}, fmt.Errorf("%w: kind=%s", ErrStringKind, kind)
	}

	if len(charSet) == 0 {
		charSet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	}

	maxLength := 32
	if length > 0 {
		maxLength = min(length, maxLength)
	}

	return valuerString{length: maxLength, kind: kind, charSet: charSet, valueSet: valueSet}, nil
}

func (v valuerString) basic(prev string) (TestValue, error) {
	if len(v.valueSet) > 0 {
		if len(v.valueSet) == 1 {
			return TestValue{}, fmt.Errorf("%w: can not infer new value since the value set contains only one item", ErrValuer)
		}

		if v.valueSet[0] == prev {
			return TestValue{Value: v.valueSet[1]}, nil
		}
		return TestValue{v.valueSet[0]}, nil
	}

	if prev == "" {
		return TestValue{Value: v.randomString(v.length)}, nil
	}

	for _, c := range v.charSet {
		if rune(prev[0]) != c {
			out := []rune(prev)
			out[0] = c
			return TestValue{Value: string(out)}, nil
		}
	}

	return TestValue{}, fmt.Errorf("%w: can not infer new basic string value", ErrValuer)
}

func (v valuerString) randomString(length int) string {
	return RandomString(length, v.charSet)
}

func (v valuerString) json() (TestValue, error) {
	return TestValue{Value: fmt.Sprintf(`{"%s":"%s", "%s":"%s"}`, v.randomString(8), v.randomString(8), v.randomString(8), v.randomString(8))}, nil
}

func (v valuerString) uuid() (TestValue, error) {
	return TestValue{Value: fmt.Sprintf("%s-%s-4%s-9%s-%s", v.randomString(8), v.randomString(4), v.randomString(3), v.randomString(3), v.randomString(12))}, nil
}

func (v valuerString) New(prev TestValue) (TestValue, error) {
	ps, ok := prev.Value.(string)

	if !ok {
		return TestValue{}, ErrPrevType
	}

	switch v.kind {
	case stringKindBasic:
		return v.basic(ps)
	case stringKindJSON:
		return v.json()
	case stringKindUUID:
		return v.uuid()
	}

	return TestValue{}, ErrStringKind
}

func (v valuerString) Zero() TestValue {
	return TestValue{Value: v.randomString(v.length)}
}

type valuerTime struct{}

func NewValuerTime() (valuerTime, error) {
	return valuerTime{}, nil
}

func (v valuerTime) New(prev TestValue) (TestValue, error) {
	return TestValue{Value: time.Now()}, nil
}

func (v valuerTime) Zero() TestValue {
	return TestValue{Value: time.Now()}
}

type valuerBoolean struct{}

func NewValuerBoolean() (valuerBoolean, error) {
	return valuerBoolean{}, nil
}

func (v valuerBoolean) New(prev TestValue) (TestValue, error) {
	p, ok := prev.Value.(bool)

	if !ok {
		return TestValue{}, ErrPrevType
	}
	return TestValue{Value: !p}, nil
}

func (v valuerBoolean) Zero() TestValue {
	return TestValue{Value: false}
}
