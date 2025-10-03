package gosqlgen

import (
	"fmt"
	"slices"
	"time"
)

type valuerNumeric struct {
	max     float64
	min     float64
	isFloat bool
}

func NewValuerNumeric(minValue, maxValue float64, isFloat bool) (valuerNumeric, error) {
	if minValue > maxValue {
		return valuerNumeric{}, Errorf("min value is greater than max value: %w", ErrValuerConstructor)
	}

	if maxValue == 0 && minValue == 0 {
		maxValue = 32
	} else if minValue == maxValue {
		return valuerNumeric{}, Errorf("min can not equal max value: %w", ErrValuerConstructor)
	}

	if maxValue-minValue < 1 && !isFloat {
		return valuerNumeric{}, Errorf("difference between min and max must be greater than one for integers: %w", ErrValuerConstructor)
	}

	// float32 has 6-9 significant bit precision
	minFloatDiff := 0.00001
	if maxValue-minValue < minFloatDiff && isFloat {
		return valuerNumeric{}, Errorf("difference between min and max must be greater than %f: %w", minFloatDiff, ErrValuerConstructor)
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
			return TestValue{}, Errorf("previous value not an int: %w", ErrValuer)
		}
		iv, err := v.otherInt(p)
		if err != nil {
			return TestValue{}, Errorf("when generating new integer value: %w", ErrValuer)
		}

		return TestValue{Value: iv}, nil
	}

	p, ok := prev.Value.(float64)
	if !ok {
		return TestValue{}, Errorf("previous value not a float: %w", ErrValuer)
	}
	fv, err := v.otherFloat(p)
	if err != nil {
		return TestValue{}, Errorf("when generating new integer value: %w", ErrValuer)
	}
	return TestValue{Value: fv}, nil
}

func (v valuerNumeric) Zero() TestValue {
	if v.isFloat {
		return TestValue{v.min}
	}
	return TestValue{int(v.min)}
}

type valuerString struct {
	length   int
	kind     stringKind
	charSet  []rune
	valueSet []string
}

func NewValuerString(length int, kind stringKind, charSet []rune, valueSet []string) (valuerString, error) {
	if !slices.Contains([]stringKind{stringKindBasic, stringKindJSON, stringKindUUID, stringKindEnum}, kind) {
		return valuerString{}, Errorf("kind=%s: %w", kind, ErrStringKind)
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
	if prev == "" {
		return TestValue{Value: v.randomString(v.length)}, nil
	}

	out := make([]rune, v.length)
	for i := range v.length {

		charSetWithoutChar := make([]rune, 0)
		if i < len(prev) {
			idx := 0
			for _, j := range v.charSet {
				if j != rune(prev[i]) {
					charSetWithoutChar = append(charSetWithoutChar, v.charSet[idx])
					idx++
				}
			}
		} else {
			charSetWithoutChar = v.charSet
		}

		if len(charSetWithoutChar) == 0 {
			return TestValue{}, Errorf("can not infer new basic string value: %w", ErrValuer)
		}

		out[i] = charSetWithoutChar[RandomInt(len(charSetWithoutChar))]
	}

	return TestValue{Value: string(out)}, nil
}

func (v valuerString) enum(prev string) (TestValue, error) {
	if len(v.valueSet) == 1 {
		return TestValue{}, Errorf("can not infer new value since the value set contains only one item: %w", ErrValuer)
	}

	if v.valueSet[0] == prev {
		return TestValue{Value: v.valueSet[1]}, nil
	}
	return TestValue{v.valueSet[0]}, nil
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
	case stringKindEnum:
		return v.enum(ps)
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
