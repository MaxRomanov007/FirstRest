package models

type Car struct {
	Producer       string  `json:"producer"`
	Model          string  `json:"model"`
	EngineCapacity float32 `json:"engine_capacity"`
	Power          float32 `json:"power"`
	Number         string  `json:"number"`
	ImagesCount    uint8   `json:"images_count"`
	Description    string  `json:"description"`
}
