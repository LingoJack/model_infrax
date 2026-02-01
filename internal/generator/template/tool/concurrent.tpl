package tool

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

// =============================================================================
// Goroutine 安全启动工具 - 自动处理 panic recover
// =============================================================================

// RecoveryOption 定义 Goroutine 恢复选项（使用位运算）
//
// 通过位运算可以灵活组合多个选项，例如：
//   - RecoveryOptionPrintStack | RecoveryOptionLogError
//
// 使用示例:
//
//	opts := RecoveryOptionPrintStack | RecoveryOptionLogError
//	Go(func() { /* your code */ }, opts)
type RecoveryOption uint32

const (
	// RecoveryOptionNone 默认选项，不做任何特殊配置
	RecoveryOptionNone RecoveryOption = 0

	// RecoveryOptionPrintStack 打印堆栈信息
	// 适用于需要详细排查问题的场景
	RecoveryOptionPrintStack RecoveryOption = 1 << 0 // 1

	// RecoveryOptionLogError 使用 logger 记录错误
	// 适用于生产环境需要记录错误日志的场景
	RecoveryOptionLogError RecoveryOption = 1 << 1 // 2

	// RecoveryOptionSilent 静默模式，不输出任何信息
	// 适用于已有其他错误处理机制的场景
	RecoveryOptionSilent RecoveryOption = 1 << 2 // 4
)

// Has 检查是否包含指定选项
//
// 使用示例:
//
//	opts := RecoveryOptionPrintStack | RecoveryOptionLogError
//	if opts.Has(RecoveryOptionPrintStack) {
//	    fmt.Println("将打印堆栈信息")
//	}
func (o RecoveryOption) Has(option RecoveryOption) bool {
	return o&option != 0
}

// GetDefaultRecoveryOptions 获取推荐的默认选项组合
//
// 该函数返回一组常用的选项组合，适用于大多数场景。
// 默认会打印堆栈并记录错误日志。
//
// 返回值:
//   - options: 推荐的默认选项组合
//
// 使用示例:
//
//	opts := GetDefaultRecoveryOptions()
//	Go(func() { /* your code */ }, opts)
func GetDefaultRecoveryOptions() RecoveryOption {
	return RecoveryOptionPrintStack | RecoveryOptionLogError
}

// PanicHandler panic 处理回调函数类型
//
// 参数:
//   - panicValue: panic 时的值
//   - stack: 堆栈信息
type PanicHandler func(panicValue interface{}, stack []byte)

// =============================================================================
// 日志函数配置 - 与具体日志库解耦
// =============================================================================

// LogFunc 日志函数类型，用于记录错误信息
//
// 参数:
//   - ctx: 上下文，用于传递 trace id 等信息
//   - format: 格式化字符串
//   - args: 格式化参数
type LogFunc func(ctx context.Context, format string, args ...interface{})

