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
