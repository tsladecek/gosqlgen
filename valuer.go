package gosqlgen

import (
	"errors"
	"fmt"
	"math/rand"
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

func (v valuerNumeric) randomInt(prev int) (int, error) {
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

func (v valuerNumeric) randomFloat(prev float64) (float64, error) {
	step := (v.max - v.min) / 10

	if prev == v.max {
		return prev - step, nil
	}

	if prev == v.min {
		return prev + step, nil
	}
	return v.min, nil
}

func (v valuerNumeric) New(prev any) (any, error) {
	if !v.isFloat {
		p, ok := prev.(int)
		if !ok {
			return 0, fmt.Errorf("%w: previous value not an int", ErrValuer)
		}
		return v.randomInt(p)
	}

	p, ok := prev.(float64)
	if !ok {
		return 0, fmt.Errorf("%w: previous value not a float", ErrValuer)
	}
	return v.randomFloat(p)
}

func (v valuerNumeric) Zero() any {
	if v.isFloat {
		return v.min
	}
	return int(v.min)
}

func (v valuerNumeric) Format(value any, typ string) (string, error) {
	switch value.(type) {
	case int, float64:
		switch typ {
		case "database/sql.NullInt16":
			return fmt.Sprintf("sql.NullInt16{Valid: true, Int16: %d}", value), nil
		case "database/sql.NullInt32":
			return fmt.Sprintf("sql.NullInt32{Valid: true, Int32: %d}", value), nil
		case "database/sql.NullInt64":
			return fmt.Sprintf("sql.NullInt16{Valid: true, Int16: %d}", value), nil
		case "database/sql.NullFloat64":
			return fmt.Sprintf("sql.NullFloat64{Valid: true, Float64: %d}", value), nil
		default:
			return fmt.Sprintf("%v", value), nil
		}
	default:
		return "", fmt.Errorf("value not a valid numeric type")
	}
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

func (v valuerString) basic(prev string) (string, error) {
	if len(v.valueSet) > 0 {
		if len(v.valueSet) == 1 {
			return "", fmt.Errorf("%w: can not infer new value since the value set contains only one item", ErrValuer)
		}

		if v.valueSet[0] == prev {
			return v.valueSet[1], nil
		}
		return v.valueSet[0], nil
	}

	if prev == "" {
		out := make([]rune, v.length)
		for i := range v.length {
			out[i] = v.charSet[rand.Intn(len(v.charSet))]
		}

		return string(out), nil
	}

	for _, c := range v.charSet {
		if rune(prev[0]) != c {
			out := []rune(prev)
			out[0] = c
			return string(out), nil
		}
	}

	return "", fmt.Errorf("%w: can not infer new basic string value", ErrValuer)
}

func (v valuerString) randomString(length int) string {
	out := make([]rune, length)
	for i := range length {
		out[i] = v.charSet[rand.Intn(len(v.charSet))]
	}

	return string(out)
}

func (v valuerString) json() (string, error) {
	return fmt.Sprintf(`{"%s":"%s", "%s":"%s"}`, v.randomString(8), v.randomString(8), v.randomString(8), v.randomString(8)), nil
}

func (v valuerString) uuid() (string, error) {
	return fmt.Sprintf("%s-%s-4%s-9%s-%s", v.randomString(8), v.randomString(4), v.randomString(3), v.randomString(3), v.randomString(12)), nil
}

func (v valuerString) New(prev any) (any, error) {
	ps, ok := prev.(string)

	if !ok {
		return "", ErrPrevType
	}

	switch v.kind {
	case stringKindBasic:
		return v.basic(ps)
	case stringKindJSON:
		return v.json()
	case stringKindUUID:
		return v.uuid()
	}

	return "", ErrStringKind
}

func (v valuerString) Zero() any {
	return v.randomString(v.length)
}

func (v valuerString) Format(value any, typ string) (string, error) {
	vv, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("%w: value=%v, type=%s", ErrValueType, value, typ)
	}

	switch typ {
	case "encoding/json.RawMessage", "[]byte":
		return fmt.Sprintf("[]byte(`%s`)", vv), nil
	case "sql.NullString":
		return fmt.Sprintf("sql.NullString{Valid: true, String: \"%s\"}", vv), nil
	default:
		return fmt.Sprintf(`"%s"`, vv), nil
	}
}

type valuerTime struct{}

func NewValuerTime() (valuerTime, error) {
	return valuerTime{}, nil
}

func (v valuerTime) New(prev any) (any, error) {
	return time.Now(), nil
}

func (v valuerTime) Zero() any {
	return time.Now()
}

func (v valuerTime) Format(value any, typ string) (string, error) {
	switch typ {
	case "time.Time":
		return "time.Now()", nil
	case "database/sql.NullTime":
		return "sql.NullTime{Valid: true, Time: time.Now()}", nil
	}
	return "", fmt.Errorf("%w: unrecognized type %s", ErrValueFormat, typ)
}

type valuerBoolean struct{}

func NewValuerBoolean() (valuerBoolean, error) {
	return valuerBoolean{}, nil
}

func (v valuerBoolean) New(prev any) (any, error) {
	p, ok := prev.(bool)

	if !ok {
		return false, ErrPrevType
	}
	return !p, nil
}

func (v valuerBoolean) Zero() any {
	return false
}

func (v valuerBoolean) Format(value any, typ string) (string, error) {
	vv, ok := value.(bool)
	if !ok {
		return "", fmt.Errorf("%w: value=%v, type=%s", ErrValueType, value, typ)
	}

	switch typ {
	case "bool":
		return fmt.Sprintf("%t", vv), nil
	case "database/sql.NullBool":
		return fmt.Sprintf("sql.NullBool{Valid: true, Bool: %t}", vv), nil
	}
	return "", fmt.Errorf("%w: unrecognized type %s", ErrValueFormat, typ)
}
