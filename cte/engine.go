package cte

import (
    "context"
    "fmt"
    "reflect"
    "runtime/debug"

    "github.com/jamestrandung/go-concurrency-117/async"
    "golang.org/x/sync/errgroup"
)

type registeredComputer struct {
    computer delegatingComputer
    metadata parsedMetadata
}

//go:generate mockery --name iEngine --case=underscore --inpackage
type iEngine interface {
    AnalyzePlan(p Plan)
    registerComputer(mp MetadataProvider)
    findAnalyzedPlan(planName string, curPlanValue reflect.Value) analyzedPlan
    getComputer(componentID string) (registeredComputer, bool)
    getPlan(planName string) (analyzedPlan, bool)
}

type Engine struct {
    computers map[string]registeredComputer
    plans     map[string]analyzedPlan
}

func NewEngine() Engine {
    return Engine{
        computers: make(map[string]registeredComputer),
        plans:     make(map[string]analyzedPlan),
    }
}

func (e Engine) AnalyzePlan(p Plan) {
    planName := extractFullNameFromValue(p)
    if _, ok := e.plans[planName]; ok {
        return
    }

    val := reflect.ValueOf(p)
    if val.Kind() != reflect.Pointer {
        panic(ErrPlanMustUsePointerReceiver.Err(reflect.TypeOf(p)))
    }

    pa := newPlanAnalyzer(e, p, val)
    ap := pa.analyze()

    e.plans[planName] = ap
}

func (e Engine) registerComputer(mp MetadataProvider) {
    computerID := extractFullNameFromValue(mp)
    if _, ok := e.computers[computerID]; ok {
        return
    }

    metadata := extractMetadata(mp, true)

    computerType := func() reflect.Type {
        cType, ok := metadata.getComputerType()
        if !ok {
            panic(ErrComputerMetaMissing.Err(reflect.TypeOf(mp)))
        }

        return extractNonPointerType(cType)
    }()

    computer := reflect.New(computerType).Interface()
    e.computers[computerID] = registeredComputer{
        computer: newDelegatingComputer(computer),
        metadata: metadata,
    }
}

func (e Engine) ExecuteMasterPlan(ctx context.Context, p MasterPlan) error {
    // Plan implementations always use pointer receivers.
    // Should be safe to extract value.
    planValue := reflect.ValueOf(p).Elem()

    planName := extractFullNameFromType(planValue.Type())
    if err := e.doExecutePlan(ctx, planName, p, planValue, p.IsSequentialCTEPlan()); err != nil {
        return swallowErrPlanExecutionEndingEarly(err)
    }

    return nil
}

func (e Engine) doExecutePlan(ctx context.Context, planName string, p MasterPlan, curPlanValue reflect.Value, isSequential bool) error {
    ap := e.findAnalyzedPlan(planName, curPlanValue)

    for _, h := range ap.preHooks {
        if err := h.hook.PreExecute(p); err != nil {
            return err
        }
    }

    err := func() error {
        if isSequential {
            return e.doExecuteSync(ctx, p, curPlanValue, ap.loaders, ap.components)
        }

        return e.doExecuteAsync(ctx, p, curPlanValue, ap.components)
    }()

    if err != nil {
        return err
    }

    for _, h := range ap.postHooks {
        if err := h.hook.PostExecute(p); err != nil {
            return err
        }
    }

    return nil
}

func (e Engine) doExecuteComputer(ctx context.Context, c delegatingComputer, p MasterPlan, loadingData LoadingData) (interface{}, error) {
    result, err := c.Compute(ctx, p, loadingData)
    if tep, ok := result.(toExecutePlan); ok {
        if err != nil {
            return tep.mp, err
        }

        return tep.mp, tep.mp.Execute(ctx)
    }

    return result, err
}

func (e Engine) doConcurrentLoading(ctx context.Context, p MasterPlan, componentCount int, loaders []loadFn) []LoadingData {
    // Data has to be loaded at the same index with the corresponding component
    loadingData := make([]LoadingData, componentCount)
    if len(loaders) == 0 {
        return loadingData
    }

    tasks := make([]async.SilentTask, 0, len(loaders))
    for idx, loader := range loaders {
        if loader == nil {
            continue
        }

        i := idx
        l := loader

        tasks = append(
            tasks,
            async.NewSilentTask(
                func(taskCtx context.Context) error {
                    data, err := l(taskCtx, p)

                    loadingData[i] = LoadingData{
                        Data: data,
                        Err:  err,
                    }

                    return nil
                },
            ),
        )
    }

    async.ForkJoinST(ctx, tasks)

    return loadingData
}