// 默认日志函数，使用 fmt.Printf 输出到标准输出
var defaultLogFunc LogFunc = func(ctx context.Context, format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// recoveryLogFunc 当前使用的日志函数
var recoveryLogFunc = defaultLogFunc

// SetRecoveryLogFunc 设置全局日志函数
//
// 该函数用于配置 recovery 工具使用的日志函数，实现与具体日志库的解耦。
// 建议在应用程序初始化时调用此函数。
//
// 参数:
//   - fn: 自定义日志函数
//
// 使用示例:
//
//	// 使用项目的 ilog
//	import "your/project/ilog"
//
//	func init() {
//	    SetRecoveryLogFunc(func(ctx context.Context, format string, args ...interface{}) {
//	        ilog.CtxErrorf(ctx, format, args...)
//	    })
//	}
//
//	// 或者使用标准库 log
//	import "log"
//
//	func init() {
//	    SetRecoveryLogFunc(func(ctx context.Context, format string, args ...interface{}) {
//	        log.Printf(format, args...)
//	    })
//	}
func SetRecoveryLogFunc(fn LogFunc) {
	if fn != nil {
		recoveryLogFunc = fn
	}
}

// ResetRecoveryLogFunc 重置日志函数为默认值
//
// 该函数用于将日志函数恢复为默认的 fmt.Printf 实现。
// 主要用于测试或需要临时切换日志实现的场景。
func ResetRecoveryLogFunc() {
	recoveryLogFunc = defaultLogFunc
}

// =============================================================================
// 核心函数 - 安全启动 Goroutine
// =============================================================================

// Go 安全启动一个 goroutine，自动处理 panic recover
//
// 该函数封装了 goroutine 的启动逻辑，自动添加 defer recover 处理，
// 避免因单个 goroutine panic 导致整个程序崩溃。
//
// 参数:
//   - fn: 要在 goroutine 中执行的函数
//   - options: 配置选项，可以通过位运算组合多个选项（可选参数，默认为打印堆栈+记录日志）
//
// 使用示例:
//
//	// 基本使用（使用默认选项）
//	Go(func() {
//	    // 你的业务逻辑
//	    doSomething()
//	})
//
//	// 使用自定义选项
//	Go(func() {
//	    doSomething()
//	}, RecoveryOptionPrintStack)
//
//	// 组合多个选项
//	Go(func() {
//	    doSomething()
//	}, RecoveryOptionPrintStack | RecoveryOptionLogError)
//
//	// 静默模式
//	Go(func() {
//	    doSomething()
//	}, RecoveryOptionSilent)
//
// 注意事项:
//   - 始终使用此函数启动 goroutine，而不是直接使用 go 关键字
//   - 如果需要自定义 panic 处理逻辑，请使用 GoWithHandler 函数
//   - 默认会打印堆栈并记录 ERROR 级别的日志
func Go(fn func(), options ...RecoveryOption) {
	// 合并所有选项，默认使用推荐配置
	opt := GetDefaultRecoveryOptions()
	if len(options) > 0 {
		opt = RecoveryOptionNone
		for _, o := range options {
			opt |= o
		}
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				handlePanic("[Go]", r, opt, nil, nil)
			}
		}()
		fn()
	}()
}

// GoWithName 安全启动一个带名称标识的 goroutine
//
// 该函数在 Go 的基础上增加了名称参数，方便在日志中识别是哪个 goroutine 发生了 panic。
//
// 参数:
//   - name: goroutine 的名称标识，用于日志输出
//   - fn: 要在 goroutine 中执行的函数
//   - options: 配置选项（可选参数）
//
// 使用示例:
//
//	// 带名称的 goroutine
//	GoWithName("UserDataSync", func() {
//	    syncUserData()
//	})
//
//	// 带名称和自定义选项
//	GoWithName("OrderProcessor", func() {
//	    processOrders()
//	}, RecoveryOptionPrintStack | RecoveryOptionLogError)
//
// 注意事项:
//   - 建议使用有意义的名称，便于问题排查
//   - 名称会出现在日志的前缀中，格式为 [GoWithName-{name}]
func GoWithName(name string, fn func(), options ...RecoveryOption) {
	// 合并所有选项，默认使用推荐配置
	opt := GetDefaultRecoveryOptions()
	if len(options) > 0 {
		opt = RecoveryOptionNone
		for _, o := range options {
			opt |= o
		}
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				prefix := fmt.Sprintf("[GoWithName-%s]", name)
				handlePanic(prefix, r, opt, nil, nil)
			}
		}()
		fn()
	}()
}

// GoWithHandler 安全启动 goroutine，支持自定义 panic 处理器
//
// 该函数允许用户自定义 panic 处理逻辑，适用于需要特殊处理的场景，
// 例如发送告警、上报监控等。
//
// 参数:
//   - fn: 要在 goroutine 中执行的函数
//   - handler: 自定义的 panic 处理函数
//   - options: 配置选项（可选参数）
//
// 使用示例:
//
//	// 自定义 panic 处理
//	GoWithHandler(func() {
//	    doRiskyOperation()
//	}, func(panicValue interface{}, stack []byte) {
//	    // 发送告警
//	    sendAlert(fmt.Sprintf("Panic occurred: %v", panicValue))
//	    // 上报监控
//	    reportMetric("goroutine_panic", 1)
//	})
//
//	// 同时使用自定义处理器和选项
//	GoWithHandler(func() {
//	    doRiskyOperation()
//	}, myPanicHandler, RecoveryOptionPrintStack)
//
// 注意事项:
//   - 自定义处理器会在默认处理逻辑之后执行
//   - 处理器中的异常也会被捕获，不会导致程序崩溃
func GoWithHandler(fn func(), handler PanicHandler, options ...RecoveryOption) {
	// 合并所有选项，默认使用推荐配置
	opt := GetDefaultRecoveryOptions()
	if len(options) > 0 {
		opt = RecoveryOptionNone
		for _, o := range options {
			opt |= o
		}
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				handlePanic("[GoWithHandler]", r, opt, handler, nil)
			}
		}()
		fn()
	}()
}

