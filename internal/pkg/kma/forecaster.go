package kma

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"resty.dev/v3"
)

const (
	UltraShortForecastInterval       = time.Hour
	UltraShortForecastAvailableAfter = 15 * time.Minute

	ShortTermForecastInterval       = 3 * time.Hour
	ShortTermForecastAvailableAfter = 10 * time.Minute
)

func GetUltraShortForecast(ctx context.Context, key string, base time.Time, nx int, ny int) (*ForecastResponse, error) {
	// 초단기 예보는 0:30, 1:30, ..., 11:30 단위로 제공된다.
	// 요청 시각(`base`)을 기준으로, 그보다 크지 않으면서 가장 가까운 기준 시각을 선택한다.
	selectedBase := selectUltraShortBase(base)

	client := resty.New()

	defer func() {
		err := client.Close()
		if err != nil {
			log.Println("error in close client", err)
		}
	}()

	resp, err := client.R().
		SetContext(ctx).
		SetQueryParam("pageNo", "1").
		SetQueryParam("numOfRows", "1000").
		SetQueryParam("dataType", "JSON").
		SetQueryParam("base_date", selectedBase.Format("20060102")).
		SetQueryParam("base_time", selectedBase.Format("1504")).
		SetQueryParam("nx", strconv.Itoa(nx)).
		SetQueryParam("ny", strconv.Itoa(ny)).
		SetQueryParam("authKey", key).
		SetResult(&ForecastResponse{}). //nolint:exhaustruct
		Get("https://apihub.kma.go.kr/api/typ02/openApi/VilageFcstInfoService_2.0/getUltraSrtFcst")
	if err != nil {
		return nil, fmt.Errorf("error in get forecast: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("%w: %s, %s", ErrForecast, resp.Status(), resp.String())
	}

	return resp.Result().(*ForecastResponse), nil //nolint:forcetypeassert
}

func GetShortTermForecast(ctx context.Context, key string, base time.Time, nx int, ny int) (*ForecastResponse, error) {
	selectedBase := selectShortTermBase(base)

	client := resty.New()

	defer func() {
		err := client.Close()
		if err != nil {
			log.Println("error in close client", err)
		}
	}()

	resp, err := client.R().SetDebug(true).
		SetContext(ctx).
		SetQueryParam("pageNo", "1").
		SetQueryParam("numOfRows", "1000").
		SetQueryParam("dataType", "JSON").
		SetQueryParam("base_date", selectedBase.Format("20060102")).
		SetQueryParam("base_time", selectedBase.Format("1504")).
		SetQueryParam("nx", strconv.Itoa(nx)).
		SetQueryParam("ny", strconv.Itoa(ny)).
		SetQueryParam("authKey", key).
		SetResult(&ForecastResponse{}). //nolint:exhaustruct
		Get("https://apihub.kma.go.kr/api/typ02/openApi/VilageFcstInfoService_2.0/getVilageFcst")
	if err != nil {
		return nil, fmt.Errorf("error in get forecast: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("%w: %s, %s", ErrForecast, resp.Status(), resp.String())
	}

	return resp.Result().(*ForecastResponse), nil //nolint:forcetypeassert
}

// 초단기 예보 기준 시각(0:30, 1:30, ..., 11:30) 중에서,
// 주어진 시각 t 기준으로 t 보다 크지 않으면서 가장 가까운 시각을 반환한다.
func selectUltraShortBase(base time.Time) time.Time {
	// 초단기 예보:
	// - 기준 시각 간격: 1시간
	// - 조회 가능 시각: 기준 시각 + 15분
	// - 후보 시작 시각: 이전날 23:30
	const (
		startHour   = 23
		startMinute = 30
	)

	return selectBase(
		base,
		UltraShortForecastInterval,
		UltraShortForecastAvailableAfter,
		startHour,
		startMinute,
	)
}

// 단기 예보 기준 시각(02, 05, 08, 11, 14, 17, 20, 23) 중에서,
// 주어진 시각 t 기준으로, t 시각에서 조회 가능한 가장 최근 기준 시각을 반환한다.
// 각 기준 시각은 기준 시각 + 10분 이후부터 조회 가능하다.
func selectShortTermBase(base time.Time) time.Time {
	// 단기 예보:
	// - 기준 시각 간격: 3시간
	// - 조회 가능 시각: 기준 시각 + 10분
	// - 후보 시작 시각: 이전날 23:00
	const (
		startHour   = 23
		startMinute = 0
	)

	return selectBase(
		base,
		ShortTermForecastInterval,
		ShortTermForecastAvailableAfter,
		startHour,
		startMinute,
	)
}

func selectBase(base time.Time, interval, availableAfter time.Duration, startHour, startMinute int) time.Time {
	// KMA API 는 한국 표준시(Asia/Seoul)를 기준으로 동작하므로,
	// base 가 다른 타임존(예: UTC)인 경우 한국 시간으로 변환해서 계산한다.
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		// 로케일 정보를 불러오지 못하면, 기존 base 의 로케이션을 그대로 사용한다.
		loc = base.Location()
	}

	baseInKST := base.In(loc)

	prevDay := baseInKST.AddDate(0, 0, -1)

	start := time.Date(prevDay.Year(), prevDay.Month(), prevDay.Day(), startHour, startMinute, 0, 0, loc)
	end := baseInKST

	var candidates []time.Time

	for t := start; !t.After(end); t = t.Add(interval) {
		candidates = append(candidates, t)
	}

	return selectBaseFromCandidates(baseInKST, candidates, availableAfter)
}

// 공통 기준 시각 선택 로직:
// - candidates: 오름차순(과거 -> 미래) 기준 시각 목록
// - 각 기준 시각은 availableAfter 이후부터 조회 가능
// - base 시각 기준으로, 조회 가능한 가장 최근 기준 시각을 반환한다.
func selectBaseFromCandidates(base time.Time, candidates []time.Time, availableAfter time.Duration) time.Time {
	// 가장 최근 기준 시각을 찾기 위해 역순으로 순회한다.
	for i := len(candidates) - 1; i >= 0; i-- {
		candidate := candidates[i]
		availableAt := candidate.Add(availableAfter)

		if base.Before(availableAt) {
			continue
		}

		return candidate
	}

	panic("no available base time found")
}
