/*
Package validation provides some helper type for validating fields
*/
package validation

import (
	"fmt"
	"reflect"
	"strings"
)

// Error represents a validation error
type Error struct {
	// An informative string that will be displayed to the user, detailing
	// the error e.g. "required field"
	MSG string

	// The name of the field in the type which has the validation error e.g. "Name"
	Field string

	IgnoreJSONTag bool
}

// Error returns a friendly error string
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.MSG)
}

// Errors is a helper type for handling multiple validation errors
type Errors struct {
	Errors []Error
}

// NewErrors creates an Errors instance and adds one error to the Errors
// field, using the field and msg values
func NewErrors(field, msg string, isExplicit bool) *Errors {
	e := &Errors{}
	if isExplicit {
		e.AddExplicitField(msg, field)
	} else {
		e.Add(msg, field)
	}
	return e
}

// Add inserts a new error into the Errors field
func (errs *Errors) Add(msg, field string) {
	errs.Errors = append(errs.Errors, Error{
		MSG:           msg,
		Field:         field,
		IgnoreJSONTag: false,
	})
}

func (errs *Errors) AddExplicitField(msg, field string) {
	errs.Errors = append(errs.Errors, Error{
		MSG:           msg,
		Field:         field,
		IgnoreJSONTag: true,
	})
}

// Has returns true if the Errors instance contains any errors
func (errs *Errors) Has() bool {
	return len(errs.Errors) > 0
}

// Error returns a friendly error string
func (errs *Errors) Error() string {
	return fmt.Sprintf("validation failed, %d errors", len(errs.Errors))
}

// ErrorJSON provides a consistent format to send back validation errors to
// a client, it is an object, that contains one field "errors", the errors field
// is an object, where the keys are a unique identifier to the object that had
// the error, keyed by ids ID, e.g. a zone ID or scene ID, whaterver has been validated,
// then that has a map which is keyed by the field name which has validation issues
// {
//    errors: {
//        "objectID": {
//            "field1": { message: "required field" },
//	    	  "field2": { message: "must be greater than 50" }
//        }
//    }
// }
type ErrorJSON struct {
	Errors map[string]map[string]map[string]string `json:"errors"`
}

// NewErrorJSON returns an ErrorJSON instance populated with all of the errors in
// the errors parameter
func NewErrorJSON(item interface{}, objectID string, errors *Errors) ErrorJSON {
	valErr := ErrorJSON{
		Errors: make(map[string]map[string]map[string]string),
	}
	valErr.Errors[objectID] = make(map[string]map[string]string)

	for _, e := range errors.Errors {
		var jsonField string
		if e.IgnoreJSONTag {
			jsonField = e.Field
		} else {
			var err error
			jsonField, err = JSONTagForField(item, e.Field)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
		valErr.Errors[objectID][jsonField] = make(map[string]string)
		valErr.Errors[objectID][jsonField]["message"] = e.MSG
	}
	return valErr
}

// JSONTagForField returns the value associated with the json tag for the
// field specified by the name parameter e.g.
// type Foo struct { Name string `json:"name"` } would return 'name'
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
