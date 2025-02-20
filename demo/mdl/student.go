package mdl

type Student struct {
	Id     uint64             `json:"id"`
	Name   string             `json:"name"`
	Age    uint8              `json:"age"`
	Class  string             `json:"class"`
	Scores map[string]float32 `json:"scores"`
}
