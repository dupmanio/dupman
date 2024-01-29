package fx

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type IRoute interface {
	Register(engine *gin.Engine)
}

const routeGroupTag = `group:"routes"`

func AsRoute(fc any, annotations ...fx.Annotation) any {
	annotations = append(
		annotations,
		fx.As(new(IRoute)),
		fx.ResultTags(routeGroupTag),
	)

	return fx.Annotate(fc, annotations...)
}

func AsRouteReceiver(fc any, annotations ...fx.Annotation) any {
	annotations = append(
		annotations,
		fx.ParamTags(routeGroupTag),
	)

	return fx.Annotate(fc, annotations...)
}

func RegisterRoutes(engine *gin.Engine, routes ...IRoute) {
	for _, route := range routes {
		route.Register(engine)
	}
}
