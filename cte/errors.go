package cte

import (
	"errors"
	"fmt"
)

type formatErr struct {
	format string
}

func makeFormatErr(format string) formatErr {
	return formatErr{
		format: format,
	}
}

func (e formatErr) Err(v ...any) error {
	return fmt.Errorf(e.format, v...)
}

var (
	// ErrPlanExecutionEndingEarly can be thrown actively by clients to end plan execution early.
	// For example, a value was retrieved from cache and thus, there's no point executing the algo
	// to calculate this value anymore. The engine will swallow this error, end execution and then
	// return a nil error to clients.
	//
	// Note: If the ending plan is nested inside another plan, the outer plan will still continue
	// as usual.
	ErrPlanExecutionEndingEarly = errors.New("CTE-001: plan execution ending early")
	// ErrRootPlanExecutionEndingEarly can be thrown actively by clients to end plan execution
	// early. For example, a value was retrieved from cache and thus, there's no point executing
	// the algo to calculate this value anymore. The engine will swallow this error, end execution
	// and then return a nil error to clients.
	//
	// Note: If the ending plan is nested inside another plan, the outer plan will also end.
	ErrRootPlanExecutionEndingEarly = errors.New("CTE-0002: plan execution ending early from root")

	ErrPlanMustUsePointerReceiver = makeFormatErr("CTE-0003: %v is using value receiver, all plans must be implemented using pointer receiver")
	ErrPlanNotAnalyzed            = makeFormatErr("CTE-0004: %v has not been analyzed yet, call AnalyzePlan on it first")
	ErrNestedPlanCannotBePointer  = makeFormatErr("CTE-0005: %v has a nested plan called %v that is a pointer")

	ErrPlanNotMeetingInoutRequirements            = makeFormatErr("CTE-0006: %v does not implement the required in-out interface %v, problem found: %v. Component having problem: %v")
	ErrPlanMissingMethod                          = makeFormatErr("missing method: [%v]")
	ErrPlanHavingAmbiguousMethods                 = makeFormatErr("required method: [%v], found ambiguous methods: [%v], components carrying the ambiguous methods: [%v]")
	ErrPlanHavingSameMethodRegisteredMoreThanOnce = makeFormatErr("required method provided more than once by the same computer key: [%v], computer key locations: [%v]")
	ErrPlanHavingMethodButSignatureMismatched     = makeFormatErr("required method: [%v], found method with mismatched signature: [%v]")

	ErrInvalidComputerType = makeFormatErr("CTE-0007: %v is not a computer")
	ErrMetadataMissing     = makeFormatErr("CTE-0008: metadata is missing for %v, it must implement the MetadataProvider interface")
	ErrNilMetadata         = makeFormatErr("CTE-0009: metadata is nil for %v")
	ErrComputerMetaMissing = makeFormatErr("CTE-0010: computer meta is missing in %v")
	ErrInoutMetaMissing    = makeFormatErr("CTE-0011: inout meta is missing in %v")
)
