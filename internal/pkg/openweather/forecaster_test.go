package openweather_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/neatflowcv/seven-skies/internal/pkg/openweather"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/forecast.json
var forecast []byte

func TestForecast(t *testing.T) {
	t.Parallel()

	var resp openweather.ForecastResponse

	err := json.Unmarshal(forecast, &resp)

	require.NoError(t, err)
	require.NotNil(t, forecast)
}
