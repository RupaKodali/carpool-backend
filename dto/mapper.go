package dto

// import (
// 	"reflect"
// )

// func MapToDTOs[T any, R any](models []T) ([]R, error) {
// 	var results []R
// 	for _, model := range models {
// 		dto, err := mapStruct[T, R](model)
// 		if err != nil {
// 			return nil, err
// 		}
// 		results = append(results, dto)
// 	}
// 	return results, nil
// }

// func mapStruct[T any, R any](src T) (R, error) {
// 	var dst R

// 	srcVal := reflect.ValueOf(src)
// 	if srcVal.Kind() == reflect.Ptr {
// 		srcVal = srcVal.Elem()
// 	}
// 	dstVal := reflect.ValueOf(&dst).Elem()

// 	srcType := srcVal.Type()

// 	for i := 0; i < srcVal.NumField(); i++ {
// 		srcField := srcVal.Field(i)
// 		srcFieldType := srcType.Field(i)

// 		// Recursively map embedded fields (anonymous struct)
// 		if srcFieldType.Anonymous {
// 			err := mapEmbeddedFields(srcField, dstVal)
// 			if err != nil {
// 				return dst, err
// 			}
// 			continue
// 		}

// 		dstField := dstVal.FieldByName(srcFieldType.Name)
// 		if !dstField.IsValid() || !dstField.CanSet() {
// 			continue
// 		}

// 		if dstField.Type() == srcField.Type() {
// 			dstField.Set(srcField)
// 		} else if srcField.Kind() == reflect.Struct && dstField.Kind() == reflect.Struct {
// 			mappedNested, err := mapSimpleStruct(srcField, dstField.Type())
// 			if err != nil {
// 				return dst, err
// 			}
// 			dstField.Set(mappedNested)
// 		}
// 	}

// 	return dst, nil
// }

// func mapEmbeddedFields(src reflect.Value, dst reflect.Value) error {
// 	srcType := src.Type()
// 	for i := 0; i < src.NumField(); i++ {
// 		srcField := src.Field(i)
// 		srcFieldType := srcType.Field(i)

// 		dstField := dst.FieldByName(srcFieldType.Name)
// 		if !dstField.IsValid() || !dstField.CanSet() {
// 			continue
// 		}

// 		if dstField.Type() == srcField.Type() {
// 			dstField.Set(srcField)
// 		}
// 	}
// 	return nil
// }

// func mapSimpleStruct(src reflect.Value, dstType reflect.Type) (reflect.Value, error) {
// 	dst := reflect.New(dstType).Elem()

// 	srcType := src.Type()
// 	for i := 0; i < src.NumField(); i++ {
// 		field := srcType.Field(i)
// 		srcField := src.Field(i)

// 		dstField := dst.FieldByName(field.Name)
// 		if dstField.IsValid() && dstField.CanSet() && dstField.Type() == srcField.Type() {
// 			dstField.Set(srcField)
// 		}
// 	}

// 	return dst, nil
// }
