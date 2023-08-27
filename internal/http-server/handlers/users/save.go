package save

import (
	"avito-test-task-2023/internal/storage"
	"errors"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"

	resp "avito-test-task-2023/internal/lib/api/response"
	"avito-test-task-2023/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Name string `json:"name" validate:"required"`
}

type Response struct {
	resp.Response
	ID int64 `json:"user_id,omitempty"`
}

type UserSaver interface {
	SaveUser(name string) error
}

func New(log *slog.Logger, userSaver UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

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

		err = userSaver.SaveUser(req.Name)
		if errors.Is(err, storage.ErrUserExists) {
			log.Info("user already exists", slog.String("name", req.Name))

			render.JSON(w, r, resp.Error("user already exists"))
			return
		}
		if err != nil {
			log.Error("failed to create user", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to create user"))
			return
		}

		log.Info("user created")

		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}
