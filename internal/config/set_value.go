package config

import (
	"fmt"
	"gontlm"
	"reflect"
	"strings"
)

func setValue(cfg *gontlm.Config, line string) error {
	name, value, err := splitNameValue(line)
	if err != nil {
		return err
	}

	return setFieldByName(cfg, name, value)
}

func splitNameValue(line string) (string, string, error) {
	items := strings.Fields(line)
	count := len(items)

	if count == 2 {
		return items[0], items[1], nil
	}

	return "", "", fmt.Errorf("malformed line %s", line)
}

func setFieldByName(cfg *gontlm.Config, name, value string) error {
	s := reflect.ValueOf(cfg).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if typeOfT.Field(i).Name == name {
			if setFieldDependsOnType(f, value) {
				return nil
			}

			break
		}
	}

	return fmt.Errorf("unknown keyword '%s'", name)
}

func setFieldDependsOnType(f reflect.Value, value string) bool {
	switch f.Type().Kind() {
	case reflect.String:
		setStringField(f, value)
	case reflect.Slice:
		setSliceField(f, value)
	default:
		return false
	}

	return true
}

func setStringField(f reflect.Value, value string) {
	f.SetString(value)
}

func setSliceField(f reflect.Value, value string) {
	slice := reflect.ValueOf(f.Interface())
	for _, v := range strings.Split(value, ",") {
		v = strings.TrimSpace(v)
		slice = reflect.Append(slice, reflect.ValueOf(v))
	}
	f.Set(slice)
}
