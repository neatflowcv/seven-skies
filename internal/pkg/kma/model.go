package kma

type ForecastResponse struct {
	Response Response `json:"response,omitzero"`
}

type Header struct {
	ResultCode string `json:"resultCode,omitempty"`
	ResultMsg  string `json:"resultMsg,omitempty"`
}

type Item struct {
	BaseDate  string `json:"baseDate,omitempty"`
	BaseTime  string `json:"baseTime,omitempty"`
	Category  string `json:"category,omitempty"`
	FcstDate  string `json:"fcstDate,omitempty"`
	FcstTime  string `json:"fcstTime,omitempty"`
	FcstValue string `json:"fcstValue,omitempty"`
	Nx        int    `json:"nx,omitempty"`
	Ny        int    `json:"ny,omitempty"`
}

type Items struct {
	Item []Item `json:"item,omitempty"`
}

type Body struct {
	DataType   string `json:"dataType,omitempty"`
	Items      Items  `json:"items,omitzero"`
	PageNo     int    `json:"pageNo,omitempty"`
	NumOfRows  int    `json:"numOfRows,omitempty"`
	TotalCount int    `json:"totalCount,omitempty"`
}

type Response struct {
	Header Header `json:"header,omitzero"`
	Body   Body   `json:"body,omitzero"`
}
