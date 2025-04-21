package utils

// Data structure for individual Kline/Candlestick
// Added Hour field
type Kline struct {
	OpenTime  int64  `json:"openTime"`
	Open      string `json:"open"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Close     string `json:"close"`
	CloseTime int64  `json:"closeTime"`
	Date      string `json:"date"` // Human-readable date YYYY-MM-DD
	Hour      int    `json:"hour"` // Hour of the day (0-23)
}
