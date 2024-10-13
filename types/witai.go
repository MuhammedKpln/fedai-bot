package types

type WitAi struct {
	Speech struct {
		Confidence float64 `json:"confidence"`
		Tokens     []struct {
			Confidence float64 `json:"confidence"`
			End        int     `json:"end"`
			Start      int     `json:"start"`
			Token      string  `json:"token"`
		} `json:"tokens"`
	} `json:"speech"`
	Text    string `json:"text"`
	Type    string `json:"type"`
	IsFinal bool   `json:"is_final,omitempty"`
}
