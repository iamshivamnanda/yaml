package yaml

import (
	"fmt"
	"reflect"
	"strings"
)

// Validator is a function that validates a value
type Validator func(value reflect.Value, yamlTag string, node *Node, param string) error

var validators = make(map[string]Validator)

// validateStruct performs validation based on struct tags
func ValidateStruct(out reflect.Value, node *Node) []error {

	var errors []error
	t := out.Type()
	if t.Kind() != reflect.Struct {
		return errors
	}

	fieldMap := make(map[string]*Node)
	if node.Kind == MappingNode {
		for i := 0; i < len(node.Content)-1; i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]
			if keyNode.Kind == ScalarNode {
				fieldMap[keyNode.Value] = valueNode
			}
		}
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("validate")
		yamlTag := field.Tag.Get("yaml")

		if tag != "" {
			if yamlTag == "" {
				yamlTag = field.Name
			}
			if childNode, exists := fieldMap[yamlTag]; exists {
				fieldValue := out.Field(i)
				if isYAMLValueType(field.Type) {
					fieldValue = fieldValue.FieldByName("Value")
				}
				errors = append(errors, applyValidationRules(fieldValue, tag, yamlTag, childNode)...)
			} else {
				if tagContainsRequired(tag) {
					errors = append(errors, CustomError{
						Line:   node.Line,
						Column: node.Column,
						Msg:    fmt.Sprintf("required field '%s' is missing", yamlTag),
					})
				}
			}
		}
	}

	return errors
}

// tagContainsRequired checks if the validate tag contains "required"
func tagContainsRequired(tag string) bool {
	rules := parseTag(tag)
	for _, rule := range rules {
		if rule.name == "required" {
			return true
		}
	}
	return false
}

// applyValidationRules applies validation rules from struct tags
func applyValidationRules(value reflect.Value, tag, yamlTag string, node *Node) []error {
	var errors []error
	rules := parseTag(tag)

	for _, rule := range rules {
		validator, exists := validators[rule.name]
		var param string
		if len(rule.params) > 0 {
			param = rule.params[0].(string)
		}
		if exists {
			err := validator(value, yamlTag, node, param)
			if err != nil {
				errors = append(errors, err)
			}
		} else {
			fmt.Println("Validator not found for: ", rule.name)
		}
	}

	return errors
}

// parseTag parses the validation tag into rules
func parseTag(tag string) []validationRule {
	var rules []validationRule
	tagParts := strings.Split(tag, ",")
	for _, part := range tagParts {
		ruleParts := strings.SplitN(part, "=", 2)
		ruleName := ruleParts[0]
		var params []interface{}
		if len(ruleParts) == 2 {
			params = append(params, ruleParts[1])
		}
		rules = append(rules, validationRule{name: ruleName, params: params})
	}
	return rules
}

// validationRule represents a single validation rule
type validationRule struct {
	name   string
	params []interface{}
}

// isEmpty checks if a value is empty
func isEmpty(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String, reflect.Array:
		return value.Len() == 0
	case reflect.Map, reflect.Slice:
		return value.IsNil() || value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return false
}

// RegisterValidator registers a custom validator
func RegisterValidator(name string, validator Validator) {
	validators[name] = validator
}

// init registers default validators
func init() {
	RegisterValidator("required", validateRequired)
	RegisterValidator("gt", validateGreaterThan)
	RegisterValidator("datetime", validateDatetime)
}
