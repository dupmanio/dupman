package wrapper

import (
	"strings"

	"github.com/dupmanio/dupman/packages/common/otel"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type FxWrapper struct {
	logger *zap.Logger
}

func NewFxWrapper(logger *zap.Logger) *FxWrapper {
	logger = logger.With(
		zap.String(string(otel.ComponentKey), "fx"),
	)

	return &FxWrapper{
		logger: logger,
	}
}

func (wrap *FxWrapper) LogEvent(event fxevent.Event) {
	switch ev := event.(type) {
	case *fxevent.OnStartExecuting:
		wrap.logger.Debug("OnStart hook executing",
			zap.String("callee", ev.FunctionName),
			zap.String("caller", ev.CallerName),
		)
	case *fxevent.OnStartExecuted:
		if ev.Err != nil {
			wrap.logger.Error("OnStart hook failed",
				zap.String("callee", ev.FunctionName),
				zap.String("caller", ev.CallerName),
				zap.Error(ev.Err),
			)
		} else {
			wrap.logger.Debug("OnStart hook executed",
				zap.String("callee", ev.FunctionName),
				zap.String("caller", ev.CallerName),
				zap.String("runtime", ev.Runtime.String()),
			)
		}
	case *fxevent.OnStopExecuting:
		wrap.logger.Debug("OnStop hook executing",
			zap.String("callee", ev.FunctionName),
			zap.String("caller", ev.CallerName),
		)
	case *fxevent.OnStopExecuted:
		if ev.Err != nil {
			wrap.logger.Error("OnStop hook failed",
				zap.String("callee", ev.FunctionName),
				zap.String("caller", ev.CallerName),
				zap.Error(ev.Err),
			)
		} else {
			wrap.logger.Debug("OnStop hook executed",
				zap.String("callee", ev.FunctionName),
				zap.String("caller", ev.CallerName),
				zap.String("runtime", ev.Runtime.String()),
			)
		}
	case *fxevent.Supplied:
		if ev.Err != nil {
			wrap.logger.Error("error encountered while applying options",
				zap.String("type", ev.TypeName),
				zap.Strings("stacktrace", ev.StackTrace),
				moduleField(ev.ModuleName),
				zap.Error(ev.Err))
		} else {
			wrap.logger.Debug("supplied",
				zap.String("type", ev.TypeName),
				zap.Strings("stacktrace", ev.StackTrace),
				moduleField(ev.ModuleName),
			)
		}
	case *fxevent.Provided:
		for _, rtype := range ev.OutputTypeNames {
			wrap.logger.Debug("provided",
				zap.String("constructor", ev.ConstructorName),
				zap.Strings("stacktrace", ev.StackTrace),
				moduleField(ev.ModuleName),
				zap.String("type", rtype),
				maybeBool("private", ev.Private),
			)
		}

		if ev.Err != nil {
			wrap.logger.Error("error encountered while applying options",
				moduleField(ev.ModuleName),
				zap.Strings("stacktrace", ev.StackTrace),
				zap.Error(ev.Err))
		}
	case *fxevent.Replaced:
		for _, rtype := range ev.OutputTypeNames {
			wrap.logger.Debug("replaced",
				zap.Strings("stacktrace", ev.StackTrace),
				moduleField(ev.ModuleName),
				zap.String("type", rtype),
			)
		}

		if ev.Err != nil {
			wrap.logger.Error("error encountered while replacing",
				zap.Strings("stacktrace", ev.StackTrace),
				moduleField(ev.ModuleName),
				zap.Error(ev.Err))
		}
	case *fxevent.Decorated:
		for _, rtype := range ev.OutputTypeNames {
			wrap.logger.Debug("decorated",
				zap.String("decorator", ev.DecoratorName),
				zap.Strings("stacktrace", ev.StackTrace),
				moduleField(ev.ModuleName),
				zap.String("type", rtype),
			)
		}

		if ev.Err != nil {
			wrap.logger.Error("error encountered while applying options",
				zap.Strings("stacktrace", ev.StackTrace),
				moduleField(ev.ModuleName),
				zap.Error(ev.Err))
		}
	case *fxevent.Run:
		if ev.Err != nil {
			wrap.logger.Error("error returned",
				zap.String("name", ev.Name),
				zap.String("kind", ev.Kind),
				moduleField(ev.ModuleName),
				zap.Error(ev.Err),
			)
		} else {
			wrap.logger.Debug("run",
				zap.String("name", ev.Name),
				zap.String("kind", ev.Kind),
				moduleField(ev.ModuleName),
			)
		}
	case *fxevent.Invoking:
		// Do not log stack as it will make logs hard to read.
		wrap.logger.Debug("invoking",
			zap.String("function", ev.FunctionName),
			moduleField(ev.ModuleName),
		)
	case *fxevent.Invoked:
		if ev.Err != nil {
			wrap.logger.Error("invoke failed",
				zap.Error(ev.Err),
				zap.String("stack", ev.Trace),
				zap.String("function", ev.FunctionName),
				moduleField(ev.ModuleName),
			)
		}
	case *fxevent.Stopping:
		wrap.logger.Debug("received signal",
			zap.String("signal", strings.ToUpper(ev.Signal.String())))
	case *fxevent.Stopped:
		if ev.Err != nil {
			wrap.logger.Error("stop failed", zap.Error(ev.Err))
		}
	case *fxevent.RollingBack:
		wrap.logger.Error("start failed, rolling back", zap.Error(ev.StartErr))
	case *fxevent.RolledBack:
		if ev.Err != nil {
			wrap.logger.Error("rollback failed", zap.Error(ev.Err))
		}
	case *fxevent.Started:
		if ev.Err != nil {
			wrap.logger.Error("start failed", zap.Error(ev.Err))
		} else {
			wrap.logger.Debug("started")
		}
	case *fxevent.LoggerInitialized:
		if ev.Err != nil {
			wrap.logger.Error("custom logger initialization failed", zap.Error(ev.Err))
		} else {
			wrap.logger.Debug("initialized custom fxevent.Logger", zap.String("function", ev.ConstructorName))
		}
	}
}

func moduleField(name string) zap.Field {
	if len(name) == 0 {
		return zap.Skip()
	}

	return zap.String("module", name)
}

func maybeBool(name string, b bool) zap.Field {
	if b {
		return zap.Bool(name, true)
	}

	return zap.Skip()
}
