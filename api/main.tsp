import "@typespec/http";

using Http;
@service(#{ title: "Seven Skies Service" })
namespace SevenSkiesService;

@doc("섭씨 온도(℃)")
scalar Celsius extends numeric;

enum WeatherCondition {
  @doc("맑음") Clear: "CLEAR",
  @doc("구름많음") Cloudy: "CLOUDY",
  @doc("비") Rain: "RAIN",
  @doc("눈") Snow: "SNOW",
  @doc("비/눈") RainSnow: "RAIN_SNOW",
}

model DailyForecast {
  @doc("해당 날짜")
  date: offsetDateTime;

  @doc("하루 동안의 예상 최고 기온(℃)")
  high: Celsius;

  @doc("하루 동안의 예상 최저 기온(℃)")
  low: Celsius;

  @doc("하루 전체를 대표하는 날씨 상태")
  condition: WeatherCondition;
}

model DailyForecastList {
  items: DailyForecast[];
}

model HourlyForecast {
  @doc("예상 시간")
  time: offsetDateTime;

  @doc("1시간 동안의 예상 기온(℃)")
  temperature: Celsius;

  @doc("1시간 동안의 예상 날씨 상태")
  condition: WeatherCondition;
}

model HourlyForecastList {
  items: HourlyForecast[];
}

@error
model Error {
  @doc("에러 코드")
  errorCode: int32;

  @doc("에러 메시지")
  errorMessage: string;
}

@route("/forecasts")
@tag("Forecasts")
interface Forecasts {
  @get
  @route("/daily")
  listDaily(): {
    @statusCode statusCode: 200;
    @body body: DailyForecastList;
  } | {
    @statusCode statusCode: 500;
    @body body: Error;
  };

  @get
  @route("/hourly")
  listHourly(): {
    @statusCode statusCode: 200;
    @body body: HourlyForecastList;
  } | {
    @statusCode statusCode: 500;
    @body body: Error;
  };
}
