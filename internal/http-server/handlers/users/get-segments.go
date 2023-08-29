package users

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"avito-test-task-2023/internal/lib/api/response"
	"avito-test-task-2023/internal/lib/logger/sl"
	"avito-test-task-2023/internal/models/segment"
)

type GetSegmentsResponse struct {
	Segments []string `json:"segments"`
}

type GetSegmentsResponseFailed struct {
	response.Response
}

type UserSegmentsGetter interface {
	GetUserSegments(userID int64) ([]*segment.Segment, error)
}

// NewUserSegmentsGetter handles the HTTP request for retrieving segments of a user.
//
// @Summary Get user segments
// @Description Retrieve segments associated with a user by user ID.
// @Tags users
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} GetSegmentsResponse
// @Failure 400 {object} GetSegmentsResponseFailed
// @Failure 500 {object} GetSegmentsResponseFailed
// @Router /users/{user_id}/segments [get]
func NewUserSegmentsGetter(log *slog.Logger, userSegmentsGetter UserSegmentsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.get-segments.NewUserSegmentsGetter"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

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

		segments, err := userSegmentsGetter.GetUserSegments(int64(userID))
		if err != nil {
			log.Error("failed to get user segments", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get user segments"))
			return
		}

		log.Info("user segments retrieved")

		segmentSlugs := make([]string, len(segments))
		for i, seg := range segments {
			segmentSlugs[i] = seg.Slug
		}

		response := GetSegmentsResponse{
			Segments: segmentSlugs,
		}

		render.JSON(w, r, response)
	}
}
