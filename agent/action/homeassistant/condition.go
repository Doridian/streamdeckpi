package homeassistant

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Doridian/go-haws"
)

type haCondition struct {
	src        string
	comparison int
	valueRaw   interface{} `yaml:"value"`
	valueNum   float64
}

type haConditionYaml struct {
	Src        string      `yaml:"src"`
	Comparison string      `yaml:"comparison"`
	Value      interface{} `yaml:"value"`
}

const (
	compareEquals = iota
	compareNotEquals
	compareLessThan
	compareLessThanOrEqual
	compareGreaterThanOrEqual
	compareGreaterThan
)

var comparisonEnumMap = map[string]int{
	"==": compareEquals,
	"=":  compareEquals,
	"eq": compareEquals,
	"is": compareEquals,

	"!=":  compareNotEquals,
	"<>":  compareNotEquals,
	"ne":  compareNotEquals,
	"not": compareNotEquals,

	"<":  compareLessThan,
	"lt": compareLessThan,

	"<=":  compareLessThanOrEqual,
	"lte": compareLessThanOrEqual,

	">":  compareGreaterThan,
	"gt": compareGreaterThan,

	">=":  compareGreaterThanOrEqual,
	"gte": compareGreaterThanOrEqual,
}

func isNumericComparison(comp int) bool {
	return comp == compareLessThan || comp == compareLessThanOrEqual || comp == compareGreaterThan || comp == compareGreaterThanOrEqual
}

func (c *haCondition) UnmarshalYAML(unmarshal func(interface{}) error) error {
	v := &haConditionYaml{}
	err := unmarshal(v)
	if err != nil {
		return err
	}

	var ok bool

	c.src = v.Src
	if c.src == "" {
		c.src = "state"
	}

	c.comparison, ok = comparisonEnumMap[v.Comparison]
	if !ok {
		return fmt.Errorf("unknown comparison: %s", v.Comparison)
	}

	if !isNumericComparison(c.comparison) {
		c.valueRaw = v.Value
		return nil
	}

	c.valueNum, ok = v.Value.(float64)
	if ok {
		return nil
	}

	valInt, ok := v.Value.(int)
	if ok {
		c.valueNum = float64(valInt)
		return nil
	}

	return errors.New("attempt to create number comparison with non-number")
}

func coerceNumber(val interface{}) (float64, bool) {
	valNum, ok := val.(float64)
	if ok {
		return valNum, true
	}

	valInt, ok := val.(int)
	if ok {
		return float64(valInt), true
	}

	valStr, ok := val.(string)
	if !ok {
		return 0, false
	}

	valNum, err := strconv.ParseFloat(valStr, 64)
	if err == nil {
		return 0, false
	}

	return valNum, true
}

func (c *haCondition) Evaluate(state *haws.State) (bool, error) {
	var val interface{}
	if c.src == "state" {
		val = state.State
	} else if state.Attributes != nil {
		val = state.Attributes[c.src]
	}

	if isNumericComparison(c.comparison) {
		var valNum float64
		if val == nil {
			valNum = 0
		} else {
			var ok bool
			valNum, ok = coerceNumber(val)
			if !ok {
				return false, errors.New("attempt to use number comparison on non-number")
			}
		}

		switch c.comparison {
		case compareLessThan:
			return valNum < c.valueNum, nil
		case compareLessThanOrEqual:
			return valNum <= c.valueNum, nil
		case compareGreaterThan:
			return valNum > c.valueNum, nil
		case compareGreaterThanOrEqual:
			return valNum >= c.valueNum, nil
		default:
			return false, errors.New("invalid number comparison, this is a bug")
		}
	}

	switch c.comparison {
	case compareEquals:
		return val == c.valueRaw, nil
	case compareNotEquals:
		return val != c.valueRaw, nil
	default:
		return false, errors.New("invalid generic comparison, this is a bug")
	}
}
