package cte

import (
	"fmt"
	"reflect"
	"strings"
)

type methodLocation struct {
	rootPlanName string
	componentStack
}

type method struct {
	owningType string
	name       string
	arguments  string // comma-separated argument types
	outputs    string // comma-separated argument types
}

func (m method) hasSameSignature(other method) bool {
	return m.name == other.name &&
		m.arguments == other.arguments &&
		m.outputs == other.outputs
}

func (m method) String() string {
	methodSignature := func() string {
		if strings.Contains(m.outputs, ",") {
			return fmt.Sprintf("%s(%s) (%s)", m.name, m.arguments, m.outputs)
		}

		if m.outputs == "" {
			return fmt.Sprintf("%s(%s)", m.name, m.arguments)
		}

		return fmt.Sprintf("%s(%s) %s", m.name, m.arguments, m.outputs)
	}()

	if m.owningType != "" {
		return m.owningType + "." + methodSignature
	}

	return methodSignature
}

var extractMethodDetails = func(rm reflect.Method, ignoreFirstReceiverArgument bool) method {
	var arguments []string
	for i := 0; i < rm.Type.NumIn(); i++ {
		if i == 0 && ignoreFirstReceiverArgument {
			continue
		}

		arguments = append(arguments, rm.Type.In(i).String())
	}

	var outputs []string
	for i := 0; i < rm.Type.NumOut(); i++ {
		outputs = append(outputs, rm.Type.Out(i).String())
	}

	return method{
		name:      rm.Name,
		arguments: strings.Join(arguments, ","),
		outputs:   strings.Join(outputs, ","),
	}
}

type structDisassembler struct {
	availableMethods             map[string]methodSet
	methodsAvailableMoreThanOnce methodSet
	methodLocations              map[method][]methodLocation
}

func newStructDisassembler() structDisassembler {
	return structDisassembler{
		availableMethods:             make(map[string]methodSet),
		methodsAvailableMoreThanOnce: make(methodSet),
		methodLocations:              make(map[method][]methodLocation),
	}
}

func (sd structDisassembler) isAvailableMoreThanOnce(m method) bool {
	return sd.methodsAvailableMoreThanOnce.has(m)
}

func (sd structDisassembler) findMethodLocations(ms methodSet, rootPlanName string) []string {
	var methodLocations []string
	for _, m := range ms.items() {
		for _, ml := range sd.methodLocations[m] {
			if ml.rootPlanName == rootPlanName {
				methodLocations = append(methodLocations, ml.componentStack.String())
			}
		}
	}

	return methodLocations
}

func (sd structDisassembler) extractAvailableMethods(t reflect.Type) []method {
	var cs componentStack
	return performMethodExtraction(sd, t, extractFullNameFromType(t), cs)
}

var addAvailableMethod = func(sd structDisassembler, rootPlanName string, cs componentStack, m method) {
	ms, ok := sd.availableMethods[m.name]
	if !ok {
		ms = make(methodSet)
		sd.availableMethods[m.name] = ms
	}

	// Even if a method is registered twice, they will be declared at different locations.
	// Hence, the provided component stack will never be the same and thus, we must record
	// the location before checking for duplicate.
	sd.methodLocations[m] = append(sd.methodLocations[m], methodLocation{
		rootPlanName:   rootPlanName,
		componentStack: cs.clone(),
	})

	if ms.has(m) {
		sd.methodsAvailableMoreThanOnce.add(m)
		return
	}

	ms.add(m)
}

var performMethodExtraction func(sd structDisassembler, t reflect.Type, rootPlanName string, cs componentStack) []method
var extractChildMethods func(sd structDisassembler, t reflect.Type, rootPlanName string, cs componentStack) []method
var extractOwnMethods func(sd structDisassembler, t reflect.Type, rootPlanName string, cs componentStack, hoistedMethods []method) []method

func init() {
	performMethodExtraction = func(sd structDisassembler, t reflect.Type, rootPlanName string, cs componentStack) []method {
		cs = cs.push(extractFullNameFromType(t))
		defer func() {
			cs = cs.pop()
		}()

		hoistedMethods := extractChildMethods(sd, t, rootPlanName, cs)
		ownMethods := extractOwnMethods(sd, t, rootPlanName, cs, hoistedMethods)

		var allMethods []method
		allMethods = append(allMethods, hoistedMethods...)
		allMethods = append(allMethods, ownMethods...)

		return allMethods
	}

	extractChildMethods = func(sd structDisassembler, t reflect.Type, rootPlanName string, cs componentStack) []method {
		var hoistedMethods []method

		actualType := t
		if actualType.Kind() == reflect.Pointer {
			actualType = t.Elem()
		}

		if actualType.Kind() == reflect.Struct {
			for i := 0; i < actualType.NumField(); i++ {
				rf := actualType.Field(i)

				// Extract methods from embedded fields
				if rf.Anonymous {
					childMethods := performMethodExtraction(sd, rf.Type, rootPlanName, cs)
					hoistedMethods = append(hoistedMethods, childMethods...)
				}
			}
		}

		return hoistedMethods
	}

	extractOwnMethods = func(sd structDisassembler, t reflect.Type, rootPlanName string, cs componentStack, hoistedMethods []method) []method {
		var ownMethods []method

		for i := 0; i < t.NumMethod(); i++ {
			rm := t.Method(i)

			m := extractMethodDetails(rm, true)
			m.owningType = t.PkgPath() + "/" + t.Name()

			isHoistedMethod := func() bool {
				for j := 0; j < len(hoistedMethods); j++ {
					hm := hoistedMethods[j]

					if hm.hasSameSignature(m) {
						return true
					}
				}

				return false
			}()

			// If a parent struct has a method carrying the same signature as one
			// that is available in an embedded field, it means this is a hoisted
			// method or the parent is overriding the same method with its own
			// implementation. In either case, we can assume this is a hoisted
			// method as it doesn't make a difference to how we should validate
			// code templates.
			if isHoistedMethod {
				continue
			}

			ownMethods = append(ownMethods, m)
			addAvailableMethod(sd, rootPlanName, cs, m)
		}

		return ownMethods
	}
}

type methodSet map[method]struct{}

func (ms methodSet) add(m method) {
	ms[m] = struct{}{}
}

func (ms methodSet) has(m method) bool {
	_, ok := ms[m]
	return ok
}

func (ms methodSet) count() int {
	return len(ms)
}

func (ms methodSet) items() []method {
	methods := make([]method, 0, len(ms))

	for key, _ := range ms {
		methods = append(methods, key)
	}

	return methods
}

func (ms methodSet) String() string {
	methods := make([]string, 0, len(ms))

	for key, _ := range ms {
		methods = append(methods, key.String())
	}

	return strings.Join(methods, ", ")
}
