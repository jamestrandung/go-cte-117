package cte

import (
	"reflect"
	"strings"
)

func isComplete(e Engine, planValue reflect.Value) error {
	sd := newStructDisassembler()
	sd.extractAvailableMethods(planValue.Type())

	var cs componentStack
	rootPlanName := extractFullNameFromType(planValue.Type())

	var verifyFn func(planName string, curPlanValue reflect.Value) error
	verifyFn = func(planName string, curPlanValue reflect.Value) error {
		ap := e.findAnalyzedPlan(planName, curPlanValue)

		cs = cs.push(planName)
		defer func() {
			cs = cs.pop()
		}()

		for _, h := range ap.preHooks {
			expectedInout, ok := h.metadata.getInoutInterface()
			if !ok {
				return ErrInoutMetaMissing.Err(reflect.TypeOf(h.hook))
			}

			err := isInterfaceSatisfied(sd, expectedInout, rootPlanName)
			if err != nil {
				cs = cs.push(reflect.TypeOf(h.hook).Name())
				return ErrPlanNotMeetingInoutRequirements.Err(planValue.Type(), expectedInout, err.Error(), cs)
			}
		}

		for _, component := range ap.components {
			if c, ok := e.computers[component.id]; ok {
				expectedInout, ok := c.metadata.getInoutInterface()
				if !ok {
					return ErrInoutMetaMissing.Err(component.id)
				}

				err := isInterfaceSatisfied(sd, expectedInout, rootPlanName)
				if err != nil {
					cs = cs.push(component.id)
					return ErrPlanNotMeetingInoutRequirements.Err(planValue.Type(), expectedInout, err.Error(), cs)
				}
			}

			if _, ok := e.plans[component.id]; ok {
				nestedPlanValue := func() reflect.Value {
					if curPlanValue.Kind() == reflect.Pointer {
						return curPlanValue.Elem().Field(component.fieldIdx)
					}

					return curPlanValue.Field(component.fieldIdx)
				}()

				if err := verifyFn(component.id, nestedPlanValue); err != nil {
					return err
				}
			}
		}

		for _, h := range ap.postHooks {
			expectedInout, ok := h.metadata.getInoutInterface()
			if !ok {
				return ErrInoutMetaMissing.Err(reflect.TypeOf(h.hook))
			}

			err := isInterfaceSatisfied(sd, expectedInout, rootPlanName)
			if err != nil {
				cs = cs.push(reflect.TypeOf(h.hook).Name())
				return ErrPlanNotMeetingInoutRequirements.Err(planValue.Type(), expectedInout, err.Error(), cs)
			}
		}

		return nil
	}

	return verifyFn(rootPlanName, planValue)
}

func isInterfaceSatisfied(sd structDisassembler, expectedInterface reflect.Type, rootPlanName string) error {
	for i := 0; i < expectedInterface.NumMethod(); i++ {
		rm := expectedInterface.Method(i)

		requiredMethod := extractMethodDetails(rm, false)

		ms, ok := sd.availableMethods[requiredMethod.name]
		if !ok {
			return ErrPlanMissingMethod.Err(requiredMethod)
		}

		if ms.count() > 1 {
			methodLocations := sd.findMethodLocations(ms, rootPlanName)
			return ErrPlanHavingAmbiguousMethods.Err(requiredMethod, ms, strings.Join(methodLocations, "; "))
		}

		foundMethod := ms.items()[0]

		if !foundMethod.hasSameSignature(requiredMethod) {
			return ErrPlanHavingMethodButSignatureMismatched.Err(requiredMethod, foundMethod)
		}

		if sd.isAvailableMoreThanOnce(foundMethod) {
			methodLocations := sd.findMethodLocations(ms, rootPlanName)
			return ErrPlanHavingSameMethodRegisteredMoreThanOnce.Err(foundMethod, strings.Join(methodLocations, "; "))
		}
	}

	return nil
}

type componentStack []string

func (s componentStack) push(componentName string) componentStack {
	return append(s, componentName)
}

func (s componentStack) pop() componentStack {
	return s[0 : len(s)-1]
}

func (s componentStack) clone() componentStack {
	result := make([]string, 0, len(s))

	for _, c := range s {
		result = append(result, c)
	}

	return result
}

func (s componentStack) String() string {
	return strings.Join(s, " >> ")
}
