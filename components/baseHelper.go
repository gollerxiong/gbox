package components

import (
	"reflect"
	"strings"
)

func StructToMap(s interface{}) map[string]interface{} {
	structValue := reflect.ValueOf(s)
	structValue = reflect.Indirect(reflect.ValueOf(s))
	kd := structValue.Kind()
	if kd == reflect.Interface {
		structValue = structValue.Elem()
	}

	if kd == reflect.Ptr {
		if structValue.IsNil() {
			return make(map[string]interface{})
		}
		structValue = structValue.Elem()
	}

	structFieldsCount := structValue.NumField()
	result := make(map[string]interface{})

	for i := 0; i < structFieldsCount; i++ {
		fieldName := structValue.Type().Field(i).Name
		jsonTag := structValue.Type().Field(i).Tag.Get("json")
		if jsonTag != "" {
			fieldName = jsonTag
		} else {
			continue
		}

		if fieldName == "-" {
			continue
		}

		if structValue.Field(i).Kind() == reflect.Struct {
			result[fieldName] = StructToMap(structValue.Field(i).Interface())
		} else {
			result[fieldName] = structValue.Field(i).Interface()
		}
	}

	return result
}

// 合并两个map成一个map，相同key的话，后面的map会覆盖前面的map的key
func MapMerge(map1 map[string]interface{}, map2 map[string]interface{}) map[string]interface{} {

	// 合并两个map
	for key, value := range map2 {
		map1[key] = value
	}

	return map1
}

// 基于第一个map，返回第二个map跟第一个map不相同的值
func DiffMapBaseFirst(mp1, mp2 map[string]interface{}) map[string]interface{} {
	diff := make(map[string]interface{})

	for key, val1 := range mp1 {
		val2, ok := mp2[key]

		// 如果元素为切片或map直接跳过不对这种元素类型做比较
		if ok {
			kind := reflect.TypeOf(val1).Kind()
			if kind == reflect.Map {
				continue
			}

			if kind == reflect.Slice || kind == reflect.Array {
				continue
			}
		}
		if ok && val2 != val1 {
			diff[key] = val2
		}
	}

	return diff
}

func MapKeys(mp map[string]interface{}) []string {
	keys := []string{}

	for key, _ := range mp {
		keys = append(keys, key)
	}

	return keys
}

func PickFieldsFromMap(mp map[string]interface{}, fields string) map[string]interface{} {
	result := make(map[string]interface{})

	if fields != "*" {
		fieldArr := strings.Split(fields, ",")

		for _, key := range fieldArr {
			val, ok := mp[key]

			if ok {
				result[key] = val
			}
		}
	} else {
		result = mp
	}

	return result
}
