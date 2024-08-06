package utility

import (
	"reflect"
	"strconv"
	"time"
)

func FormatDate(date, currentISOFormat, newISOFormat string) (string, error) {
	t, err := time.Parse(currentISOFormat, date)
	if err != nil {
		return date, err
	}
	return t.Format(newISOFormat), nil
}

func GetUnixTime(date, currentISOFormat, newISOFormat string) (int, error) {
	t, err := time.Parse(currentISOFormat, date)
	if err != nil {
		return 0, err
	}
	return int(t.Unix()), nil
}
func GetUnixString(date, currentISOFormat, newISOFormat string) (string, error) {
	t, err := time.Parse(currentISOFormat, date)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(t.Unix())), nil
}

func ConvertStringInterfaceToStringFloat(originalMap map[string]interface{}) map[string]float64 {
	convertedMap := make(map[string]float64)
	for key, value := range originalMap {
		if val, ok := value.(float64); ok {
			convertedMap[key] = val
		} else if val, ok := value.(string); ok {
			if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
				convertedMap[key] = floatVal
			}
		}
	}
	return convertedMap
}

func RemoveKey(p interface{}, key string) {
	val := reflect.ValueOf(p).Elem()
	val.FieldByName(key).Set(reflect.Zero(val.FieldByName(key).Type()))
}

func CopyStruct(src, dst interface{}) {
	srcValue := reflect.ValueOf(src).Elem()
	dstValue := reflect.ValueOf(dst).Elem()

	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)
		dstField := dstValue.FieldByName(srcValue.Type().Field(i).Name)

		if dstField.IsValid() {
			dstField.Set(srcField)
		}
	}
}
