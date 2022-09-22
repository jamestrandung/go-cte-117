package cte

import (
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
	loaders      []loadFn
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

//go:generate mockery --name iPlanAnalyzer --case=underscore --inpackage
type iPlanAnalyzer interface {
	extractLoaders() []loadFn
}

type planAnalyzer struct {
	itself     iPlanAnalyzer
	engine     iEngine
	plan       Plan
	planValue  reflect.Value
	preHooks   []preHook
	postHooks  []postHook
	components []parsedComponent
}

func newPlanAnalyzer(e iEngine, p Plan, pValue reflect.Value) *planAnalyzer {
	if pValue.Kind() == reflect.Pointer {
		pValue = pValue.Elem()
	}

	result := &planAnalyzer{
		engine:    e,
		plan:      p,
		planValue: pValue,
	}

	result.itself = result

	return result
}

func (pa *planAnalyzer) analyze() analyzedPlan {
	for i := 0; i < pa.planValue.NumField(); i++ {
		isPointerType, fieldType, fieldPointerType := extractFieldTypes(pa.planValue.Type().Field(i))

		fa := newFieldAnalyzer(pa, i, isPointerType, fieldType, fieldPointerType)

		component, pre, post := fa.itself.analyze()

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

	loaders := pa.itself.extractLoaders()

	return analyzedPlan{
		pType:        extractUnderlyingType(pa.planValue),
		isMasterPlan: isMasterPlan,
		isSequential: pa.plan.IsSequentialCTEPlan(),
		components:   pa.components,
		loaders:      loaders,
		preHooks:     pa.preHooks,
		postHooks:    pa.postHooks,
	}
}

func (pa *planAnalyzer) extractLoaders() []loadFn {
	// Loaders have to maintain the same index with the corresponding component.
	// Hence, cannot simply use append() on an empty slice.
	loaders := make([]loadFn, len(pa.components))

	foundLoader := false
	for idx, component := range pa.components {
		if c, ok := pa.engine.getComputer(component.id); ok {
			if c.computer.loadFn != nil {
				loaders[idx] = c.computer.loadFn
				foundLoader = true
			}
		}
	}

	if !foundLoader {
		var tmp []loadFn
		return tmp
	}

	return loaders
}

//go:generate mockery --name iFieldAnalyzer --case=underscore --inpackage
type iFieldAnalyzer interface {
	analyze() (*parsedComponent, *preHook, *postHook)
	handleHooks() (*preHook, *postHook)
	handleNestedPlan() *parsedComponent
	handleComputer() *parsedComponent
	createComputerComponent(componentID string) *parsedComponent
}

type fieldAnalyzer struct {
	itself           iFieldAnalyzer
	pa               *planAnalyzer
	fieldIdx         int
	isPointerType    bool
	fieldType        reflect.Type
	fieldPointerType reflect.Type
}

var newFieldAnalyzer = func(
	pa *planAnalyzer,
	fieldIdx int,
	isPointerType bool,
	fieldType reflect.Type,
	fieldPointerType reflect.Type,
) *fieldAnalyzer {
	fa := &fieldAnalyzer{
		pa:               pa,
		fieldIdx:         fieldIdx,
		isPointerType:    isPointerType,
		fieldType:        fieldType,
		fieldPointerType: fieldPointerType,
	}

	fa.itself = fa

	return fa
}

func (fa *fieldAnalyzer) analyze() (*parsedComponent, *preHook, *postHook) {
	pre, post := fa.itself.handleHooks()
	if pre != nil || post != nil {
		return nil, pre, post
	}

	if nestedPlan := fa.itself.handleNestedPlan(); nestedPlan != nil {
		return nestedPlan, nil, nil
	}

	return fa.itself.handleComputer(), nil, nil
}

func (fa *fieldAnalyzer) handleHooks() (*preHook, *postHook) {
	// Hook types might be embedded in a parent plan struct. Hence, we need to check if the type
	// is a hook but not a plan so that we don't register a plan as a hook.
	typeOrPointerTypeIsPlanType := fa.fieldType.Implements(planType) || fa.fieldPointerType.Implements(planType)

	// Hooks might be implemented with value or pointer receivers.
	isPreHookType := fa.fieldType.Implements(preHookType) || fa.fieldPointerType.Implements(preHookType)
	isPostHookType := fa.fieldType.Implements(postHookType) || fa.fieldPointerType.Implements(postHookType)

	if typeOrPointerTypeIsPlanType || (!isPreHookType && !isPostHookType) {
		return nil, nil
	}

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

	component := fa.itself.createComputerComponent(componentID)
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

func (fa *fieldAnalyzer) createComputerComponent(componentID string) *parsedComponent {
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
			panic(ErrParallelPlanCannotContainSyncResult.Err(fa.pa.planValue.Type(), extractShortName(componentID)))
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

	if fa.fieldType.ConvertibleTo(sideEffectType) {
		return &parsedComponent{
			id:       componentID,
			fieldIdx: fa.fieldIdx,
		}
	}

	if fa.fieldType.ConvertibleTo(syncSideEffectType) {
		if !fa.pa.plan.IsSequentialCTEPlan() {
			panic(ErrParallelPlanCannotContainSyncSideEffect.Err(fa.pa.planValue.Type(), extractShortName(componentID)))
		}

		return &parsedComponent{
			id:       componentID,
			fieldIdx: fa.fieldIdx,
		}
	}

	return nil
}
