package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "John Doe",
				Age:    18,
				Email:  "johndoe@example.com",
				Role:   UserRole("admin"),
				Phones: []string{"12345678901", "45678901234"},
			},
			expectedErr: ValidationErrors(nil),
		},
		{
			in: User{
				ID:     "123456789012345678901234567890",
				Name:   "John Doe",
				Age:    60,
				Email:  "john.doe-example.com1",
				Role:   UserRole("user"),
				Phones: []string{"123", "456"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: fmt.Errorf("length 30 does not equal 36")},
				{Field: "Age", Err: fmt.Errorf("value 60 is greater than max 50")},
				{Field: "Email", Err: fmt.Errorf(`value "john.doe-example.com1" does not match regexp "^\\w+@\\w+\\.\\w+$"`)},
				{Field: "Role", Err: fmt.Errorf(`value user not in list admin,stuff`)},
				{Field: "Phones", Err: fmt.Errorf("length 3 does not equal 11")},
				{Field: "Phones", Err: fmt.Errorf("length 3 does not equal 11")},
			},
		},
		{
			in: App{
				Version: "1.0.0",
			},
			expectedErr: ValidationErrors(nil),
		},
		{
			in: App{
				Version: "v1.0.0",
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: fmt.Errorf("length 6 does not equal 5")},
			},
		},
		{
			in: Token{
				Header:    make([]byte, 0),
				Payload:   make([]byte, 0),
				Signature: make([]byte, 0),
			},
			expectedErr: ValidationErrors(nil),
		},
		{
			in:          Response{Code: 200},
			expectedErr: ValidationErrors(nil),
		},
		{
			in: Response{Code: 302},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: fmt.Errorf("value 302 not in list 200,404,500")},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			fmt.Println(err.Error())
			require.Equal(t, tt.expectedErr.Error(), err.Error())
			_ = tt
		})
	}
}