// GoWithContext 安全启动支持 context 的 goroutine
//
// 该函数支持 context 取消机制，适用于需要优雅关闭的场景。
//
// 参数:
//   - ctx: 上下文，用于控制 goroutine 的生命周期
//   - fn: 要在 goroutine 中执行的函数，接收 context 参数
//   - options: 配置选项（可选参数）
//
// 使用示例:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	GoWithContext(ctx, func(ctx context.Context) {
//	    for {
//	        select {
//	        case <-ctx.Done():
//	            // 收到取消信号，优雅退出
//	            return
//	        default:
//	            // 执行业务逻辑
//	            doWork()
//	        }
//	    }
//	})
//
// 注意事项:
//   - 函数内部应该检查 ctx.Done() 以支持优雅关闭
//   - context 取消时不会触发 panic 处理
func GoWithContext(ctx context.Context, fn func(ctx context.Context), options ...RecoveryOption) {
	// 合并所有选项，默认使用推荐配置
	opt := GetDefaultRecoveryOptions()
	if len(options) > 0 {
		opt = RecoveryOptionNone
		for _, o := range options {
			opt |= o
		}
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				handlePanic("[GoWithContext]", r, opt, nil, ctx)
			}
		}()
		fn(ctx)
	}()
}

// =============================================================================
// 批量并发控制
// =============================================================================

// GoGroup 并发执行多个任务，等待所有任务完成
//
// 该函数封装了 sync.WaitGroup 的使用，简化并发任务的编写。
// 所有任务都会自动处理 panic recover。
//
// 参数:
//   - fns: 要并发执行的函数列表
//   - options: 配置选项（可选参数）
//
// 使用示例:
//
//	// 并发执行多个任务
//	GoGroup([]func(){
//	    func() { task1() },
//	    func() { task2() },
//	    func() { task3() },
//	})
//
//	// 使用自定义选项
//	GoGroup([]func(){
//	    func() { task1() },
//	    func() { task2() },
//	}, RecoveryOptionPrintStack)
//
// 注意事项:
//   - 该函数会阻塞直到所有任务完成
//   - 单个任务 panic 不会影响其他任务的执行
//   - 如果需要收集任务结果，请使用 GoGroupWithResult
func GoGroup(fns []func(), options ...RecoveryOption) {
	// 合并所有选项，默认使用推荐配置
	opt := GetDefaultRecoveryOptions()
	if len(options) > 0 {
		opt = RecoveryOptionNone
		for _, o := range options {
			opt |= o
		}
	}

	var wg sync.WaitGroup
	for i, fn := range fns {
		wg.Add(1)
		go func(index int, f func()) {
			defer func() {
				wg.Done()
				if r := recover(); r != nil {
					prefix := fmt.Sprintf("[GoGroup-task%d]", index)
					handlePanic(prefix, r, opt, nil, nil)
				}
			}()
			f()
		}(i, fn)
	}
	wg.Wait()
}

