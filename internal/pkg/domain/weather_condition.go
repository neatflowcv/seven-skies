package domain

type WeatherCondition string

const (
	WeatherConditionUnknown  WeatherCondition = "UNKNOWN"
	WeatherConditionClear    WeatherCondition = "CLEAR"
	WeatherConditionCloudy   WeatherCondition = "CLOUDY"
	WeatherConditionRain     WeatherCondition = "RAIN"
	WeatherConditionSnow     WeatherCondition = "SNOW"
	WeatherConditionRainSnow WeatherCondition = "RAIN_SNOW"
)

// IsWorseThan은 현재 날씨 조건이 other보다 나쁜 조건인지 확인합니다.
// 우선순위: RAIN > RAIN_SNOW > SNOW > CLOUDY > CLEAR.
func (c WeatherCondition) IsWorseThan(other WeatherCondition) bool {
	return c.priority() > other.priority()
}

// priority는 날씨 조건의 우선순위를 반환합니다.
// 우선순위가 높을수록 더 나쁜 날씨 조건입니다.
// 우선순위: RAIN(4) > RAIN_SNOW(3) > SNOW(2) > CLOUDY(1) > CLEAR(0) > UNKNOWN(-1).
func (c WeatherCondition) priority() int {
	switch c {
	case WeatherConditionClear:
		return 0
	case WeatherConditionCloudy:
		return 1
	case WeatherConditionSnow:
		return 2 //nolint:mnd
	case WeatherConditionRainSnow:
		return 3 //nolint:mnd
	case WeatherConditionRain:
		return 4 //nolint:mnd
	case WeatherConditionUnknown:
		return -1
	default:
		return -1
	}
}
