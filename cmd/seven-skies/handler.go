package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/neatflowcv/seven-skies/api"
	"github.com/oapi-codegen/runtime/types"
)

var _ api.StrictServerInterface = (*Handler)(nil)

type Handler struct{}

func NewHandler() http.Handler {
	handler := &Handler{}

	return api.HandlerFromMux(api.NewStrictHandler(handler, nil), chi.NewMux())
}

func (h *Handler) GetSevenSkiesV1Skies( //nolint:ireturn
	ctx context.Context,
	_ api.GetSevenSkiesV1SkiesRequestObject,
) (api.GetSevenSkiesV1SkiesResponseObject, error) {
	return api.GetSevenSkiesV1Skies200JSONResponse([]api.Sky{
		{
			Date:    types.Date{Time: time.Now()},
			High:    20, //nolint:mnd
			Low:     10, //nolint:mnd
			Weather: api.Weather("Sunny"),
		},
	}), nil
}
