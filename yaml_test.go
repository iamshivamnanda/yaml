package yaml_test

import (
	"fmt"
	"testing"

	yaml "github.com/iamshivamnanda/yaml/v3"
)

func TestUnmarshal(t *testing.T) {
	data := []byte(`
name: John Doe
age: 30
email: john.doe@example.com
created_at: 2023-06-21T15:00:00Z
`)

	type Person struct {
		Name      string `yaml:"name" validate:"required"`
		Age       int    `yaml:"age" validate:"gt=18"`
		CreatedAt string `yaml:"created_at" validate:"datetime=2006-01-02T15:04:05Z07:00"`
		Email     string `yaml:"email" validate:"email"`
	}
	var person Person

	err := yaml.Unmarshal(data, &person)
	if err != nil {
		t.Errorf("Unmarshal failed: %v", err)
	}

	expectedName := "John Doe"
	if person.Name != expectedName {
		t.Errorf("Expected name to be %q, but got %q", expectedName, person.Name)
	}

	expectedAge := 30
	if person.Age != expectedAge {
		t.Errorf("Expected age to be %d, but got %d", expectedAge, person.Age)
	}

	expectedEmail := "john.doe@example.com"
	if person.Email != expectedEmail {
		t.Errorf("Expected email to be %q, but got %q", expectedEmail, person.Email)
	}

	data = []byte(`
age: 16
email: john.doe@example.com
created_at: 2023-06-21T15:00:00Z
`)
	var person2 Person
	if err := yaml.Unmarshal(data, &person2); err != nil {
		fmt.Println("Error: ", err)
		t.Log("Error: ", err)
	} else {
		fmt.Printf("Parsed Config: %+v\n", person2)
	}
}
