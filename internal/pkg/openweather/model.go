package openweather

type ForecastResponse struct {
	Cod     string `json:"cod,omitempty"`
	Message int    `json:"message,omitempty"`
	Cnt     int    `json:"cnt,omitempty"`
	List    []List `json:"list,omitempty"`
	City    City   `json:"city,omitzero"`
}

type Main struct {
	Temp      float64 `json:"temp,omitempty"`
	FeelsLike float64 `json:"feels_like,omitempty"`
	TempMin   float64 `json:"temp_min,omitempty"`
	TempMax   float64 `json:"temp_max,omitempty"`
	Pressure  int     `json:"pressure,omitempty"`
	SeaLevel  int     `json:"sea_level,omitempty"`
	GrndLevel int     `json:"grnd_level,omitempty"`
	Humidity  int     `json:"humidity,omitempty"`
	TempKf    float64 `json:"temp_kf,omitempty"`
}

type Weather struct {
	ID          int    `json:"id,omitempty"`
	Main        string `json:"main,omitempty"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
}

type Clouds struct {
	All int `json:"all,omitempty"`
}

type Wind struct {
	Speed float64 `json:"speed,omitempty"`
	Deg   int     `json:"deg,omitempty"`
	Gust  float64 `json:"gust,omitempty"`
}

type Sys struct {
	Pod string `json:"pod,omitempty"`
}

type List struct {
	Dt         int       `json:"dt,omitempty"`
	Main       Main      `json:"main,omitzero"`
	Weather    []Weather `json:"weather,omitempty"`
	Clouds     Clouds    `json:"clouds,omitzero"`
	Wind       Wind      `json:"wind,omitzero"`
	Visibility int       `json:"visibility,omitempty"`
	Pop        float64   `json:"pop,omitempty"`
	Sys        Sys       `json:"sys,omitzero"`
	DtTxt      string    `json:"dt_txt,omitempty"`
}

type Coord struct {
	Lat float64 `json:"lat,omitempty"`
	Lon float64 `json:"lon,omitempty"`
}

type City struct {
	ID         int    `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Coord      Coord  `json:"coord,omitzero"`
	Country    string `json:"country,omitempty"`
	Population int    `json:"population,omitempty"`
	Timezone   int    `json:"timezone,omitempty"`
	Sunrise    int    `json:"sunrise,omitempty"`
	Sunset     int    `json:"sunset,omitempty"`
}
