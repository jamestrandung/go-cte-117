package cte

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestMethod_HasSameSignature(t *testing.T) {
	scenarios := []struct {
		desc     string
		m1       method
		m2       method
		expected bool
	}{
		{
			desc: "different name",
			m1: method{
				owningType: "owningType",
				name:       "name1",
				arguments:  "arguments",
				outputs:    "outputs",
			},
			m2: method{
				owningType: "owningType",
				name:       "name2",
				arguments:  "arguments",
				outputs:    "outputs",
			},
			expected: false,
		},
		{
			desc: "different arguments",
			m1: method{
				owningType: "owningType",
				name:       "name",
				arguments:  "arguments1",
				outputs:    "outputs",
			},
			m2: method{
				owningType: "owningType",
				name:       "name",
				arguments:  "arguments2",
				outputs:    "outputs",
			},
			expected: false,
		},
		{
			desc: "different outputs",
			m1: method{
				owningType: "owningType",
				name:       "name",
				arguments:  "arguments",
				outputs:    "outputs1",
			},
			m2: method{
				owningType: "owningType",
				name:       "name",
				arguments:  "arguments",
				outputs:    "outputs2",
			},
			expected: false,
		},
		{
			desc: "same signature",
			m1: method{
				owningType: "owningType1",
				name:       "name",
				arguments:  "arguments",
				outputs:    "outputs",
			},
			m2: method{
				owningType: "owningType2",
				name:       "name",
				arguments:  "arguments",
				outputs:    "outputs",
			},
			expected: true,
		},
	}

	for _, scenario := range scenarios {
		s := scenario

		t.Run(s.desc, func(t *testing.T) {
			assert.Equal(t, s.expected, s.m1.hasSameSignature(s.m2))
		})
	}
}

func TestMethod_String(t *testing.T) {
	scenarios := []struct {
		desc     string
		m        method
		expected string
	}{
		{
			desc: "has multiple outputs",
			m: method{
				name:      "name",
				arguments: "arguments",
				outputs:   "output1,output2",
			},
			expected: "name(arguments) (output1,output2)",
		},
		{
			desc: "has no outputs",
			m: method{
				name:      "name",
				arguments: "arguments",
				outputs:   "",
			},
			expected: "name(arguments)",
		},
		{
			desc: "has one output",
			m: method{
				name:      "name",
				arguments: "arguments",
				outputs:   "output",
			},
			expected: "name(arguments) output",
		},
		{
			desc: "has owning type",
			m: method{
				owningType: "owningType",
				name:       "name",
				arguments:  "arguments",
				outputs:    "output1,output2",
			},
			expected: "owningType.name(arguments) (output1,output2)",
		},
	}

	for _, scenario := range scenarios {
		s := scenario

		t.Run(s.desc, func(t *testing.T) {
			assert.Equal(t, s.expected, s.m.String())
		})
	}
}

func TestExtractMethodDetails(t *testing.T) {
	m, ok := reflect.TypeOf(dummy{}).MethodByName("DoDummy")
	if ok = assert.Equal(t, true, ok, "dummy should have a method called DoDummy"); !ok {
		return
	}

	t.Run("ignoreFirstReceiverArgument == true", func(t *testing.T) {
		expected := method{
			name:      "DoDummy",
			arguments: "string,int",
			outputs:   "float64,[]cte.dummy",
		}

		assert.Equal(t, expected, extractMethodDetails(m, true))
	})

	t.Run("ignoreFirstReceiverArgument == false", func(t *testing.T) {
		expected := method{
			name:      "DoDummy",
			arguments: "cte.dummy,string,int",
			outputs:   "float64,[]cte.dummy",
		}

		assert.Equal(t, expected, extractMethodDetails(m, false))
	})
}

func TestStructDisassembler_IsAvailableMoreThanOnce(t *testing.T) {
	// Dummy data
	rootPlanName := "rootPlanName"
	cs := componentStack{}
	m := method{
		owningType: "owningType",
		name:       "name",
		arguments:  "arguments",
		outputs:    "output1,output2",
	}

	sd := newStructDisassembler()
	assert.Equal(t, false, sd.isAvailableMoreThanOnce(m))

	sd.addAvailableMethod(rootPlanName, cs, m)
	assert.Equal(t, false, sd.isAvailableMoreThanOnce(m))

	sd.addAvailableMethod(rootPlanName, cs, m)
	assert.Equal(t, true, sd.isAvailableMoreThanOnce(m))
}

func TestStructDisassembler_FindAvailableMethods(t *testing.T) {
	sd := newStructDisassembler()

	func() {
		ms, ok := sd.findAvailableMethods("dummy")
		assert.Equal(t, 0, ms.count())
		assert.False(t, ok)
	}()

	sd.availableMethods["dummy"] = methodSet{method{name: "dummy"}: struct{}{}}

	func() {
		ms, ok := sd.findAvailableMethods("dummy")
		assert.Equal(t, 1, ms.count())
		assert.True(t, ok)
	}()

	sd.availableMethods["empty"] = methodSet{}

	func() {
		ms, ok := sd.findAvailableMethods("empty")
		assert.Equal(t, 0, ms.count())
		assert.False(t, ok)
	}()
}

