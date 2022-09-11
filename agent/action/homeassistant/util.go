package homeassistant

import "strconv"

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