// GoGroupWithLimit 并发执行多个任务，限制最大并发数
//
// 该函数在 GoGroup 的基础上增加了并发数限制，防止创建过多 goroutine。
//
// 参数:
//   - fns: 要并发执行的函数列表
//   - limit: 最大并发数
//   - options: 配置选项（可选参数）
//
// 使用示例:
//
//	// 限制最多 5 个并发
//	GoGroupWithLimit(tasks, 5)
//
//	// 使用自定义选项
//	GoGroupWithLimit(tasks, 10, RecoveryOptionLogError)
//
// 注意事项:
//   - limit 必须大于 0，否则会 panic
//   - 适用于需要控制资源使用的场景
func GoGroupWithLimit(fns []func(), limit int, options ...RecoveryOption) {
	if limit <= 0 {
		panic("[GoGroupWithLimit] limit must be greater than 0")
	}

	// 合并所有选项，默认使用推荐配置
	opt := GetDefaultRecoveryOptions()
	if len(options) > 0 {
		opt = RecoveryOptionNone
		for _, o := range options {
			opt |= o
		}
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, limit)

	for i, fn := range fns {
		wg.Add(1)
		semaphore <- struct{}{} // 获取信号量

		go func(index int, f func()) {
			defer func() {
				<-semaphore // 释放信号量
				wg.Done()
				if r := recover(); r != nil {
					prefix := fmt.Sprintf("[GoGroupWithLimit-task%d]", index)
					handlePanic(prefix, r, opt, nil, nil)
				}
			}()
			f()
		}(i, fn)
	}
	wg.Wait()
}

// =============================================================================
// 带结果返回的并发任务
// =============================================================================

// TaskResult 任务执行结果
//
// 包含任务执行的结果值、错误信息和 panic 信息。
//
// 字段:
//   - Index: 任务索引，对应输入函数列表中的位置
//   - Value: 任务返回值
//   - Err: 任务返回的错误
//   - Panic: 如果发生 panic，此字段存储 panic 值
//   - Stack: 如果发生 panic，此字段存储堆栈信息
type TaskResult[T any] struct {
	Index int         // 任务索引
	Value T           // 任务返回值
	Err   error       // 任务错误
	Panic interface{} // panic 值（如果有）
	Stack []byte      // panic 堆栈（如果有）
}

// HasPanic 检查任务是否发生了 panic
func (r *TaskResult[T]) HasPanic() bool {
	return r.Panic != nil
}

// IsSuccess 检查任务是否成功（无错误且无 panic）
func (r *TaskResult[T]) IsSuccess() bool {
	return r.Err == nil && r.Panic == nil
}

// TaskFunc 带返回值的任务函数类型
//
// 返回值:
//   - T: 任务结果
//   - error: 任务错误
type TaskFunc[T any] func() (T, error)

// GoGroupWithResult 并发执行多个任务，收集所有任务的结果
//
// 该函数支持泛型，可以收集任意类型的任务结果。
// 所有任务都会自动处理 panic recover，单个任务失败不影响其他任务。
//
// 参数:
//   - fns: 要并发执行的任务函数列表
//   - options: 配置选项（可选参数）
//
// 返回值:
//   - []TaskResult[T]: 所有任务的执行结果列表，顺序与输入函数列表一致
//
// 使用示例:
//
//	// 并发获取多个用户信息
//	results := GoGroupWithResult([]TaskFunc[*User]{
//	    func() (*User, error) { return getUserByID(1) },
//	    func() (*User, error) { return getUserByID(2) },
//	    func() (*User, error) { return getUserByID(3) },
//	})
//
//	// 处理结果
//	for _, r := range results {
//	    if r.HasPanic() {
//	        log.Printf("任务 %d panic: %v", r.Index, r.Panic)
//	        continue
//	    }
//	    if r.Err != nil {
//	        log.Printf("任务 %d 错误: %v", r.Index, r.Err)
//	        continue
//	    }
//	    log.Printf("任务 %d 成功: %+v", r.Index, r.Value)
//	}
//
//	// 并发下载多个文件
//	downloadResults := GoGroupWithResult([]TaskFunc[[]byte]{
//	    func() ([]byte, error) { return downloadFile("url1") },
//	    func() ([]byte, error) { return downloadFile("url2") },
//	})
//
// 注意事项:
//   - 该函数会阻塞直到所有任务完成
//   - 结果列表的顺序与输入函数列表的顺序一致
//   - 单个任务 panic 会被捕获并记录在对应的 TaskResult.Panic 中
//   - 即使某些任务失败，其他任务仍会继续执行
func GoGroupWithResult[T any](fns []TaskFunc[T], options ...RecoveryOption) []TaskResult[T] {
	// 合并所有选项，默认使用推荐配置
	opt := GetDefaultRecoveryOptions()
	if len(options) > 0 {
		opt = RecoveryOptionNone
		for _, o := range options {
			opt |= o
		}
	}

	results := make([]TaskResult[T], len(fns))
	var wg sync.WaitGroup

	for i, fn := range fns {
		wg.Add(1)
		go func(index int, f TaskFunc[T]) {
			defer func() {
				wg.Done()
				if r := recover(); r != nil {
					stack := debug.Stack()
					results[index].Panic = r
					results[index].Stack = stack

					// 使用统一的 panic 处理逻辑记录日志
					prefix := fmt.Sprintf("[GoGroupWithResult-task%d]", index)
					handlePanic(prefix, r, opt, nil, nil)
				}
			}()

			results[index].Index = index
			value, err := f()
			results[index].Value = value
			results[index].Err = err
		}(i, fn)
	}
	wg.Wait()

	return results
}

