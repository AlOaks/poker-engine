package models

type Player struct {
	ID        int     `json:"id"`
	HoleCards [2]Card `json:"holeCards"`
	Folded    bool    `json:"folded"`
	Equity    float64 `json:"equity"`
	Winner    bool    `json:"winner"`
}
