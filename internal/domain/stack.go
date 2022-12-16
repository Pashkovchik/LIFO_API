package domain

import "time"

type StackData struct {
	ID          uint64    `json:"id,omitempty"`
	StackID     uint64    `json:"stackID,omitempty"`
	Info        string    `json:"info"`
	CreatedDate time.Time `json:"createdDate,omitempty"`
}

func (sd *StackData) Fields() []interface{} {
	return []interface{}{
		&sd.ID,
		&sd.StackID,
		&sd.Info,
		&sd.CreatedDate,
	}
}

type Stack struct {
	ID          uint64      `json:"id,omitempty"`
	Name        string      `json:"name"`
	CreatedDate time.Time   `json:"createdDate,omitempty"`
	StackData   []StackData `json:"stackData,omitempty"`
}

func (s *Stack) Fields() []interface{} {
	return []interface{}{
		&s.ID,
		&s.Name,
		&s.CreatedDate,
	}
}