// GoGroupWithResultAndLimit 并发执行多个任务，限制最大并发数并收集结果
//
// 该函数结合了 GoGroupWithLimit 和 GoGroupWithResult 的功能。
//
// 参数:
//   - fns: 要并发执行的任务函数列表
//   - limit: 最大并发数
//   - options: 配置选项（可选参数）
//
// 返回值:
//   - []TaskResult[T]: 所有任务的执行结果列表
//
// 使用示例:
//
//	// 限制并发数为 5，批量处理订单
//	results := GoGroupWithResultAndLimit(orderTasks, 5)
//
//	// 统计成功和失败数量
//	var success, failed int
//	for _, r := range results {
//	    if r.IsSuccess() {
//	        success++
//	    } else {
//	        failed++
//	    }
//	}
//
// 注意事项:
//   - limit 必须大于 0，否则会 panic
//   - 适用于大量任务需要限流执行的场景
func GoGroupWithResultAndLimit[T any](fns []TaskFunc[T], limit int, options ...RecoveryOption) []TaskResult[T] {
	if limit <= 0 {
		panic("[GoGroupWithResultAndLimit] limit must be greater than 0")
	}

	// 合并所有选项，默认使用推荐配置
	opt := GetDefaultRecoveryOptions()
	if len(options) > 0 {
		opt = RecoveryOptionNone
		for _, o := range options {
			opt |= o
		}
	}

	results := make([]TaskResult[T], len(fns))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, limit)

	for i, fn := range fns {
		wg.Add(1)
		semaphore <- struct{}{} // 获取信号量

		go func(index int, f TaskFunc[T]) {
			defer func() {
				<-semaphore // 释放信号量
				wg.Done()
				if r := recover(); r != nil {
					stack := debug.Stack()
					results[index].Panic = r
					results[index].Stack = stack

					prefix := fmt.Sprintf("[GoGroupWithResultAndLimit-task%d]", index)
					handlePanic(prefix, r, opt, nil, nil)
				}
			}()

			results[index].Index = index
			value, err := f()
			results[index].Value = value
			results[index].Err = err
		}(i, fn)
	}
	wg.Wait()

	return results
}

// =============================================================================
// 定时任务支持
// =============================================================================

// GoTicker 安全启动一个定时执行的 goroutine
//
// 该函数封装了 time.Ticker 的使用，自动处理 panic recover。
// 返回一个 stop 函数，调用后可以停止定时任务。
//
// 参数:
//   - interval: 执行间隔
//   - fn: 要定时执行的函数
//   - options: 配置选项（可选参数）
//
// 返回值:
//   - stop: 停止函数，调用后停止定时任务
//
// 使用示例:
//
//	// 每 5 秒执行一次
//	stop := GoTicker(5*time.Second, func() {
//	    checkHealth()
//	})
//
//	// 10 分钟后停止
//	time.AfterFunc(10*time.Minute, stop)
//
//	// 或者手动停止
//	// stop()
//
// 注意事项:
//   - 定时任务会立即执行一次，然后按间隔重复执行
//   - 单次执行 panic 不会停止后续执行
//   - 使用 stop() 可以优雅地停止定时任务
func GoTicker(interval time.Duration, fn func(), options ...RecoveryOption) (stop func()) {
	// 合并所有选项，默认使用推荐配置
	opt := GetDefaultRecoveryOptions()
	if len(options) > 0 {
		opt = RecoveryOptionNone
		for _, o := range options {
			opt |= o
		}
	}

	ticker := time.NewTicker(interval)
	stopChan := make(chan struct{})

	go func() {
		// 立即执行一次
		safeExecute("[GoTicker]", fn, opt)

		for {
			select {
			case <-ticker.C:
				safeExecute("[GoTicker]", fn, opt)
			case <-stopChan:
				ticker.Stop()
				return
			}
		}
	}()

	return func() {
		close(stopChan)
	}
}

