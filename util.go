package jsonschema

import (
	"math"
)

func getInteger(val interface{}) (result int, canconv bool) {
	val_f, ok := val.(float64)
	if !ok {
		// must number
		return 0, false
	}

	if math.Mod(val_f, 1) != 0 {
		// must integer, not float.
		return 0, false
	}

	if val_f < 0 {
		// must greater or equals to 0
		return 0, false
	}

	return int(val_f), true
}

func convInterfaceArrayToStringArray(val []interface{}) []string {
	ret := make([]string, 0)
	for _, v := range val {
		str, ok := v.(string)
		if !ok {
			return nil
		}

		ret = append(ret, str)
	}

	return ret
}
