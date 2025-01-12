package otel

import (
	"fmt"
	"regexp"
	"runtime"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func GetFunctionCallAttributes(skipCaller int) (string, []attribute.KeyValue) {
	pc, file, no, _ := runtime.Caller(skipCaller + 1)
	details := runtime.FuncForPC(pc)
	re := regexp.MustCompile(`^(.+)\.\(\*([^)]+)\)\.(\w+)$`)
	functionParts := re.FindStringSubmatch(details.Name())

	attributes := []attribute.KeyValue{
		semconv.CodeFilepath(file),
		semconv.CodeLineNumber(no),
		semconv.CodeNamespace(functionParts[1]),
		semconv.CodeFunction(functionParts[3]),
	}

	return fmt.Sprintf("%s.%s", functionParts[2], functionParts[3]), attributes
}