func (e Engine) doExecuteSync(
    ctx context.Context,
    p MasterPlan,
    curPlanValue reflect.Value,
    loaders []loadFn,
    components []parsedComponent,
) error {
    loadingData := e.doConcurrentLoading(ctx, p, len(components), loaders)

    for idx, component := range components {
        if c, ok := e.computers[component.id]; ok {
            result, err := func() (result any, err error) {
                defer func() {
                    if r := recover(); r != nil {
                        err = fmt.Errorf("panic executing sync task: %v \n %s", r, debug.Stack())
                    }
                }()

                return e.doExecuteComputer(ctx, c.computer, p, loadingData[idx])
            }()

            // Register Result/SyncResult in a sequential plan's field
            if component.requireSet {
                field := curPlanValue.Field(component.fieldIdx)

                casted := func() reflect.Value {
                    if component.isSyncResult {
                        return reflect.ValueOf(newSyncResult(result)).Convert(component.fieldType)
                    }

                    task := async.Completed(result, err)
                    return reflect.ValueOf(newResult(task)).Convert(component.fieldType)
                }()

                resultTakingIntoAccountPointerType := func() reflect.Value {
                    if !component.isPointerType {
                        return casted
                    }

                    rv := reflect.New(component.fieldType)
                    rv.Elem().Set(casted)

                    return rv
                }()

                field.Set(resultTakingIntoAccountPointerType)
            }

            if err != nil {
                return err
            }

            continue
        }

        // Nested plan gets executed synchronously
        if ap, ok := e.plans[component.id]; ok {
            // Nested plan is always a value, never a pointer. Hence, no need to call Elem().
            nestedPlanValue := curPlanValue.Field(component.fieldIdx)

            err := e.doExecutePlan(ctx, component.id, p, nestedPlanValue, ap.isSequential)
            if err != nil && err != ErrPlanExecutionEndingEarly {
                return err
            }
        }
    }

    return nil
}

func (e Engine) doExecuteAsync(ctx context.Context, p MasterPlan, curPlanValue reflect.Value, components []parsedComponent) error {
    tasks := make([]async.SilentTask, 0, len(components))
    for _, component := range components {
        componentID := component.id

        if c, ok := e.computers[componentID]; ok {
            task := async.NewTask(
                func(taskCtx context.Context) (interface{}, error) {
                    data, err := c.computer.Load(taskCtx, p)

                    return e.doExecuteComputer(
                        taskCtx, c.computer, p, LoadingData{
                            Data: data,
                            Err:  err,
                        },
                    )
                },
            )

            tasks = append(tasks, task)

            // Register Result in a parallel plan's field
            if component.requireSet {
                field := curPlanValue.Field(component.fieldIdx)

                resultTakingIntoAccountPointerType := func() reflect.Value {
                    casted := reflect.ValueOf(newResult(task)).Convert(component.fieldType)
                    if !component.isPointerType {
                        return casted
                    }

                    rv := reflect.New(component.fieldType)
                    rv.Elem().Set(casted)

                    return rv
                }()

                field.Set(resultTakingIntoAccountPointerType)
            }

            continue
        }

        // Nested plan gets executed asynchronously by wrapping it inside a task
        if ap, ok := e.plans[componentID]; ok {
            // Nested plan is always a value, never a pointer. Hence, no need to call Elem().
            nestedPlanValue := curPlanValue.Field(component.fieldIdx)

            task := async.NewSilentTask(
                func(taskCtx context.Context) error {
                    err := e.doExecutePlan(taskCtx, componentID, p, nestedPlanValue, ap.isSequential)
                    if err != nil && err != ErrPlanExecutionEndingEarly {
                        return err
                    }

                    return nil
                },
            )

            tasks = append(tasks, task)
        }
    }

    g, groupCtx := errgroup.WithContext(ctx)
    for _, task := range tasks {
        t := task
        g.Go(
            func() error {
                return t.ExecuteSync(groupCtx).Error()
            },
        )
    }

    return g.Wait()
}

func (e Engine) VerifyConfigurations() error {
    for _, p := range e.plans {
        if p.isMasterPlan {
            rp := reflect.New(p.pType)

            if err := newCompletenessValidator(e, rp).validate(); err != nil {
                return err
            }
        }
    }

    return nil
}

func (e Engine) findExistingPlanOrCreate(planName string) analyzedPlan {
    if existing, ok := e.plans[planName]; ok {
        return existing
    }

    return analyzedPlan{}
}

func (e Engine) findAnalyzedPlan(planName string, curPlanValue reflect.Value) analyzedPlan {
    if len(planName) == 0 {
        panic(ErrPlanNotAnalyzed.Err(curPlanValue.Type()))
    }

    ap, ok := e.plans[planName]
    if !ok || len(ap.components) == 0 {
        panic(ErrPlanNotAnalyzed.Err(planName))
    }

    return ap
}

func (e Engine) getComputer(componentID string) (registeredComputer, bool) {
    c, ok := e.computers[componentID]
    return c, ok
}

func (e Engine) getPlan(planName string) (analyzedPlan, bool) {
    p, ok := e.plans[planName]
    return p, ok
}
