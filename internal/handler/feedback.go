package handler

import (
	"context"
	"log/slog"
	"time"

	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/models"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/service"
	"github.com/danielgtaylor/huma/v2"
)

type FeedbackHandler struct {
	service service.FeedbackService
}

func NewFeedbackHandler(fbService *service.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{service: *fbService}
}

// HTTP handler
func (fbH *FeedbackHandler) PostFeedback(ctx context.Context, input *models.FeedbackRequest) (*models.RespWithBody, error) {
	feedback := models.FeedbackModel{
		Name:    input.Body.Name,
		Email:   input.Body.Email,
		Phone:   input.Body.Phone,
		Message: input.Body.Message,
	}

	_, err := fbH.service.Save(feedback)
	if err != nil {
		slog.Error("Failed to save feedback: " + err.Error())
		return nil, huma.Error500InternalServerError("Internal server error")
	}

	return &models.RespWithBody{
		Body: models.DefaultResponse{
			Status:    "success",
			Message:   "Feedback saved successfully",
			Timestamp: time.Now(),
		},
	}, nil
}
