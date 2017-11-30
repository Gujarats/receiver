package receiver

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/magiconair/properties/assert"
)

// struct for test : tag written correctly
type testStructOk struct {
	Lat      float64 `request:"latitude,required"`
	Lon      float64 `request:"longitude,required"`
	Name     string  `request:"name,required"`
	Distance int64   `request:"distance,optional"`
}

// struct for test : tag written incorrecly
type testStructNotOk struct {
	Lat      float64 `request:"latitude,required"`
	Lon      float64 `request:"longitude,required"`
	Name     string  `request:"name,required"`
	Distance int64   `request:"name,optional"`
}

func TestReceiver(t *testing.T) {
	testObjects := []struct {
		r http.Request

		// This is for testing struct assuming that the struct and the tag written correcly
		data *testStructOk

		// This is for testing struct assuming that the tag written incorrecly
		dataNotOk   *testStructNotOk
		ok          bool
		expectedErr error
	}{
		// test 0
		{
			r: http.Request{
				Form: url.Values{
					"latitude":  []string{"123.123456"},
					"longitude": []string{"12.123456"},
					"name":      []string{"gujarat"},
					"distance":  []string{"1234"},
				},
			},
			data: &testStructOk{},
			ok:   true,
		},

		// test 1
		{
			r: http.Request{
				Form: url.Values{
					"latitude":  []string{"123.123456"},
					"longitude": []string{"12.123456"},
					"name":      []string{"gujarat"},
					"distance":  []string{},
				},
			},
			data: &testStructOk{},
			ok:   true,
		},

		// test 2
		{
			r: http.Request{
				Form: url.Values{
					"latitude":  []string{"123.123456"},
					"longitude": []string{"12.123456"},
					"name":      []string{},
					"distance":  []string{},
				},
			},
			data:        &testStructOk{},
			ok:          false,
			expectedErr: errors.New("required value form is empty on the key = name"),
		},

		// test 3
		{
			r: http.Request{
				Form: url.Values{
					"latitude":  []string{},
					"longitude": []string{"12.123456"},
					"name":      []string{},
					"distance":  []string{},
				},
			},
			data:        &testStructOk{},
			ok:          false,
			expectedErr: errors.New("required value form is empty on the key = latitude"),
		},

		// test 4
		{
			r: http.Request{
				Form: url.Values{
					"latitude":  []string{"123.123456"},
					"longitude": []string{"12.123456"},
					"name":      []string{"gujarat"},
					"distance":  []string{"1234"},
				},
			},
			dataNotOk:   &testStructNotOk{},
			ok:          false,
			expectedErr: errors.New("Parse integer failed check your tag key name"),
		},
	}

	for index, testObject := range testObjects {
		if testObject.data != nil {
			err := SetData(testObject.data, &testObject.r)
			if err != nil && testObject.ok {
				t.Errorf("%+v\n", err)
			}

			if isAnyFieldEmpty(testObject.data) && testObject.ok {
				t.Errorf("Failed index = %v, struct = %+v is empty", index, testObject.data)
			}

			if !testObject.ok {
				assert.Equal(t, testObject.expectedErr, err)
			}
		}

		if testObject.dataNotOk != nil {
			err := SetData(testObject.dataNotOk, &testObject.r)
			if err != nil && testObject.ok {
				t.Errorf("%+v\n", err)
			}

			if isAnyFieldEmpty(testObject.dataNotOk) && testObject.ok {
				t.Errorf("Failed index = %v, struct = %+v is empty", index, testObject.data)
			}

			if !testObject.ok {
				assert.Equal(t, testObject.expectedErr, err)
			}

		}
	}
}

// passing struct and check all the field,
// false if there is empty value.
func isAnyFieldEmpty(input interface{}) bool {
	object := reflect.ValueOf(input)

	dataType := reflect.TypeOf(input)
	for index := 0; index < object.Elem().NumField(); index++ {
		field := dataType.Elem().Field(index)
		tagValue := field.Tag.Get(tagKey)
		_, validator := getRequirements(tagValue)
		if isZeroOfUnderlyingType(object.Elem().Field(index).Interface()) && validator == "required" {
			return true
		}
	}

	return false
}

// compare the object value to non assign value.
func isZeroOfUnderlyingType(objectValue interface{}) bool {
	return reflect.DeepEqual(objectValue, reflect.Zero(reflect.TypeOf(objectValue)).Interface())
}
