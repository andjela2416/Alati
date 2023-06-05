package main

import (
	"context"
	"encoding/json"
	"example.com/mod/store"
	tracer "example.com/mod/tracer"
	"github.com/google/uuid"
	"io"
	"net/http"
)

func decodeBody(ctx context.Context, r io.Reader) (*store.Config, error) {
	span := tracer.StartSpanFromContext(ctx, "decodeBody")
	defer span.Finish()

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var c store.Config
	if err := dec.Decode(&c); err != nil {
		tracer.LogError(span, err)
		return nil, err
	}
	return &c, nil
}

func decodeGroup(ctx context.Context, r io.Reader) (*store.Group, error) {
	span := tracer.StartSpanFromContext(ctx, "decodeBody")
	defer span.Finish()

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var g store.Group
	if err := dec.Decode(&g); err != nil {
		tracer.LogError(span, err)
		return nil, err
	}
	return &g, nil
}

func renderJSON(ctx context.Context, w http.ResponseWriter, v interface{}) {
	span := tracer.StartSpanFromContext(ctx, "decodeBody")
	defer span.Finish()

	js, err := json.Marshal(v)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
func createId() string {
	return uuid.New().String()
}
