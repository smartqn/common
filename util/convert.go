package util

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func CoerceString(v interface{}) (string, error) {
	switch v := v.(type) {
	case string:
		return v, nil
	case int, int16, int32, int64, uint, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v), nil
	case float32, float64:
		return fmt.Sprintf("%f", v), nil
	}
	return fmt.Sprintf("%s", v), nil
}

func CoerceInt(v interface{}) (int, error) {
	switch v := v.(type) {
	case string:
		i64, err := strconv.ParseInt(v, 10, 0)
		return int(i64), err
	case int, int16, int32, int64:
		return int(reflect.ValueOf(v).Int()), nil
	case uint, uint16, uint32, uint64:
		return int(reflect.ValueOf(v).Uint()), nil
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	}
	return 0, errors.New("invalid value type")
}

func CoerceFloat(v interface{}) (float64, error) {
	switch v := v.(type) {
	case string:
		if len(v) > 0 {
			if iv, err := strconv.ParseFloat(v, 64); err == nil {
				return iv, nil
			}
		}
	case int, int16, int32, int64:
		return float64(reflect.ValueOf(v).Int()), nil
	case uint, uint16, uint32, uint64:
		return float64(reflect.ValueOf(v).Uint()), nil
	case float32:
		return float64(v), nil
	case float64:
		return float64(v), nil
	}
	return 0, errors.New("invalid value type")
}
