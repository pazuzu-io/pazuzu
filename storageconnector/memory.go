package storageconnector

type Memory struct {
	features []Feature
}

func NewMemory(features []Feature) *Memory {
	return &Memory{features}
}