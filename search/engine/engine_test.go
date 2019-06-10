package engine

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/go-cmp/cmp"
)

type MockDock struct {
	ID      string
	isValid bool
	value   interface{}
}

func (m MockDock) GetID() string {
	return m.ID
}

func (m MockDock) Content() (interface{}, error) {
	if !m.isValid {
		return nil, errors.New("error")
	}
	return m.value, nil
}

func TestIndexing(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Error(err)
		return
	}

	docs := []MockDock{
		{ID: "1", isValid: true, value: map[string]string{"key": "value"}},
		{ID: "2", isValid: false, value: map[string]string{"key": "value"}},
	}

	for i, v := range docs {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			err := engine.Index(v)
			if v.isValid && err != nil {
				t.Error(err)
			}
			if !v.isValid && err == nil {
				t.Error("an error should be returned")
			}

		})
	}
}

func TestFinding(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Error(err)
		return
	}

	docs := []MockDock{
		{ID: "1", isValid: true, value: map[string]string{"name": "john", "surname": "smith"}},
		{ID: "2", isValid: false, value: map[string]string{"name": "john", "surname": "renbourn"}},
		{ID: "3", isValid: true, value: map[string]string{"name": "michael", "surname": "smith"}},
	}

	for _, v := range docs {
		engine.Index(v)
	}

	queries := map[string][]string{
		"michael":  []string{"3"},
		"john":     []string{"1"},
		"smith":    []string{"1", "3"},
		"renbourn": []string{},
	}

	for i, v := range queries {
		t.Run(i, func(t *testing.T) {
			res, err := engine.Find(i)
			if err != nil {
				t.Error(err)
			}
			if !cmp.Equal(res, v, cmpopts.EquateEmpty()) {
				t.Error(res, v, cmp.Diff(res, v))
			}
		})
	}
}
