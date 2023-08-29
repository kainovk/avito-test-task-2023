package users

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"avito-test-task-2023/internal/lib/api/response"
	"avito-test-task-2023/internal/lib/logger/sl"
)

type ConfigureSegmentsRequest struct {
	SegmentsToAdd    []SegmentRequest `json:"segments_to_add"`
	SegmentsToDelete []string         `json:"segments_to_delete"`
}

type SegmentRequest struct {
	Slug     string     `json:"slug" validate:"required"`
	DeleteAt *time.Time `json:"delete_at"`
}

type ConfigureSegmentsResponse struct {
	response.Response
}

type UserSegmentConfigurer interface {
	ConfigureUserSegments(userID int64, segAdd []SegmentRequest, segDel []string) error
}

// NewUserSegmentConfigurer handles the HTTP request for configuring user segments.
//
// @Summary Configure user segments
// @Description Configure user segments by adding and/or deleting segments for a user.
// @Tags users
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param request body ConfigureSegmentsRequest true "Request body"
// @Success 200 {object} ConfigureSegmentsResponse
// @Failure 400 {object} ConfigureSegmentsResponse
// @Failure 500 {object} ConfigureSegmentsResponse
// @Router /users/{user_id}/configure-segments [post]
func NewUserSegmentConfigurer(log *slog.Logger, userSegmentConfigurer UserSegmentConfigurer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.configure-segments.NewUserSegmentConfigurer"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req ConfigureSegmentsRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("empty request"))
			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request body"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		userIDStr := chi.URLParam(r, "user_id")
		if userIDStr == "" {
			log.Info("user_id param is empty")

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			log.Error("failed to parse user_id")

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request"))
		}

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		err = userSegmentConfigurer.ConfigureUserSegments(int64(userID), req.SegmentsToAdd, req.SegmentsToDelete)
		if err != nil {
			log.Error("failed to configure user segments", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to configure user segments"))
			return
		}

		log.Info("user segments updated")

		render.JSON(w, r, ConfigureSegmentsResponse{
			Response: response.OK(),
		})
	}
}
