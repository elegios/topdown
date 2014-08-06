package types

type Update struct {
	Pos     Position
	Content interface{}
}

type SpeechCharUpdate struct {
	Speech    string `json:"speech"`
	Character string `json:"character"`
}

type SpeechPointUpdate struct {
	Speech string   `json:"speech"`
	Pos    Position `json:"position"`
}

type NudgeUpdate struct {
	Nudge     string  `json:"nudge"`
	Amount    float32 `json:"amount"`
	Character string  `json:"character,omitempty"`
	Target    string  `json:"target"`
}
