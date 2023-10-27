package controller

import (
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/common/pagination"
	commonServices "github.com/dupmanio/dupman/packages/common/service"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/notify/model"
	"github.com/dupmanio/dupman/packages/notify/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type NotificationController struct {
	httpSvc         *commonServices.HTTPService
	notificationSvc *service.NotificationService
}

func NewNotificationController(
	httpSvc *commonServices.HTTPService,
	notificationSvc *service.NotificationService,
) (*NotificationController, error) {
	return &NotificationController{
		httpSvc:         httpSvc,
		notificationSvc: notificationSvc,
	}, nil
}

func (ctrl *NotificationController) Create(ctx *gin.Context) {
	var (
		entity = &model.Notification{}

		payload  *dto.NotificationOnCreate
		response dto.NotificationOnResponse
	)

	if err := ctx.ShouldBind(&payload); err != nil {
		ctrl.httpSvc.HTTPValidationError(ctx, err)

		return
	}

	_ = copier.Copy(&entity, &payload)

	website, err := ctrl.notificationSvc.Create(entity)
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	_ = copier.Copy(&response, &website)

	ctrl.httpSvc.HTTPResponse(ctx, http.StatusCreated, response)
}

func (ctrl *NotificationController) GetAll(ctx *gin.Context) {
	var response dto.NotificationsOnResponse

	pager := pagination.Paginate(ctx)

	notifications, err := ctrl.notificationSvc.GetAllForCurrentUser(ctx, pager)
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	_ = copier.Copy(&response, &notifications)

	ctrl.httpSvc.HTTPPaginatedResponse(ctx, http.StatusOK, response, pager)
}

func (ctrl *NotificationController) GetCount(ctx *gin.Context) {
	count, err := ctrl.notificationSvc.GetCountForCurrentUser(ctx)
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, count)
}

func (ctrl *NotificationController) MarkAsRead(ctx *gin.Context) {
	notificationID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusBadRequest, fmt.Sprintf("invalid notification ID: %s", err))

		return
	}

	// @todo: implement 404

	if err := ctrl.notificationSvc.MarkAsRead(ctx, notificationID); err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, gin.H{})
}

func (ctrl *NotificationController) MarkAllAsRead(ctx *gin.Context) {
	if err := ctrl.notificationSvc.MarkAllAsReadForCurrentUser(ctx); err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, gin.H{})
}

func (ctrl *NotificationController) DeleteAll(ctx *gin.Context) {
	if err := ctrl.notificationSvc.DeleteAllForCurrentUser(ctx); err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, gin.H{})
}
