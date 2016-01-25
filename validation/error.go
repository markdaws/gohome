package validation

import (
	"fmt"
	"reflect"
	"strings"
)

type Error struct {
	MSG   string
	Field string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.MSG)
}

type Errors struct {
	Errors []Error
}

func (errs *Errors) Add(msg, field string) {
	errs.Errors = append(errs.Errors, Error{
		MSG:   msg,
		Field: field,
	})
}
func (errs *Errors) Has() bool {
	return len(errs.Errors) > 0
}
func (errs *Errors) Error() string {
	return fmt.Sprintf("validation failed, %d errors", len(errs.Errors))
}

type ErrorJSON struct {
	Errors map[string]map[string]string `json:"errors"`
}

func NewErrorJSON(item interface{}, clientID string, errors *Errors) ErrorJSON {
	valErr := ErrorJSON{
		Errors: make(map[string]map[string]string),
	}
	for _, e := range errors.Errors {
		jsonField, err := JSONTagForField(item, e.Field)
		if err != nil {
			fmt.Println(err)
			//TODO: Log?
			continue
		}
		valErr.Errors[clientID+"."+jsonField] = make(map[string]string)
		valErr.Errors[clientID+"."+jsonField]["message"] = e.MSG
	}
	return valErr
}

func JSONTagForField(s interface{}, name string) (string, error) {
	f, ok := reflect.TypeOf(s).Elem().FieldByName(name)
	if !ok {
		return "", fmt.Errorf("invalid field: %s", name)
	}

	parts := strings.Split(string(f.Tag), ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("unexpected tag format")
	}
	return strings.Replace(parts[1], "\"", "", -1), nil
}
