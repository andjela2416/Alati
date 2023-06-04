package main

import (
	"encoding/json"
	"example.com/mod/store"
	"github.com/google/uuid"
	"io"
	"net/http"
)

func decodeBody(r io.Reader) (*store.Config, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var c store.Config
	if err := dec.Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func decodeGroup(r io.Reader) (*store.Group, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var g store.Group
	if err := dec.Decode(&g); err != nil {
		return nil, err
	}
	return &g, nil
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
func createId() string {
	return uuid.New().String()
}
