package cte

import (
	"reflect"
	"strings"
)

//go:generate mockery --name iCompletenessValidator --case=underscore --inpackage
type iCompletenessValidator interface {
	doValidate(planName string, cs componentStack, curPlanValue reflect.Value) error
	verifyComponentCompleteness(pm parsedMetadata, cs componentStack, componentID string, planType string) error
	isInterfaceSatisfied(expectedInterface reflect.Type) error
}

type completenessValidator struct {
	itself       iCompletenessValidator
	engine       iEngine
	planValue    reflect.Value
	rootPlanName string
	sd           *structDisassembler
}

func newCompletenessValidator(engine Engine, planValue reflect.Value) *completenessValidator {
	rootPlanName := extractFullNameFromType(planValue.Type())

	sd := newStructDisassembler()
	sd.extractAvailableMethods(planValue.Type())

	result := &completenessValidator{
		engine:       engine,
		planValue:    planValue,
		rootPlanName: rootPlanName,
		sd:           sd,
	}

	result.itself = result

	return result
}

func (v *completenessValidator) validate() error {
	var cs componentStack
	return v.itself.doValidate(v.rootPlanName, cs, v.planValue)
}

func (v *completenessValidator) doValidate(planName string, cs componentStack, curPlanValue reflect.Value) error {
	ap := v.engine.findAnalyzedPlan(planName, curPlanValue)

	cs = cs.push(planName)
	defer func() {
		cs = cs.pop()
	}()

	for _, h := range ap.preHooks {
		err := v.itself.verifyComponentCompleteness(h.metadata, cs, reflect.TypeOf(h.hook).String(), v.planValue.Type().String())
		if err != nil {
			return err
		}
	}

	for _, component := range ap.components {
		if c, ok := v.engine.getComputer(component.id); ok {
			err := v.itself.verifyComponentCompleteness(c.metadata, cs, component.id, v.planValue.Type().String())
			if err != nil {
				return err
			}

			continue
		}

		if _, ok := v.engine.getPlan(component.id); ok {
			nestedPlanValue := func() reflect.Value {
				if curPlanValue.Kind() == reflect.Pointer {
					return curPlanValue.Elem().Field(component.fieldIdx)
				}

				return curPlanValue.Field(component.fieldIdx)
			}()

			if err := v.itself.doValidate(component.id, cs, nestedPlanValue); err != nil {
				return err
			}
		}
	}

	for _, h := range ap.postHooks {
		err := v.itself.verifyComponentCompleteness(h.metadata, cs, reflect.TypeOf(h.hook).String(), v.planValue.Type().String())
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *completenessValidator) verifyComponentCompleteness(pm parsedMetadata, cs componentStack, componentID string, planType string) error {
	expectedInout, ok := pm.getInoutInterface()
	if !ok {
		return ErrInoutMetaMissing.Err(componentID)
	}

	err := v.itself.isInterfaceSatisfied(expectedInout)
	if err != nil {
		cs = cs.push(componentID)
		return ErrPlanNotMeetingInoutRequirements.Err(planType, expectedInout, err.Error(), cs)
	}

	return nil
}

func (v *completenessValidator) isInterfaceSatisfied(expectedInterface reflect.Type) error {
	for i := 0; i < expectedInterface.NumMethod(); i++ {
		rm := expectedInterface.Method(i)

		requiredMethod := extractMethodDetails(rm, false)

		ms, ok := v.sd.itself.findAvailableMethods(requiredMethod.name)
		if !ok {
			return ErrPlanMissingMethod.Err(requiredMethod)
		}

		if ms.count() > 1 {
			methodLocations := v.sd.itself.findMethodLocations(ms, v.rootPlanName)
			return ErrPlanHavingAmbiguousMethods.Err(requiredMethod, ms, strings.Join(methodLocations, "; "))
		}

		foundMethod := ms.items()[0]

		if !foundMethod.hasSameSignature(requiredMethod) {
			return ErrPlanHavingMethodButSignatureMismatched.Err(requiredMethod, foundMethod)
		}

		if v.sd.itself.isAvailableMoreThanOnce(foundMethod) {
			methodLocations := v.sd.itself.findMethodLocations(ms, v.rootPlanName)
			return ErrPlanHavingSameMethodRegisteredMoreThanOnce.Err(foundMethod, strings.Join(methodLocations, "; "))
		}
	}

	return nil
}
