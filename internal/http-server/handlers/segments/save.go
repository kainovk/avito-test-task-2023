package segments

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	resp "avito-test-task-2023/internal/lib/api/response"
	"avito-test-task-2023/internal/lib/logger/sl"
	"avito-test-task-2023/internal/storage"
)

type SaveRequest struct {
	Name string `json:"name" validate:"required"`
}

type SaveResponse struct {
	resp.Response
}

type SegmentSaver interface {
	SaveSegment(name string) error
}

func NewSegmentSaver(log *slog.Logger, segmentSaver SegmentSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.segments.save.NewSegmentSaver"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req SaveRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))
			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		err = segmentSaver.SaveSegment(req.Name)
		if errors.Is(err, storage.ErrSegmentExists) {
			log.Info("segment already exists", slog.String("name", req.Name))

			render.JSON(w, r, resp.Error("segment already exists"))
			return
		}
		if err != nil {
			log.Error("failed to create segment", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to create segment"))
			return
		}

		log.Info("segment created")

		render.JSON(w, r, SaveResponse{
			Response: resp.OK(),
		})
	}
}
