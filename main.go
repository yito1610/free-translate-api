package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bregydoc/gtranslate"
	"github.com/rs/cors"
)

type TranslateResponse struct {
	TranslatedText string `json:"translatedText,omitempty"`
	Status         bool   `json:"status"`
	Message        string `json:"message"`
}

func TranslateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	text := r.URL.Query().Get("text")
	to := r.URL.Query().Get("to")

	if text == "" || to == "" {
		sendErrorResponse(w, "Text and to parameters are required", http.StatusBadRequest)
		return
	}

	translated, err := gtranslate.TranslateWithParams(text, gtranslate.TranslationParams{
		From: "auto",
		To:   to,
	})
	if err != nil {
		sendErrorResponse(w, "Translation failed", http.StatusInternalServerError)
		return
	}

	response := TranslateResponse{
		TranslatedText: translated,
		Status:         true,
		Message:        "",
	}

	sendJSONResponse(w, response, http.StatusOK)
}

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := TranslateResponse{
		Status:  false,
		Message: message,
	}

	sendJSONResponse(w, response, statusCode)
}

func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/translate", TranslateHandler)

	handler := cors.Default().Handler(mux)

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
