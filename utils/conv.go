package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// convert interface to string.
func GetString(v interface{}) string {
	switch result := v.(type) {
	case string:
		return result
	case []byte:
		return string(result)
	default:
		if v != nil {
			return fmt.Sprintf("%v", result)
		}
	}
	return ""
}

// convert interface to int.
func GetInt(v interface{}) int {
	switch result := v.(type) {
	case int:
		return result
	case int32:
		return int(result)
	case int64:
		return int(result)
	case float32:
		return int(result)
	case float64:
		return int(result)
	default:
		if d := GetString(v); d != "" {
			value, _ := strconv.Atoi(d)
			return value
		}
	}
	return 0
}

// convert interface to int64.
func GetInt64(v interface{}) int64 {
	switch result := v.(type) {
	case int:
		return int64(result)
	case int32:
		return int64(result)
	case int64:
		return result
	case float64:
		return int64(GetFloat64(v))
	default:
		if d := GetString(v); d != "" {
			value, _ := strconv.ParseInt(d, 10, 64)
			return value
		}
	}
	return 0
}

// convert interface to float64.
func GetFloat64(v interface{}) float64 {
	switch result := v.(type) {
	case float64:
		return result
	default:
		if d := GetString(v); d != "" {
			value, _ := strconv.ParseFloat(d, 64)
			return value
		}
	}
	return 0
}

// convert interface to bool.
func GetBool(v interface{}) bool {
	switch result := v.(type) {
	case bool:
		return result
	default:
		if d := GetString(v); d != "" {
			value, _ := strconv.ParseBool(d)
			return value
		}
	}
	return false
}

//struct转换为map
func GetStructMap(v interface{}) map[string]interface{} {
	t := reflect.TypeOf(v)
	vo := reflect.ValueOf(v)
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[strings.ToLower(t.Field(i).Name)] = vo.Field(i).Interface()
	}
	return data
}

// convert interface to byte slice.
func GetByteArray(v interface{}) []byte {
	switch result := v.(type) {
	case []byte:
		return result
	case string:
		return []byte(result)
	default:
		return nil
	}
}
