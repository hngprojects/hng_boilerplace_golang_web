package utility

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
	"strings"
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

func FormatInspectionPeriod(t interface{}) string {
	timeStampStr, ok := t.(string)
	if !ok {
		return ""
	}

	timeStamp, err := strconv.Atoi(timeStampStr)
	if err != nil || timeStamp < 1 {
		return ""
	}

	inspectionTime := time.Unix(int64(timeStamp), 0)
	return inspectionTime.Format("2006-01-02 15:04:05")
}

func NumberFormat(t interface{}) float64 {
	num, ok := t.(float64)
	if !ok {
		numInt, ok := t.(int)
		if ok {
			num = float64(numInt)
		}
		return num
	}
	return num
}

func Add(num1, num2 interface{}) float64 {
	first, ok := num1.(float64)
	if !ok {
		firstInt, ok := num1.(int)
		if ok {
			first = float64(firstInt)
		}
	}
	second, ok := num2.(float64)
	if !ok {
		secondInt, ok := num1.(int)
		if ok {
			second = float64(secondInt)
		}
	}
	return first + second
}

func ConvertIntValues(m map[string]interface{}) {
	for key, value := range m {
		switch v := value.(type) {
		case float64:
			if intValue := int(v); float64(intValue) == v {
				m[key] = intValue
			}
		case map[string]interface{}:
			ConvertIntValues(v)
		}
	}
}

func StructToMap(obj interface{}) (map[string]interface{}, error) {

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return map[string]interface{}{}, err
	}

	result := make(map[string]interface{})

	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return map[string]interface{}{}, err
	}

	ConvertIntValues(result)

	return result, nil
}

func GetConstants(pkgImportPath string) (map[string]string, error) {
	// pkgImportPath example  ./services/names
	constants := map[string]string{}

	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, pkgImportPath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
					for _, spec := range genDecl.Specs {
						if valueSpec, ok := spec.(*ast.ValueSpec); ok {
							for i, ident := range valueSpec.Names {
								actionSpec, ok := valueSpec.Values[i].(*ast.BasicLit)
								if ok {
									constants[strings.ReplaceAll(ident.Name, `"`, "")] = strings.ReplaceAll(actionSpec.Value, `"`, "")
								}

							}
						}
					}
				}
			}
		}
	}

	return constants, nil
}
