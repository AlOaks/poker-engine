package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONColumn[T any] []T

func (j *JSONColumn[T]) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion failed")
	}

	return json.Unmarshal(bytes, j)
}

func (j JSONColumn[T]) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func UniqueSet(xs []Card) []Card {
	seen := make(map[Rank]struct{}, len(xs))
	out := make([]Card, 0, len(xs))
	for _, x := range xs {
		if _, ok := seen[x.Rank]; ok {
			continue
		}
		seen[x.Rank] = struct{}{}
		out = append(out, x)
	}
	return out
}
