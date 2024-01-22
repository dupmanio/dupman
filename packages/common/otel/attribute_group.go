package otel

import (
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

const AttributeKeyPrefix = "dupman."

const (
	ComponentKey = attribute.Key(AttributeKeyPrefix + "component")

	TraceIDKey = attribute.Key(AttributeKeyPrefix + "trace.id")
	SpanIDKey  = attribute.Key(AttributeKeyPrefix + "trace.span.id")

	MessageTypeKey = attribute.Key(AttributeKeyPrefix + "message.type")

	PaginationPageKey  = attribute.Key(AttributeKeyPrefix + "pagination.page")
	PaginationLimitKey = attribute.Key(AttributeKeyPrefix + "pagination.limit")

	GinHandlerKey    = attribute.Key(AttributeKeyPrefix + "gin.handler")
	GinHandlerNumKey = attribute.Key(AttributeKeyPrefix + "gin.handler.num")

	WebsiteIDKey = attribute.Key(AttributeKeyPrefix + "website.id")

	NotificationIDKey = attribute.Key(AttributeKeyPrefix + "notification.id")

	UserIDKey = attribute.Key(AttributeKeyPrefix + "user.id")

	RouteKey    = attribute.Key(AttributeKeyPrefix + "route")
	MigratorKey = attribute.Key(AttributeKeyPrefix + "migrator")

	DurationKey  = attribute.Key(AttributeKeyPrefix + "duration")
	RenewedAtKey = attribute.Key(AttributeKeyPrefix + "renewed_at")

	VaultPathKey = attribute.Key(AttributeKeyPrefix + "vault.path")
)

func TraceID(val string) attribute.KeyValue {
	return TraceIDKey.String(val)
}

func SpanID(val string) attribute.KeyValue {
	return SpanIDKey.String(val)
}

func MessageType(val string) attribute.KeyValue {
	return MessageTypeKey.String(val)
}

func PaginationPage(val int) attribute.KeyValue {
	return PaginationPageKey.Int(val)
}

func PaginationLimit(val int) attribute.KeyValue {
	return PaginationLimitKey.Int(val)
}

func WebsiteID(val uuid.UUID) attribute.KeyValue {
	return WebsiteIDKey.String(val.String())
}

func NotificationID(val uuid.UUID) attribute.KeyValue {
	return NotificationIDKey.String(val.String())
}

func UserID(val uuid.UUID) attribute.KeyValue {
	return UserIDKey.String(val.String())
}

func VaultPath(val string) attribute.KeyValue {
	return VaultPathKey.String(val)
}
