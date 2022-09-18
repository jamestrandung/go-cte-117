package cte

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummy struct {
	pointer    *string
	nonPointer string
}

func (d dummy) DoDummy(arg1 string, arg2 int) (float64, []dummy) {
	return 0, nil
}

func TestSwallowErrPlanExecutionEndingEarly(t *testing.T) {
	scenarios := []struct {
		desc     string
		err      error
		expected error
	}{
		{
			desc:     "ErrPlanExecutionEndingEarly",
			err:      ErrPlanExecutionEndingEarly,
			expected: nil,
		},
		{
			desc:     "ErrRootPlanExecutionEndingEarly",
			err:      ErrRootPlanExecutionEndingEarly,
			expected: nil,
		},
		{
			desc:     "other errors",
			err:      assert.AnError,
			expected: assert.AnError,
		},
	}

	for _, scenario := range scenarios {
		s := scenario

		t.Run(
			s.desc, func(t *testing.T) {
				assert.Equal(t, s.expected, swallowErrPlanExecutionEndingEarly(s.err))
			},
		)
	}
}

func TestExtractFullNameFromValue(test *testing.T) {
	defer func(original func(t reflect.Type) string) {
		extractFullNameFromType = original
	}(extractFullNameFromType)

	extractFullNameFromType = func(t reflect.Type) string {
		assert.Equal(test, reflect.TypeOf(dummy{}), t)
		return "dummy"
	}

	assert.Equal(test, "dummy", extractFullNameFromValue(dummy{}))
}

func TestExtractFullNameFromType(t *testing.T) {
	scenarios := []struct {
		desc     string
		t        reflect.Type
		expected string
	}{
		{
			desc:     "pointer",
			t:        reflect.TypeOf(&dummy{}),
			expected: "github.com/jamestrandung/go-cte-117/cte/dummy",
		},
		{
			desc:     "non-pointer",
			t:        reflect.TypeOf(dummy{}),
			expected: "github.com/jamestrandung/go-cte-117/cte/dummy",
		},
	}

	for _, scenario := range scenarios {
		s := scenario

		t.Run(
			s.desc, func(t *testing.T) {
				assert.Equal(t, s.expected, extractFullNameFromType(s.t))
			},
		)
	}
}

func TestExtractShortName(t *testing.T) {
	scenarios := []struct {
		desc     string
		fullName string
		expected string
	}{
		{
			desc:     "containing /",
			fullName: "github.com/jamestrandung/go-cte-117/cte/dummy",
			expected: "dummy",
		},
		{
			desc:     "not containing /",
			fullName: "dummy",
			expected: "dummy",
		},
	}

	for _, scenario := range scenarios {
		s := scenario

		t.Run(
			s.desc, func(t *testing.T) {
				assert.Equal(t, s.expected, extractShortName(s.fullName))
			},
		)
	}
}

func TestExtractFieldTypes(t *testing.T) {
	d := reflect.TypeOf(dummy{})

	pointerField := d.Field(0)
	t.Run(
		"pointer field", func(t *testing.T) {
			isPointerType, valueType, pointerType := extractFieldTypes(pointerField)

			assert.Equal(t, true, isPointerType)
			assert.Equal(t, reflect.TypeOf(""), valueType)
			assert.Equal(t, reflect.PointerTo(reflect.TypeOf("")), pointerType)
		},
	)

	nonPointerField := d.Field(1)
	t.Run(
		"non-pointer field", func(t *testing.T) {
			isPointerType, valueType, pointerType := extractFieldTypes(nonPointerField)

			assert.Equal(t, false, isPointerType)
			assert.Equal(t, reflect.TypeOf(""), valueType)
			assert.Equal(t, reflect.PointerTo(reflect.TypeOf("")), pointerType)
		},
	)
}

func TestExtractUnderlyingType(test *testing.T) {
	defer func(original func(t reflect.Type) reflect.Type) {
		extractNonPointerType = original
	}(extractNonPointerType)

	scenarios := []struct {
		desc string
		test func(t *testing.T)
	}{
		{
			desc: "pointer",
			test: func(test *testing.T) {
				s := "string"
				v := reflect.ValueOf(&s)

				expected := reflect.TypeOf(1)
				extractNonPointerType = func(t reflect.Type) reflect.Type {
					assert.Equal(test, v.Type(), t)

					return expected
				}

				actual := extractUnderlyingType(v)
				assert.Equal(test, expected, actual)
			},
		},
		{
			desc: "non-pointer",
			test: func(test *testing.T) {
				s := "string"
				v := reflect.ValueOf(s)

				expected := reflect.TypeOf(1)
				extractNonPointerType = func(t reflect.Type) reflect.Type {
					assert.Equal(test, v.Type(), t)

					return expected
				}

				actual := extractUnderlyingType(v)
				assert.Equal(test, expected, actual)
			},
		},
	}

	for _, scenario := range scenarios {
		s := scenario

		test.Run(s.desc, s.test)
	}
}

func TestExtractNonPointerType(t *testing.T) {
	scenarios := []struct {
		desc     string
		t        reflect.Type
		expected reflect.Type
	}{
		{
			desc:     "pointer",
			t:        reflect.TypeOf(&dummy{}),
			expected: reflect.TypeOf(dummy{}),
		},
		{
			desc:     "non-pointer",
			t:        reflect.TypeOf(dummy{}),
			expected: reflect.TypeOf(dummy{}),
		},
	}

	for _, scenario := range scenarios {
		s := scenario

		t.Run(
			s.desc, func(t *testing.T) {
				assert.Equal(t, s.expected, extractNonPointerType(s.t))
			},
		)
	}
}
