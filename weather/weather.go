package weather

type Temperature struct {
	Temp float32 `json:"temp,omitempty"`
	Date string  `json:"date,omitempty"`
}

type Windspeed struct {
	North float32 `json:"north,omitempty"`
	West  float32 `json:"west,omitempty"`
	Date  string  `json:"date,omitempty"`
}

type Weather struct {
	North float32 `json:"north,omitempty"`
	West  float32 `json:"west,omitempty"`
	Temp  float32 `json:"temp,omitempty"`
	Date  string  `json:"date,omitempty"`
}
