package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dupmanio/dupman/packages/common/otel"
	"github.com/dupmanio/dupman/packages/common/pagination"
	commonServices "github.com/dupmanio/dupman/packages/common/service"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/notify/model"
	"github.com/dupmanio/dupman/packages/notify/server"
	"github.com/dupmanio/dupman/packages/notify/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type NotificationController struct {
	server          *server.Server
	httpSvc         *commonServices.HTTPService
	notificationSvc *service.NotificationService
	ot              *otel.OTel
}

func NewNotificationController(
	server *server.Server,
	httpSvc *commonServices.HTTPService,
	notificationSvc *service.NotificationService,
	ot *otel.OTel,
) (*NotificationController, error) {
	return &NotificationController{
		server:          server,
		httpSvc:         httpSvc,
		notificationSvc: notificationSvc,
		ot:              ot,
	}, nil
}

func (ctrl *NotificationController) Create(ctx *gin.Context) {
	var (
		entity = &model.Notification{}

		payload  *dto.NotificationOnCreate
		response dto.NotificationOnResponse
	)

	ctrl.httpSvc.EnrichSpanWithControllerAttributes(ctx)

	if err := ctx.ShouldBind(&payload); err != nil {
		ctrl.httpSvc.HTTPValidationErrorWithOTelLog(ctx, err)

		return
	}

	_ = copier.Copy(&entity, &payload)

	notification, err := ctrl.notificationSvc.Create(ctx, entity)
	if err != nil {
		ctrl.httpSvc.HTTPErrorWithOTelLog(ctx, "Unable to create Notification", http.StatusInternalServerError, err)

		return
	}

	_ = copier.Copy(&response, &notification)

	// @todo: move to notification service.
	if err = ctrl.notificationSvc.SendNotificationToChannel(ctx, response); err != nil {
		ctrl.httpSvc.HTTPErrorWithOTelLog(ctx, "Unable to send Notification to channel", http.StatusInternalServerError, err)

		return
	}

	ctrl.ot.LogInfoEvent(ctx, "Notification has been created successfully", otel.NotificationID(notification.ID))
	ctrl.httpSvc.HTTPResponse(ctx, http.StatusCreated, response)
}

func (ctrl *NotificationController) GetAll(ctx *gin.Context) {
	var response dto.NotificationsOnResponse

	ctrl.httpSvc.EnrichSpanWithControllerAttributes(ctx)
	pager := pagination.Paginate(ctx)

	notifications, err := ctrl.notificationSvc.GetAllForCurrentUser(ctx, pager)
	if err != nil {
		ctrl.httpSvc.HTTPErrorWithOTelLog(ctx, "Unable to get notifications", http.StatusInternalServerError, err)

		return
	}

	_ = copier.Copy(&response, &notifications)

	ctrl.httpSvc.HTTPPaginatedResponse(ctx, http.StatusOK, response, pager)
}

func (ctrl *NotificationController) GetCount(ctx *gin.Context) {
	ctrl.httpSvc.EnrichSpanWithControllerAttributes(ctx)

	count, err := ctrl.notificationSvc.GetCountForCurrentUser(ctx)
	if err != nil {
		ctrl.httpSvc.HTTPErrorWithOTelLog(ctx, "Unable to get notifications count", http.StatusInternalServerError, err)

		return
	}

	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, count)
}

func (ctrl *NotificationController) Realtime(ctx *gin.Context) {
	const heartbeatInterval = 10 * time.Second

	ctrl.httpSvc.EnrichSpanWithControllerAttributes(ctx)

	pubSub, err := ctrl.notificationSvc.SubscribeToUserNotifications(ctx)
	if err != nil {
		ctrl.httpSvc.HTTPErrorWithOTelLog(
			ctx,
			"Unable to subscribe to realtime Notifications",
			http.StatusInternalServerError,
			err,
		)

		return
	}

	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Transfer-Encoding", "chunked")

	for {
		select {
		case message := <-pubSub.Channel():
			ctx.SSEvent("notification", message.Payload)
			ctx.Writer.Flush()
		case now := <-time.After(heartbeatInterval):
			ctx.SSEvent("heartbeat", now)
			ctx.Writer.Flush()
		// @todo: fix graceful shutdown.
		case sig := <-ctrl.server.Interrupt:
			ctx.SSEvent("close", sig)
			ctx.Writer.Flush()

			return
		}
	}
}

func (ctrl *NotificationController) MarkAsRead(ctx *gin.Context) {
	ctrl.httpSvc.EnrichSpanWithControllerAttributes(ctx)

	notificationID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctrl.httpSvc.HTTPErrorWithOTelLog(
			ctx,
			"Invalid Notification ID",
			http.StatusBadRequest,
			fmt.Errorf("invalid Notification ID: %w", err),
		)

		return
	}

	// @todo: implement 404

	if err := ctrl.notificationSvc.MarkAsRead(ctx, notificationID); err != nil {
		ctrl.httpSvc.HTTPErrorWithOTelLog(
			ctx,
			"Unable to mark Notification as read",
			http.StatusInternalServerError,
			err,
			otel.NotificationID(notificationID),
		)

		return
	}

	ctrl.ot.LogInfoEvent(ctx, "Notification has been marked as read successfully", otel.NotificationID(notificationID))
	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, gin.H{})
}

func (ctrl *NotificationController) MarkAllAsRead(ctx *gin.Context) {
	ctrl.httpSvc.EnrichSpanWithControllerAttributes(ctx)

	if err := ctrl.notificationSvc.MarkAllAsReadForCurrentUser(ctx); err != nil {
		ctrl.httpSvc.HTTPErrorWithOTelLog(ctx, "Unable to mark Notifications as read", http.StatusInternalServerError, err)

		return
	}

	ctrl.ot.LogInfoEvent(ctx, "Notifications has been marked as read successfully")
	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, gin.H{})
}

func (ctrl *NotificationController) DeleteAll(ctx *gin.Context) {
	ctrl.httpSvc.EnrichSpanWithControllerAttributes(ctx)

	if err := ctrl.notificationSvc.DeleteAllForCurrentUser(ctx); err != nil {
		ctrl.httpSvc.HTTPErrorWithOTelLog(ctx, "Unable to delete Notifications", http.StatusInternalServerError, err)

		return
	}

	ctrl.ot.LogInfoEvent(ctx, "Notifications has been deleted successfully")
	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, gin.H{})
}
