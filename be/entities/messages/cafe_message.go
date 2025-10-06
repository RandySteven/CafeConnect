package messages

type CafeMessage struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    uint64  `json:"radius"`
}
