package cte

import (
	"reflect"
	"strings"
)

func isComplete(e Engine, planValue reflect.Value) error {
	planName := extractFullNameFromType(extractUnderlyingType(planValue))

	sd := newStructDisassembler()
	sd.extractAvailableMethods(planValue.Type())

	var cs componentStack

	var verifyFn func(planName string, curPlanValue reflect.Value) error
	verifyFn = func(planName string, curPlanValue reflect.Value) error {
		ap := e.findAnalyzedPlan(planName, curPlanValue)

		cs = cs.Push(planName)
		defer func() {
			cs = cs.Pop()
		}()

		for _, h := range ap.preHooks {
			expectedInout, ok := h.metadata.getInoutInterface()
			if !ok {
				return ErrInoutMetaMissing.Err(reflect.TypeOf(h.hook))
			}

			err := isInterfaceSatisfied(sd, expectedInout)
			if err != nil {
				cs = cs.Push(reflect.TypeOf(h.hook).Name())
				return ErrPlanNotMeetingInoutRequirements.Err(planValue.Type(), expectedInout, err.Error(), cs)
			}
		}

		for _, component := range ap.components {
			if c, ok := e.computers[component.id]; ok {
				expectedInout, ok := c.metadata.getInoutInterface()
				if !ok {
					return ErrInoutMetaMissing.Err(component.id)
				}

				err := isInterfaceSatisfied(sd, expectedInout)
				if err != nil {
					cs = cs.Push(component.id)
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

			err := isInterfaceSatisfied(sd, expectedInout)
			if err != nil {
				cs = cs.Push(reflect.TypeOf(h.hook).Name())
				return ErrPlanNotMeetingInoutRequirements.Err(planValue.Type(), expectedInout, err.Error(), cs)
			}
		}

		return nil
	}

	return verifyFn(planName, planValue)
}

func isInterfaceSatisfied(sd structDisassembler, expectedInterface reflect.Type) error {
	for i := 0; i < expectedInterface.NumMethod(); i++ {
		rm := expectedInterface.Method(i)

		requiredMethod := extractMethodDetails(rm, false)

		ms, ok := sd.availableMethods[requiredMethod.name]
		if !ok {
			return ErrPlanMissingMethod.Err(requiredMethod)
		}

		if ms.count() > 1 {
			return ErrPlanHavingAmbiguousMethods.Err(requiredMethod, ms)
		}

		foundMethod := ms.items()[0]

		if !foundMethod.hasSameSignature(requiredMethod) {
			return ErrPlanHavingMethodButSignatureMismatched.Err(requiredMethod, foundMethod)
		}

		if sd.isAvailableMoreThanOnce(foundMethod) {
			return ErrPlanHavingSameMethodRegisteredMoreThanOnce.Err(foundMethod)
		}
	}

	return nil
}

type componentStack []string

func (s componentStack) Push(componentName string) componentStack {
	return append(s, componentName)
}

func (s componentStack) Pop() componentStack {
	return s[0 : len(s)-1]
}

func (s componentStack) String() string {
	return strings.Join(s, " >> ")
}
