package yaml

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// validateRequired checks if a value is required
func validateRequired(value reflect.Value, yamlTag string, node *Node, param string) error {
	if isEmpty(value) {
		return CustomError{
			Line:   node.Line,
			Column: node.Column,
			Msg:    fmt.Sprintf("field '%s' is required", yamlTag),
		}
	}
	return nil
}

// validateGreaterThan checks if a value is greater than a threshold
func validateGreaterThan(value reflect.Value, yamlTag string, node *Node, param string) error {
	threshold, err := strconv.Atoi(param)
	if err != nil {
		return err
	}
	if value.Int() <= int64(threshold) {
		return CustomError{
			Line:   node.Line,
			Column: node.Column,
			Msg:    fmt.Sprintf("field '%s' must be greater than %d", yamlTag, threshold),
		}
	}
	return nil
}

// validateDatetime checks if a value matches a datetime format
func validateDatetime(value reflect.Value, yamlTag string, node *Node, param string) error {
	if _, err := time.Parse(param, value.String()); err != nil {
		return CustomError{
			Line:   node.Line,
			Column: node.Column,
			Msg:    fmt.Sprintf("field '%s' has invalid datetime format", yamlTag),
		}
	}
	return nil
}
