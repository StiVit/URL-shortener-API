package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/StiVit/URL-shortener-API/internal/lib/api/response"
	"github.com/StiVit/URL-shortener-API/internal/lib/logger/sl"
	"github.com/StiVit/URL-shortener-API/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("Alias is empty")
			render.JSON(w, r, response.Error("not found"))
			return 
		}

		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("URL not found", "alias", alias)

			render.JSON(w, r, response.Error("Not found"))

			return
		}
		if err != nil {
			log.Error("failed to get URL", sl.Err(err))

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		log.Info("Got URL", slog.String("url", resURL))

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}