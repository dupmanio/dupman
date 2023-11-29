package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/api/service"
	"github.com/dupmanio/dupman/packages/common/otel"
	commonServices "github.com/dupmanio/dupman/packages/common/service"
	"github.com/dupmanio/dupman/packages/domain/dto"
	domainErrors "github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type UserController struct {
	httpSvc *commonServices.HTTPService
	userSvc *service.UserService
	ot      *otel.OTel
}

func NewUserController(
	httpSvc *commonServices.HTTPService,
	userSvc *service.UserService,
	ot *otel.OTel,
) (*UserController, error) {
	return &UserController{
		httpSvc: httpSvc,
		userSvc: userSvc,
		ot:      ot,
	}, nil
}

func (ctrl *UserController) Create(ctx *gin.Context) {
	var (
		entity = &model.User{}

		payload  *dto.UserOnCreate
		response dto.UserAccount
	)

	ctrl.httpSvc.EnrichSpanWithControllerAttributes(ctx)

	if err := ctx.ShouldBind(&payload); err != nil {
		ctrl.httpSvc.HTTPValidationErrorWithOTelLog(ctx, err)

		return
	}

	_ = copier.Copy(&entity, &payload)

	user, err := ctrl.userSvc.Create(ctx, entity)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, domainErrors.ErrUserAlreadyExists) {
			statusCode = http.StatusBadRequest
		}

		ctrl.httpSvc.HTTPErrorWithOTelLog(ctx, "Unable to create User", statusCode, err)

		return
	}

	_ = copier.Copy(&response, &user)

	ctrl.ot.LogInfoEvent(ctx, "User has been created successfully", otel.UserID(user.ID))
	ctrl.httpSvc.HTTPResponse(ctx, http.StatusCreated, response)
}

func (ctrl *UserController) Update(ctx *gin.Context) {
	var (
		entity = &model.User{}

		payload  *dto.UserOnUpdate
		response dto.UserAccount
	)

	ctrl.httpSvc.EnrichSpanWithControllerAttributes(ctx)

	if err := ctx.ShouldBind(&payload); err != nil {
		ctrl.httpSvc.HTTPValidationErrorWithOTelLog(ctx, err)

		return
	}

	_ = copier.Copy(&entity, &payload)

	user, err := ctrl.userSvc.Update(ctx, entity)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, domainErrors.ErrUserDoesNotExist) {
			statusCode = http.StatusNotFound
		}

		ctrl.httpSvc.HTTPErrorWithOTelLog(ctx, "Unable to update User", statusCode, err, otel.UserID(payload.ID))

		return
	}

	_ = copier.Copy(&response, &user)

	ctrl.ot.LogInfoEvent(ctx, "User has been updated successfully", otel.UserID(user.ID))
	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, response)
}

func (ctrl *UserController) GetContactInfo(ctx *gin.Context) {
	var response dto.ContactInfo

	ctrl.httpSvc.EnrichSpanWithControllerAttributes(ctx)

	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctrl.httpSvc.HTTPErrorWithOTelLog(
			ctx,
			"Invalid User ID",
			http.StatusBadRequest,
			fmt.Errorf("invalid user ID: %w", err),
		)

		return
	}

	user, err := ctrl.userSvc.GetSingle(ctx, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, domainErrors.ErrUserDoesNotExist) {
			statusCode = http.StatusNotFound
		}

		ctrl.httpSvc.HTTPErrorWithOTelLog(ctx, "Unable to load User", statusCode, err, otel.UserID(userID))

		return
	}

	_ = copier.Copy(&response, &user)

	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, response)
}
