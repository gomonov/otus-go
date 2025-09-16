package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	intKind    = "int"
	stringKind = "string"
)

var intConstraints = map[string]bool{
	"min": true,
	"max": true,
	"in":  true,
}

var strConstraints = map[string]bool{
	"len":    true,
	"regexp": true,
	"in":     true,
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

type Constraint struct {
	Name  string
	Value string
}

type Constraints []Constraint

func (v ValidationErrors) Error() string {
	length := 0
	for _, e := range v {
		length += len(e.Field) + len(e.Err.Error())
	}

	length += len(v) * (len(": ") + len("\n"))

	var builder strings.Builder
	builder.Grow(length)

	for _, err := range v {
		builder.WriteString(err.Field)
		builder.WriteString(": ")
		builder.WriteString(err.Err.Error())
		builder.WriteString("\n")
	}

	return builder.String()
}

func Validate(v interface{}) error {
	var validateErrors ValidationErrors

	valStruct := reflect.ValueOf(v)
	typeStruct := valStruct.Type()

	if valStruct.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, got a %T", v)
	}

	for i := 0; i < valStruct.NumField(); i++ {
		typeField := typeStruct.Field(i)
		valField := valStruct.Field(i)

		errors, err := validateField(valField, typeField)
		if err != nil {
			return err
		}
		if errors != nil {
			validateErrors = append(validateErrors, errors...)
		}
	}

	return validateErrors
}

//nolint:gocognit
func validateField(fieldVal reflect.Value, fieldType reflect.StructField) (ValidationErrors, error) {
	fieldName := fieldType.Name

	tag := fieldType.Tag.Get("validate")
	if tag == "" {
		return nil, nil
	}

	switch fieldVal.Kind().String() {
	case intKind:
		constraints, err := parseConstraints(tag, reflect.Int)
		if err != nil {
			return nil, err
		}

		validateErrors, err := validateInt(fieldVal.Int(), fieldName, constraints)
		if err != nil {
			return nil, err
		}

		return validateErrors, nil
	case stringKind:
		constraints, err := parseConstraints(tag, reflect.String)
		if err != nil {
			return nil, err
		}

		validateErrors, err := validateStr(fieldVal.String(), fieldName, constraints)
		if err != nil {
			return nil, err
		}

		return validateErrors, nil
	case "slice":
		var validateErrorsTotal ValidationErrors
		elemKind := fieldVal.Type().Elem().Kind()

		switch elemKind.String() {
		case intKind:
			constraints, errParse := parseConstraints(tag, reflect.Int)
			if errParse != nil {
				return nil, errParse
			}

			for _, v := range fieldVal.Interface().([]int) {
				validateErrors, err := validateInt(int64(v), fieldName, constraints)
				if err != nil {
					return nil, err
				}

				validateErrorsTotal = append(validateErrorsTotal, validateErrors...)
			}
			return validateErrorsTotal, nil
		case stringKind:
			constraints, errParse := parseConstraints(tag, reflect.String)
			if errParse != nil {
				return nil, errParse
			}

			for _, v := range fieldVal.Interface().([]string) {
				validateErrors, err := validateStr(v, fieldName, constraints)
				if err != nil {
					return nil, err
				}

				validateErrorsTotal = append(validateErrorsTotal, validateErrors...)
			}
			return validateErrorsTotal, nil
		default:
			return nil, fmt.Errorf("unsupported type: %s", elemKind)
		}
	default:
		return nil, fmt.Errorf("unsupported type: %s", fieldVal.Kind())
	}
}

func validateInt(value int64, fieldName string, constraints Constraints) (ValidationErrors, error) {
	var validateErrors ValidationErrors
	for _, constraint := range constraints {
		switch constraint.Name {
		case "min":
			limit, err := strconv.Atoi(constraint.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid value %q for constraint %q: %w", constraint.Value, fieldName, err)
			}
			if value < int64(limit) {
				validateError := ValidationError{Field: fieldName, Err: fmt.Errorf("value %d is less than %d", value, limit)}
				validateErrors = append(validateErrors, validateError)
			}
		case "max":
			limit, err := strconv.Atoi(constraint.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid value %q for constraint %q: %w", constraint.Value, fieldName, err)
			}
			if value > int64(limit) {
				validateError := ValidationError{Field: fieldName, Err: fmt.Errorf("value %d is greater than max %d", value, limit)}
				validateErrors = append(validateErrors, validateError)
			}
		case "in":
			found := false
			for _, numStr := range strings.Split(constraint.Value, ",") {
				num, err := strconv.Atoi(strings.TrimSpace(numStr))
				if err != nil {
					return nil, fmt.Errorf("invalid value %q for constraint %q: %w", constraint.Value, fieldName, err)
				}
				if value == int64(num) {
					found = true
					break
				}
			}
			if !found {
				validateError := ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("value %d not in list %s", value, constraint.Value),
				}
				validateErrors = append(validateErrors, validateError)
			}
		default:
			return nil, fmt.Errorf("unknown int constraint: %s", constraint.Name)
		}
	}

	return validateErrors, nil
}

func validateStr(value string, fieldName string, constraints Constraints) (ValidationErrors, error) {
	var validateErrors ValidationErrors
	for _, constraint := range constraints {
		switch constraint.Name {
		case "len":
			length, err := strconv.Atoi(constraint.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid value %q for constraint %q: %w", constraint.Value, fieldName, err)
			}
			if length != len(value) {
				validateError := ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("length %d does not equal %d", len(value), length),
				}
				validateErrors = append(validateErrors, validateError)
			}
		case "regexp":
			re, err := regexp.Compile(constraint.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid regexp %q: %w", constraint.Value, err)
			}
			if !re.MatchString(value) {
				validateError := ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("value %q does not match regexp %q", value, constraint.Value),
				}
				validateErrors = append(validateErrors, validateError)
			}
		case "in":
			found := false
			for _, str := range strings.Split(constraint.Value, ",") {
				if strings.TrimSpace(str) == value {
					found = true
					break
				}
			}
			if !found {
				validateError := ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("value %s not in list %s", value, constraint.Value),
				}
				validateErrors = append(validateErrors, validateError)
			}
		default:
			return nil, fmt.Errorf("unknown string constraint: %s", constraint.Name)
		}
	}

	return validateErrors, nil
}

func parseConstraints(tag string, kind reflect.Kind) (Constraints, error) {
	constraints := strings.Split(tag, "|")
	result := make(Constraints, 0, len(constraints))

	for _, constraint := range constraints {
		constraintParts := strings.SplitN(constraint, ":", 2)
		if len(constraintParts) != 2 {
			return nil, fmt.Errorf("invalid tag format: %s", constraint)
		}

		name := constraintParts[0]
		value := constraintParts[1]

		var validConstraints map[string]bool

		switch kind.String() {
		case intKind:
			validConstraints = intConstraints
		case stringKind:
			validConstraints = strConstraints
		default:
			return nil, fmt.Errorf("unsupported kind: %v", kind)
		}

		if !validConstraints[name] {
			return nil, fmt.Errorf("invalid constraint: %s for type %s", name, kind.String())
		}

		result = append(result, Constraint{Name: name, Value: value})
	}

	return result, nil
}