// GoAfter 安全启动一个延迟执行的 goroutine
//
// 该函数封装了 time.AfterFunc 的使用，自动处理 panic recover。
//
// 参数:
//   - delay: 延迟时间
//   - fn: 要执行的函数
//   - options: 配置选项（可选参数）
//
// 返回值:
//   - timer: time.Timer 对象，可用于取消延迟任务
//
// 使用示例:
//
//	// 5 秒后执行
//	timer := GoAfter(5*time.Second, func() {
//	    sendNotification()
//	})
//
//	// 如果需要取消
//	// timer.Stop()
//
// 注意事项:
//   - 如果在延迟期间需要取消，调用 timer.Stop()
//   - 函数只会执行一次
func GoAfter(delay time.Duration, fn func(), options ...RecoveryOption) *time.Timer {
	// 合并所有选项，默认使用推荐配置
	opt := GetDefaultRecoveryOptions()
	if len(options) > 0 {
		opt = RecoveryOptionNone
		for _, o := range options {
			opt |= o
		}
	}

	return time.AfterFunc(delay, func() {
		safeExecute("[GoAfter]", fn, opt)
	})
}

// =============================================================================
// 内部辅助函数
// =============================================================================

// handlePanic 处理 panic 的统一逻辑
//
// 参数:
//   - prefix: 日志前缀，用于标识来源
//   - panicValue: panic 时的值
//   - opt: 配置选项
//   - handler: 自定义处理器（可为 nil）
//   - ctx: 上下文（可为 nil，nil 时使用 context.Background()）
func handlePanic(prefix string, panicValue interface{}, opt RecoveryOption, handler PanicHandler, ctx context.Context) {
	// 静默模式直接返回
	if opt.Has(RecoveryOptionSilent) {
		// 但仍然执行自定义处理器
		if handler != nil {
			safeCallHandler(handler, panicValue, nil)
		}
		return
	}

	// 确保 ctx 不为 nil
	if ctx == nil {
		ctx = context.Background()
	}

	// 获取堆栈信息
	var stack []byte
	if opt.Has(RecoveryOptionPrintStack) || handler != nil {
		stack = debug.Stack()
	}

	// 记录错误日志
	if opt.Has(RecoveryOptionLogError) {
		if opt.Has(RecoveryOptionPrintStack) {
			recoveryLogFunc(ctx, "%s panic recovered: %v\nstack:\n%s", prefix, panicValue, string(stack))
		} else {
			recoveryLogFunc(ctx, "%s panic recovered: %v", prefix, panicValue)
		}
	} else if opt.Has(RecoveryOptionPrintStack) {
		// 仅打印堆栈，不记录日志
		fmt.Printf("%s panic recovered: %v\nstack:\n%s\n", prefix, panicValue, string(stack))
	}

	// 执行自定义处理器
	if handler != nil {
		safeCallHandler(handler, panicValue, stack)
	}
}

// safeExecute 安全执行函数，捕获 panic
//
// 参数:
//   - prefix: 日志前缀
//   - fn: 要执行的函数
//   - opt: 配置选项
func safeExecute(prefix string, fn func(), opt RecoveryOption) {
	defer func() {
		if r := recover(); r != nil {
			handlePanic(prefix, r, opt, nil, nil)
		}
	}()
	fn()
}

