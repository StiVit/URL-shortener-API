package delete

import (
	"log/slog"
	"net/http"

	resp "github.com/StiVit/URL-shortener-API/internal/lib/api/response"
	"github.com/StiVit/URL-shortener-API/internal/lib/logger/sl"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLDeleter interface {
	DeleteURL(aliasToDelete string) (error)
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.delete.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("Alias is missing")
			render.JSON(w, r, resp.Error("alias is required"))
			return
		}

		log.Info("deleting URL", slog.String("alias", alias))

		err := urlDeleter.DeleteURL(alias)
		if err != nil {
			log.Error("Falied to delete URL", sl.Err(err))
			render.JSON(w, r, resp.Error("failed ot delete URL"))
			return
		}

		log.Info("URL deleted successfully", slog.String("alias", alias))
		render.JSON(w, r, resp.Ok())
	}
}