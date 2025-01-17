package config

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateClientConfig(conf *ClientConfig) error {
	// implement this
	return validation.ValidateStruct(conf,
		validation.Field(&conf.Version, validation.Required),
		validation.Field(&conf.Host, validation.Required),
		nestedFields(&conf.Log,
			validation.Field(&conf.Log.Level, validation.In(
				LogLevelDebug,
				LogLevelInfo,
				LogLevelWarning,
				LogLevelError,
				LogLevelFatal,
			)),
		),
		validation.Field(&conf.Namespaces, validation.By(validateNamespaces)),
		// ... etc
	)
}

func ValidateServerConfig(conf *ServerConfig) error {
	// implement this
	return nil
}

func validateNamespaces(value interface{}) error {
	namespaces, ok := value.([]*Namespace)
	if !ok {
		return errors.New("can't convert value to namespaces")
	}

	m := map[string]int{}
	for _, n := range namespaces {
		if n == nil {
			continue
		}
		m[n.Name]++
	}

	dup := []string{}
	for k, v := range m {
		if v > 1 {
			dup = append(dup, k)
		}
	}

	if len(dup) > 0 {
		return fmt.Errorf("duplicate namespaces are not allowed [%s]", strings.Join(dup, ","))
	}

	return nil
}

// ozzo-validation helper for nested validation struct
// https://github.com/go-ozzo/ozzo-validation/issues/136
func nestedFields(target interface{}, fieldRules ...*validation.FieldRules) *validation.FieldRules {
	return validation.Field(target, validation.By(func(value interface{}) error {
		valueV := reflect.Indirect(reflect.ValueOf(value))
		if valueV.CanAddr() {
			addr := valueV.Addr().Interface()
			return validation.ValidateStruct(addr, fieldRules...)
		}
		return validation.ValidateStruct(target, fieldRules...)
	}))
}
