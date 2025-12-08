//nolint:ireturn
package main

import (
	"context"
	"net/http"
	"time"

	"github.com/neatflowcv/seven-skies/api"
	"github.com/neatflowcv/seven-skies/internal/app/flow"
)

var _ api.StrictServerInterface = (*Handler)(nil)

type Handler struct {
	flow *flow.Flow
}

func NewHandler(flow *flow.Flow) *Handler {
	return &Handler{
		flow: flow,
	}
}

func (h *Handler) ForecastsListDaily(
	ctx context.Context,
	req api.ForecastsListDailyRequestObject,
) (api.ForecastsListDailyResponseObject, error) {
	timezone := "Asia/Seoul"
	if req.Params.Timezone != nil {
		timezone = *req.Params.Timezone
	}

	now := time.Now()
	if req.Params.Now != nil {
		now = *req.Params.Now
	}

	weathers, err := h.flow.ListDailyWeathers(ctx, timezone, now)
	if err != nil {
		return api.ForecastsListDaily500JSONResponse{ //nolint:nilerr
			ErrorCode:    http.StatusInternalServerError,
			ErrorMessage: err.Error(),
		}, nil
	}

	var items []api.DailyForecast
	for _, weather := range weathers {
		items = append(items, api.DailyForecast{
			Date:      weather.Date,
			High:      api.Celsius(weather.High),
			Low:       api.Celsius(weather.Low),
			Condition: api.WeatherCondition(weather.Condition),
		})
	}

	return api.ForecastsListDaily200JSONResponse{
		Items: items,
	}, nil
}

func (h *Handler) ForecastsListHourly(
	ctx context.Context,
	req api.ForecastsListHourlyRequestObject,
) (api.ForecastsListHourlyResponseObject, error) {
	panic("unimplemented")
}
