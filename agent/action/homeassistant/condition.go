package homeassistant

import (
	"errors"
	"fmt"
	"strconv"
)

type haCondition struct {
	valueStr   string `yaml:"value"`
	comparison int
	valueNum   int64
}

type haConditionYaml struct {
	Comparison string `yaml:"comparison"`
	Value      string `yaml:"value"`
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
	c.comparison, ok = comparisonEnumMap[v.Comparison]
	if !ok {
		return fmt.Errorf("unknown comparison: %s", v.Comparison)
	}

	if !isNumericComparison(c.comparison) {
		c.valueStr = v.Value
		return nil
	}

	c.valueNum, err = strconv.ParseInt(v.Value, 10, 64)
	return err
}

func (c *haCondition) Evaluate(stateStr string) (bool, error) {
	if isNumericComparison(c.comparison) {
		stateNum, err := strconv.ParseInt(stateStr, 10, 64)
		if err != nil {
			return false, err
		}
		switch c.comparison {
		case compareLessThan:
			return stateNum < c.valueNum, nil
		case compareLessThanOrEqual:
			return stateNum <= c.valueNum, nil
		case compareGreaterThan:
			return stateNum > c.valueNum, nil
		case compareGreaterThanOrEqual:
			return stateNum >= c.valueNum, nil
		default:
			return false, errors.New("invalid number comparison, this is a bug")
		}
	}

	switch c.comparison {
	case compareEquals:
		return stateStr == c.valueStr, nil
	case compareNotEquals:
		return stateStr != c.valueStr, nil
	default:
		return false, errors.New("invalid string comparison, this is a bug")
	}
}
