# 기상 예보 수집처

기상 예보 데이터 업데이트 간격이 다르고, 요구하는 입출력도 너무 다르다.

모노레포 형식으로 기상 예보를 수집하는 별도의 서비스들과 기상 데이터를 관리하는 서비스를 분리한다.

## 기상청 API 허브

- <https://apihub.kma.go.kr/>

### 기상청 API

많은 API가 제공되나, 기상 예보를 위해서는 2개만 보면 된다.

#### 초단기 예보

<https://apihub.kma.go.kr/api/typ02/openApi/VilageFcstInfoService_2.0/getUltraSrtFcst>

- 제공하는 범위가 짧다. TODO: 조사 보충 필요
- 00:30 부터 1시간 간격으로 새로운 발표를 한다.
  - 00:30, 01:30, ..., 23:30
- 발표 시간으로부터 15분 후에 해당 데이터를 조회할 수 있다.

##### 출력 카테고리 - 초단기

| 코드 | 의미         | 단위                                                                   |
| ---- | ------------ | ---------------------------------------------------------------------- |
| T1H  | 기온         | ℃                                                                      |
| REH  | 습도         | %                                                                      |
| RN1  | 1시간 강수량 | mm                                                                     |
| SKY  | 하늘상태     | 맑음(1), 구름많음(3), 흐림(4)                                          |
| PTY  | 강수형태     | 없음(0), 비(1), 비/눈(2), 눈(3), 빗방울(5), 빗방울눈날림(6), 눈날림(7) |

#### 단기 예보

<https://apihub.kma.go.kr/api/typ02/openApi/VilageFcstInfoService_2.0/getVilageFcst>

- 제공하는 범위가 대충 1주일 되는 것 같다. TODO: 조사 보충 필요
- 02:00시부터 3시간 간격으로 새로운 발표를 한다.
  - 02:00, 05:00, ..., 23:00
- 발표 시간으로부터 10분 후에 해당 데이터를 조회할 수 있다.

##### 출력 카테고리 - 단기

| 코드 | 의미         | 단위                                       |
| ---- | ------------ | ------------------------------------------ |
| TMP  | 1시간 기온   | ℃                                          |
| TMN  | 일 최저기온  | ℃                                          |
| TMX  | 일 최고기온  | ℃                                          |
| REH  | 습도         | %                                          |
| POP  | 강수확률     | %                                          |
| PTY  | 강수형태     | 없음(0), 비(1), 비/눈(2), 눈(3), 소나기(4) |
| PCP  | 1시간 강수량 | mm                                         |
| SNO  | 1시간 신적설 | cm                                         |
| SKY  | 하늘상태     | enum: 맑음(1), 구름많음(3), 흐림(4)        |

#### 입력 (공통)

- pageNo: 페이지 번호, 1부터 시작
- numOfRows: 한 페이지 결과 수
  - 예시: 1000
- dataType: 반환 형식
  - enum: XML, JSON
- base_date: 발표 일자, 해당 발표시간껄 보여준다는 건지? 이후라는건지? TODO: 조사 보충 필요
  - 예시: 20251203
- base_time: 발표 시간, 해당 발표시간껄 보여준다는 건지? 이후라는건지? TODO: 조사 보충 필요
  - 예시: 0630
  - API 마다 발표 시간이 정해져 있다.
- nx: 예보지점 X 좌표
- ny: 예보지점 Y 좌표
- authKey: 인증키

> 예보지점 (nx, ny)는 (위도, 경도)가 아니다.
> 홈페이지에서 동네예보 지점 좌표(위경도) 참고자료를 다운로드 받아 확인할 수 있다.

### 평가

다소 구식 인터페이스이지만, 한국에서 사실상 유일한 공공 기상 데이터 출처이므로 필수적으로 활용해야 한다.

## OpenWeather

- <https://openweathermap.org/>

### 사용법 및 제한

<https://openweathermap.org/api> 여기에 접속해서 일단 가입해야 한다.
가입 후 몇 시간 정도는 지나야 API 사용할 수 있다.

<https://openweathermap.org/full-price#current> 여기에서 Free 티어가 사용가능한 API가 명확히 나온다.

- 3-hour Forecast 5 days 이거 하나만 쓰면 된다.
- 2시간 마다 데이터가 업데이트 된다.
  - 즉, 하루 12번만 호출하면 된다.
- 1분에 최대 60회 가능하며, 한 달에 1,000,000회 호출이 가능하다.
  - (1,000,000 / 30)

### OpenWeather API

#### Call 5 day / 3 hour forecast data

<https://openweathermap.org/forecast5>

형식: api.openweathermap.org/data/2.5/forecast?units=metric&lat={lat}&lon={lon}&appid={API key}

<https://openweathermap.org/weather-conditions>
