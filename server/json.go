package main

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
)

// requestのBodyをJSONとして読み込み、targetにデコードをする
func readRequestJSON(req *http.Request, target any) error {
	contentType := req.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return err
	}
	if mediaType != "application/json" {
		return fmt.Errorf("expected Content-Type: application/json, got %q", contentType)
	}

	decode := json.NewDecoder(req.Body)
	decode.DisallowUnknownFields()
	return decode.Decode(target)
}

// JSONをレスポンスとして返す
func renderJSON(w http.ResponseWriter, v any) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
