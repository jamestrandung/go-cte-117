package cte

import (
	"fmt"
	"reflect"
)

var (
	planType           = reflect.TypeOf((*Plan)(nil)).Elem()
	preHookType        = reflect.TypeOf((*Pre)(nil)).Elem()
	postHookType       = reflect.TypeOf((*Post)(nil)).Elem()
	resultType         = reflect.TypeOf(Result{})
	syncResultType     = reflect.TypeOf(SyncResult{})
	sideEffectType     = reflect.TypeOf(SideEffect{})
	syncSideEffectType = reflect.TypeOf(SyncSideEffect{})
)

type analyzedPlan struct {
	pType        reflect.Type
	isMasterPlan bool
	isSequential bool
	components   []parsedComponent
	preHooks     []preHook
	postHooks    []postHook
}

type parsedComponent struct {
	id            string
	fieldIdx      int
	fieldType     reflect.Type
	isSyncResult  bool
	requireSet    bool
	isPointerType bool
}

type preHook struct {
	hook     Pre
	metadata parsedMetadata
}

type postHook struct {
	hook     Post
	metadata parsedMetadata
}

type planAnalyzer struct {
	engine     Engine
	plan       Plan
	planValue  reflect.Value
	preHooks   []preHook
	postHooks  []postHook
	components []parsedComponent
}

func (pa *planAnalyzer) analyze() analyzedPlan {
	for i := 0; i < pa.planValue.NumField(); i++ {
		isPointerType, fieldType, fieldPointerType := extractFieldTypes(pa.planValue.Type().Field(i))

		fa := fieldAnalyzer{
			pa:               pa,
			fieldIdx:         i,
			isPointerType:    isPointerType,
			fieldType:        fieldType,
			fieldPointerType: fieldPointerType,
		}

		component, pre, post := fa.analyze()

		if component != nil {
			pa.components = append(pa.components, *component)
		}

		if pre != nil {
			pa.preHooks = append(pa.preHooks, *pre)
		}

		if post != nil {
			pa.postHooks = append(pa.postHooks, *post)
		}
	}

	_, isMasterPlan := pa.plan.(MasterPlan)

	return analyzedPlan{
		pType:        extractUnderlyingType(pa.planValue),
		isMasterPlan: isMasterPlan,
		isSequential: pa.plan.IsSequentialCTEPlan(),
		components:   pa.components,
		preHooks:     pa.preHooks,
		postHooks:    pa.postHooks,
	}
}

type fieldAnalyzer struct {
	pa               *planAnalyzer
	fieldIdx         int
	isPointerType    bool
	fieldType        reflect.Type
	fieldPointerType reflect.Type
}

func (fa *fieldAnalyzer) analyze() (*parsedComponent, *preHook, *postHook) {
	pre, post := fa.handleHooks()
	if pre != nil || post != nil {
		return nil, pre, post
	}

	if nestedPlan := fa.handleNestedPlan(); nestedPlan != nil {
		return nestedPlan, nil, nil
	}

	return fa.handleComputer(), nil, nil
}

func (fa *fieldAnalyzer) handleHooks() (*preHook, *postHook) {
	// Hook types might be embedded in a parent plan struct. Hence, we need to check if the type
	// is a hook but not a plan so that we don't register a plan as a hook.
	typeAndPointerTypeIsNotPlanType := !fa.fieldType.Implements(planType) && !fa.fieldPointerType.Implements(planType)

	// Hooks might be implemented with value or pointer receivers.
	isPreHookType := fa.fieldType.Implements(preHookType) || fa.fieldPointerType.Implements(preHookType)
	isPostHookType := fa.fieldType.Implements(postHookType) || fa.fieldPointerType.Implements(postHookType)

	if typeAndPointerTypeIsNotPlanType && (isPreHookType || isPostHookType) {
		// Call to Interface() returns a pointer value which is acceptable for
		// both scenarios where fieldType uses pointer or value receiver to
		// implement an interface
		hook := reflect.New(fa.fieldType).Interface()

		mp, ok := hook.(MetadataProvider)
		if !ok {
			panic(ErrMetadataMissing.Err(fa.fieldType))
		}

		if isPreHookType {
			return &preHook{
				hook:     hook.(Pre),
				metadata: extractMetadata(mp, false),
			}, nil
		}

		return nil, &postHook{
			hook:     hook.(Post),
			metadata: extractMetadata(mp, false),
		}
	}

	return nil, nil
}

func (fa *fieldAnalyzer) handleNestedPlan() *parsedComponent {
	if !fa.fieldPointerType.Implements(planType) {
		return nil
	}

	if fa.isPointerType {
		panic(ErrNestedPlanCannotBePointer.Err(fa.pa.planValue.Type(), fa.fieldType))
	}

	// Dynamically analyze nested plans
	fa.pa.engine.AnalyzePlan(reflect.New(fa.fieldType).Interface().(Plan))

	return &parsedComponent{
		id:       extractFullNameFromType(fa.fieldType),
		fieldIdx: fa.fieldIdx,
	}
}

func (fa *fieldAnalyzer) handleComputer() *parsedComponent {
	componentID := extractFullNameFromType(fa.fieldType)

	component := func() *parsedComponent {
		if fa.fieldType.ConvertibleTo(resultType) {
			// Both sequential & parallel plans can contain Result fields
			return &parsedComponent{
				id:            componentID,
				fieldIdx:      fa.fieldIdx,
				fieldType:     fa.fieldType,
				requireSet:    true,
				isPointerType: fa.isPointerType,
			}
		}

		if fa.fieldType.ConvertibleTo(syncResultType) {
			if !fa.pa.plan.IsSequentialCTEPlan() {
				panic(fmt.Errorf("parallel plan cannot contain SyncResult field: %s", extractShortName(componentID)))
			}

			return &parsedComponent{
				id:            componentID,
				fieldIdx:      fa.fieldIdx,
				fieldType:     fa.fieldType,
				isSyncResult:  true,
				requireSet:    true,
				isPointerType: fa.isPointerType,
			}
		}

		if fa.fieldType.ConvertibleTo(sideEffectType) || fa.fieldType.ConvertibleTo(syncSideEffectType) {
			return &parsedComponent{
				id:       componentID,
				fieldIdx: fa.fieldIdx,
			}
		}

		return nil
	}()

	if component == nil {
		return nil
	}

	// Dynamically register computers that are actually used in a plan
	mp, ok := reflect.New(fa.fieldType).Interface().(MetadataProvider)
	if !ok {
		panic(ErrMetadataMissing.Err(fa.fieldType))
	}

	fa.pa.engine.registerComputer(mp)

	return component
}
