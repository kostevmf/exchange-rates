package ondate

import (
	"errors"
	"exchange-rates/internal/exrate"
	"exchange-rates/internal/lib/logger/sl"
	"exchange-rates/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"time"

	resp "exchange-rates/internal/lib/api/response"
)

type CurRatesByDateGetter interface {
	GetCurrencyRatesByDate(date time.Time) ([]exrate.Options, error)
}

type Response struct {
	resp.Response
	CurRates []exrate.Options `json:"cur_rates,omitempty"`
}

func New(log *slog.Logger, cRatGetter CurRatesByDateGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.current-rates.ondate.New"

		log = log.With(slog.String("op", op))

		date := chi.URLParam(r, "date")
		if date == "" {
			log.Info("date is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		t, err := time.Parse(time.DateOnly, date)
		if err != nil {
			log.Error("failed to parse date", "date", date)

			render.JSON(w, r, resp.Error("failed to parse date"))

			return
		}

		cRates, err := cRatGetter.GetCurrencyRatesByDate(t)
		if err != nil {
			if errors.Is(err, storage.ErrCurRatesNotFound) {
				log.Info(err.Error())

				render.JSON(w, r, resp.Error(err.Error()))

				return
			}

			log.Error("failed to get exchange rates by date", sl.Error(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
			CurRates: cRates,
		})
	}
}
