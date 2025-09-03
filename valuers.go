package gosqlgen

import (
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"time"
)

var (
	ErrValuer     = errors.New("failed to infer new value")
	ErrStringKind = errors.New("unrecognized string kind")
	ErrPrevType   = errors.New("type of previous value does not match valuer")
)

type valuerNumeric struct {
	max     int
	min     int
	isFloat bool
}

func NewValuerNumeric(minValue, maxValue int, isFloat bool) (valuerNumeric, error) {
	if maxValue == 0 && minValue == 0 {
		maxValue = 32
	}
	return valuerNumeric{max: maxValue, min: minValue, isFloat: isFloat}, nil
}

func (v valuerNumeric) randomInt(prev int) (int, error) {
	r := rand.Intn(v.max+1-v.min) + v.min
	if r != prev {
		return r, nil
	}

	if r == v.max {
		if prev != v.max-1 {
			return v.max - 1, nil
		}
		return v.max - 2, nil
	}
	if r == v.min {
		if prev != v.max+1 {
			return v.max + 1, nil
		}
		return v.max + 2, nil
	}

	return 0, ErrValuer
}

func (v valuerNumeric) New(prev any) (any, error) {
	var p int
	switch prev.(type) {
	case int, int8, int16, int32, int64:
		p = prev.(int)
	case float64, float32:
		p = int(prev.(float64))
	default:
		return 0, ErrPrevType
	}

	r, err := v.randomInt(p)
	if err != nil {
		return 0, err
	}

	if v.isFloat {
		return float64(r), nil
	}

	return r, nil
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

	return valuerString{length: length, kind: kind, charSet: charSet, valueSet: valueSet}, nil
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

	return "", nil
}

func (v valuerString) json(prev string) (string, error) {
	return "", nil
}

func (v valuerString) uuid(prev string) (string, error) {
	return "", nil
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
		return v.json(ps)
	case stringKindUUID:
		return v.uuid(ps)
	}

	return "", ErrStringKind
}

type valuerTime struct{}

func NewValuerTime() (valuerTime, error) {
	return valuerTime{}, nil
}

func (v valuerTime) New(prev any) (any, error) {
	return time.Now(), nil
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