func TestStructDisassembler_FindMethodLocations(t *testing.T) {
	methodToLookFor1 := method{
		owningType: "owningType1",
		name:       "name",
		arguments:  "arguments",
		outputs:    "output1,output2",
	}

	methodToLookFor2 := method{
		owningType: "owningType2",
		name:       "name",
		arguments:  "arguments",
		outputs:    "output1,output2",
	}

	methodFromAnotherRootPlan := method{
		owningType: "owningType3",
		name:       "name",
		arguments:  "arguments",
		outputs:    "output1,output2",
	}

	sd := newStructDisassembler()
	sd.addAvailableMethod("rootPlanName1", componentStack{"dummyComponent1"}, methodToLookFor1)
	sd.addAvailableMethod("rootPlanName1", componentStack{"dummyComponent2"}, methodToLookFor2)
	sd.addAvailableMethod("rootPlanName2", componentStack{"dummyComponent1"}, methodFromAnotherRootPlan)

	ms := make(methodSet)
	ms.add(methodToLookFor1)
	ms.add(methodToLookFor2)

	result := sd.findMethodLocations(ms, "rootPlanName1")
	assert.Equal(t, []string{"dummyComponent1", "dummyComponent2"}, result)
}

func TestAddAvailableMethod(t *testing.T) {
	// Dummy data
	rootPlanName := "rootPlanName"

	cs := componentStack{}
	cs = cs.push("dummyComponent1")

	m := method{
		owningType: "owningType",
		name:       "name",
		arguments:  "arguments",
		outputs:    "output1,output2",
	}

	sd := newStructDisassembler()
	assert.Equal(t, 0, len(sd.availableMethods))
	assert.Equal(t, 0, sd.methodsAvailableMoreThanOnce.count())
	assert.Equal(t, 0, len(sd.methodLocations))

	sd.addAvailableMethod(rootPlanName, cs, m)
	assert.Equal(t, 1, len(sd.availableMethods))
	assert.Equal(t, 0, sd.methodsAvailableMoreThanOnce.count())
	assert.Equal(t, 1, len(sd.methodLocations))
	assert.Equal(t, 1, len(sd.methodLocations[m]))
	assert.Equal(t, sd.methodLocations[m][0], methodLocation{
		rootPlanName:   rootPlanName,
		componentStack: cs,
	})

	cs = cs.push("dummyComponent2")

	sd.addAvailableMethod(rootPlanName, cs, m)
	assert.Equal(t, 1, len(sd.availableMethods))
	assert.Equal(t, 1, sd.methodsAvailableMoreThanOnce.count())
	assert.Equal(t, 1, len(sd.methodLocations))
	assert.Equal(t, 2, len(sd.methodLocations[m]))
	assert.Equal(t, sd.methodLocations[m][1], methodLocation{
		rootPlanName:   rootPlanName,
		componentStack: cs,
	})

	// Changes to component stack must not affect the recorded locations
	assert.NotEqual(t, sd.methodLocations[m][0], methodLocation{
		rootPlanName:   rootPlanName,
		componentStack: cs,
	})

	cs = cs.push("dummyComponent3")

	assert.Equal(t, 1, len(sd.methodLocations[m][0].componentStack))
	assert.Equal(t, 2, len(sd.methodLocations[m][1].componentStack))
	assert.Equal(t, 3, len(cs))
}

