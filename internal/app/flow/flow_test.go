//nolint:dupl
package flow_test

import (
	"context"
	"testing"
	"time"

	"github.com/neatflowcv/seven-skies/internal/app/flow"
	"github.com/neatflowcv/seven-skies/internal/pkg/domain"
	"github.com/neatflowcv/seven-skies/internal/pkg/repository/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestFlow_ListDailyWeathers(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mock := mocks.NewMockRepository(ctrl)
	now := time.Now()
	mock.EXPECT().ListWeathers(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*domain.Weather{
		domain.NewWeather(
			"1",
			domain.WeatherSourceOpenWeather,
			now,
			now.Add(time.Hour*24),
			domain.WeatherConditionClear,
			domain.Temperature{Value: domain.Celsius(20)},
		),
	}, nil)
	flow := flow.NewFlow(mock)

	weathers, err := flow.ListDailyWeathers(context.Background(), "Asia/Seoul", now)

	require.NoError(t, err)
	require.Len(t, weathers, 1)
	require.InDelta(t, 20.0, weathers[0].High, 0.01)
	require.InDelta(t, 20.0, weathers[0].Low, 0.01)
	require.Equal(t, "CLEAR", weathers[0].Condition)
}

func TestFlow_ListDailyWeathers_SameSourceAndTargetDate(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mock := mocks.NewMockRepository(ctrl)
	now := time.Now()
	mock.EXPECT().ListWeathers(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*domain.Weather{
		domain.NewWeather(
			"1",
			domain.WeatherSourceOpenWeather,
			now,
			now,
			domain.WeatherConditionClear,
			domain.Temperature{Value: domain.Celsius(20)},
		),
		domain.NewWeather(
			"2",
			domain.WeatherSourceOpenWeather,
			now,
			now.Add(time.Hour*24),
			domain.WeatherConditionCloudy,
			domain.Temperature{Value: domain.Celsius(30)},
		),
	}, nil)
	flow := flow.NewFlow(mock)

	weathers, err := flow.ListDailyWeathers(context.Background(), "Asia/Seoul", time.Now())

	require.NoError(t, err)
	require.Len(t, weathers, 1)
	require.InDelta(t, 30.0, weathers[0].High, 0.01)
	require.InDelta(t, 30.0, weathers[0].Low, 0.01)
	require.Equal(t, "CLOUDY", weathers[0].Condition)
}

func TestFlow_ListDailyWeathers_DifferentSourceAndTargetDate(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mock := mocks.NewMockRepository(ctrl)
	now := time.Now()
	mock.EXPECT().ListWeathers(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*domain.Weather{
		domain.NewWeather(
			"1",
			domain.WeatherSourceOpenWeather,
			now,
			now,
			domain.WeatherConditionClear,
			domain.Temperature{Value: domain.Celsius(20)},
		),
		domain.NewWeather(
			"2",
			domain.WeatherSourceKMA,
			now,
			now.Add(time.Hour*24),
			domain.WeatherConditionCloudy,
			domain.Temperature{Value: domain.Celsius(30)},
		),
	}, nil)
	flow := flow.NewFlow(mock)

	weathers, err := flow.ListDailyWeathers(context.Background(), "Asia/Seoul", time.Now())

	require.NoError(t, err)
	require.Len(t, weathers, 1)
	require.InDelta(t, 30.0, weathers[0].High, 0.01)
	require.InDelta(t, 20.0, weathers[0].Low, 0.01)
	require.Equal(t, "CLOUDY", weathers[0].Condition)
}