// safeCallHandler 安全调用自定义处理器
//
// 参数:
//   - handler: 自定义处理器
//   - panicValue: panic 时的值
//   - stack: 堆栈信息
func safeCallHandler(handler PanicHandler, panicValue interface{}, stack []byte) {
	defer func() {
		if r := recover(); r != nil {
			recoveryLogFunc(context.Background(), "[safeCallHandler] panic handler itself panicked: %v", r)
		}
	}()
	handler(panicValue, stack)
}

// =============================================================================
// SafeGo Builder - 链式调用风格（高级用法）
// =============================================================================

// SafeGoBuilder 安全 Goroutine 构建器，支持链式调用
//
// 使用示例:
//
//	NewSafeGo(func() {
//	    doWork()
//	}).
//	    WithName("MyWorker").
//	    WithOptions(RecoveryOptionPrintStack).
//	    WithHandler(func(p interface{}, s []byte) {
//	        sendAlert(fmt.Sprintf("Worker panic: %v", p))
//	    }).
//	    Run()
type SafeGoBuilder struct {
	fn      func()                   // 要执行的函数
	name    string                   // goroutine 名称
	options RecoveryOption           // 配置选项
	handler PanicHandler             // 自定义处理器
	ctx     context.Context          // 上下文
	ctxFn   func(ctx context.Context) // 带上下文的函数
}

// NewSafeGo 创建新的 SafeGoBuilder 实例
//
// 参数:
//   - fn: 要在 goroutine 中执行的函数
//
// 返回值:
//   - *SafeGoBuilder: 构建器实例，支持链式调用
//
// 使用示例:
//
//	builder := NewSafeGo(func() {
//	    doWork()
//	})
func NewSafeGo(fn func()) *SafeGoBuilder {
	return &SafeGoBuilder{
		fn:      fn,
		options: GetDefaultRecoveryOptions(),
	}
}

// NewSafeGoWithContext 创建支持 context 的 SafeGoBuilder 实例
//
// 参数:
//   - ctx: 上下文
//   - fn: 带上下文的函数
//
// 返回值:
//   - *SafeGoBuilder: 构建器实例
//
// 使用示例:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//
//	NewSafeGoWithContext(ctx, func(ctx context.Context) {
//	    doWorkWithContext(ctx)
//	}).Run()
func NewSafeGoWithContext(ctx context.Context, fn func(ctx context.Context)) *SafeGoBuilder {
	return &SafeGoBuilder{
		ctx:     ctx,
		ctxFn:   fn,
		options: GetDefaultRecoveryOptions(),
	}
}

// WithName 设置 goroutine 名称
//
// 参数:
//   - name: goroutine 名称，用于日志标识
//
// 返回值:
//   - *SafeGoBuilder: 返回自身以支持链式调用
func (b *SafeGoBuilder) WithName(name string) *SafeGoBuilder {
	b.name = name
	return b
}

// WithOptions 设置配置选项
//
// 参数:
//   - options: 配置选项，支持位运算组合
//
// 返回值:
//   - *SafeGoBuilder: 返回自身以支持链式调用
func (b *SafeGoBuilder) WithOptions(options RecoveryOption) *SafeGoBuilder {
	b.options = options
	return b
}

// WithHandler 设置自定义 panic 处理器
//
// 参数:
//   - handler: 自定义处理函数
//
// 返回值:
//   - *SafeGoBuilder: 返回自身以支持链式调用
func (b *SafeGoBuilder) WithHandler(handler PanicHandler) *SafeGoBuilder {
	b.handler = handler
	return b
}

// Run 启动 goroutine
//
// 使用示例:
//
//	NewSafeGo(myFunc).
//	    WithName("Worker").
//	    Run()
func (b *SafeGoBuilder) Run() {
	prefix := "[SafeGo]"
	if b.name != "" {
		prefix = fmt.Sprintf("[SafeGo-%s]", b.name)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				handlePanic(prefix, r, b.options, b.handler, b.ctx)
			}
		}()

		if b.ctxFn != nil && b.ctx != nil {
			b.ctxFn(b.ctx)
		} else if b.fn != nil {
			b.fn()
		}
	}()
}