func TestPerformMethodExtraction(test *testing.T) {
	t := reflect.TypeOf(dummy{})
	rootPlanName := "rootPlanName"
	cs := componentStack{}

	sd := newStructDisassembler()
	sdMock := &mockIStructDisassembler{}
	sd.itself = sdMock

	childMethods := []method{
		{
			owningType: "owningType1",
			name:       "name",
			arguments:  "arguments",
			outputs:    "output1,output2",
		},
	}

	sdMock.On("extractChildMethods", mock.Anything, mock.Anything, mock.Anything).
		Return(func(t reflect.Type, rootPlanNameIn string, csIn componentStack) []method {
			assert.Equal(test, reflect.TypeOf(dummy{}), t)
			assert.Equal(test, rootPlanName, rootPlanNameIn)
			if assert.Equal(test, 1, len(csIn)) {
				assert.Equal(test, "github.com/jamestrandung/go-cte-117/cte/dummy", csIn[0])
			}

			return childMethods
		}).Once()

	sdMock.On("extractOwnMethods", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(func(t reflect.Type, rootPlanNameIn string, csIn componentStack, hoistedMethods []method) []method {
			assert.Equal(test, reflect.TypeOf(dummy{}), t)
			assert.Equal(test, rootPlanName, rootPlanNameIn)
			if assert.Equal(test, 1, len(csIn)) {
				assert.Equal(test, "github.com/jamestrandung/go-cte-117/cte/dummy", csIn[0])
			}
			assert.Equal(test, childMethods, hoistedMethods)

			return []method{
				{
					owningType: "owningType2",
					name:       "name",
					arguments:  "arguments",
					outputs:    "output1,output2",
				},
			}
		}).Once()

	expected := []method{
		{
			owningType: "owningType1",
			name:       "name",
			arguments:  "arguments",
			outputs:    "output1,output2",
		},
		{
			owningType: "owningType2",
			name:       "name",
			arguments:  "arguments",
			outputs:    "output1,output2",
		},
	}

	actual := sd.performMethodExtraction(t, rootPlanName, cs)
	assert.Equal(test, 0, len(cs))
	assert.Equal(test, expected, actual)
	mock.AssertExpectationsForObjects(test, sdMock)
}

type field1 struct{}

type field2 struct{}

type extractChildMethods_struct struct {
	field1
	f field2
}

func TestExtractChildMethods(test *testing.T) {
	rootPlanName := "rootPlanName"
	cs := componentStack{}
	cs = cs.push("dummy")

	sd := newStructDisassembler()
	sdMock := &mockIStructDisassembler{}
	sd.itself = sdMock

	expected := []method{
		{
			owningType: "owningType1",
			name:       "name",
			arguments:  "arguments",
			outputs:    "output1,output2",
		},
		{
			owningType: "owningType2",
			name:       "name",
			arguments:  "arguments",
			outputs:    "output1,output2",
		},
	}

	sdMock.On("performMethodExtraction", mock.Anything, mock.Anything, mock.Anything).
		Return(func(t reflect.Type, rootPlanNameIn string, csIn componentStack) []method {
			assert.Equal(test, reflect.TypeOf(field1{}), t, "performMethodExtraction must only be called on embedded field1, not non-embedded field2")
			assert.Equal(test, rootPlanName, rootPlanNameIn)
			if assert.Equal(test, 1, len(csIn)) {
				assert.Equal(test, "dummy", csIn[0])
			}

			return expected
		}).Times(2)

	test.Run("non-struct type", func(t *testing.T) {
		actual := sd.extractChildMethods(reflect.TypeOf(true), rootPlanName, cs)
		assert.Equal(test, []method(nil), actual)
	})

	test.Run("non-pointer type", func(t *testing.T) {
		actual := sd.extractChildMethods(reflect.TypeOf(extractChildMethods_struct{}), rootPlanName, cs)
		assert.Equal(test, expected, actual)
	})

	test.Run("pointer type", func(t *testing.T) {
		actual := sd.extractChildMethods(reflect.TypeOf(&extractChildMethods_struct{}), rootPlanName, cs)
		assert.Equal(test, expected, actual)
	})

	mock.AssertExpectationsForObjects(test, sdMock)
}

type extractOwnMethods_struct struct{}

func (extractOwnMethods_struct) Do()   {}
func (extractOwnMethods_struct) Does() {}

func TestExtractOwnMethods(test *testing.T) {
	defer func(original func(rm reflect.Method, ignoreFirstReceiverArgument bool) method) {
		extractMethodDetails = original
	}(extractMethodDetails)

	rootPlanName := "rootPlanName"
	cs := componentStack{}

	sd := newStructDisassembler()
	sdMock := &mockIStructDisassembler{}
	sd.itself = sdMock

	extractMethodDetails = func(rm reflect.Method, ignoreFirstReceiverArgument bool) method {
		assert.True(test, ignoreFirstReceiverArgument)

		return method{
			name: "dummy",
		}
	}

	sdMock.On("addAvailableMethod", mock.Anything, mock.Anything, mock.Anything).
		Return(func(rootPlanNameIn string, csIn componentStack, m method) {
			assert.Equal(test, rootPlanName, rootPlanNameIn)
			assert.Equal(test, cs, csIn)

			expected := method{
				owningType: "github.com/jamestrandung/go-cte-117/cte/extractOwnMethods_struct",
				name:       "dummy",
			}

			assert.Equal(test, expected, m)
		}).Times(2)

	test.Run("hoisted methods do not contain owned methods", func(t *testing.T) {
		expected := []method{
			{
				owningType: "github.com/jamestrandung/go-cte-117/cte/extractOwnMethods_struct",
				name:       "dummy",
			},
			{
				owningType: "github.com/jamestrandung/go-cte-117/cte/extractOwnMethods_struct",
				name:       "dummy",
			},
		}

		actual := sd.extractOwnMethods(reflect.TypeOf(extractOwnMethods_struct{}), rootPlanName, cs, []method{})

		assert.Equal(test, expected, actual)
	})

	test.Run("hoisted methods contain owned methods", func(t *testing.T) {
		hoistedMethods := []method{
			{
				name: "dummy",
			},
		}

		actual := sd.extractOwnMethods(reflect.TypeOf(extractOwnMethods_struct{}), rootPlanName, cs, hoistedMethods)

		assert.Equal(test, []method(nil), actual)
	})

	mock.AssertExpectationsForObjects(test, sdMock)
}
