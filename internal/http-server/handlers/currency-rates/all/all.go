package all

import (
	"errors"
	"exchange-rates/internal/exrate"
	"exchange-rates/internal/lib/logger/sl"
	"exchange-rates/internal/storage"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"

	resp "exchange-rates/internal/lib/api/response"
)

type CurRatesGetter interface {
	GetAllCurrencyRates() ([]exrate.Options, error)
}

type Response struct {
	resp.Response
	CurRates []exrate.Options `json:"cur_rates,omitempty"`
}

func New(log *slog.Logger, cRatGetter CurRatesGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.current-rates.all.New"

		log = log.With(slog.String("op", op))

		cRates, err := cRatGetter.GetAllCurrencyRates()
		if err != nil {
			if errors.Is(err, storage.ErrCurRatesNotFound) {
				log.Info(err.Error())

				render.JSON(w, r, resp.Error(err.Error()))

				return
			}

			log.Error("failed to get exchange rates", sl.Error(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
			CurRates: cRates,
		})
	}
}
