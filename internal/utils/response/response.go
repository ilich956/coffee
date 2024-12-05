package response

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type Error struct {
	ErrorDescription string `json:"ErrorDescription"`
}

type Message struct {
	MessageDescription string `json:"MessageDescription"`
}

func SendError(w http.ResponseWriter, statusCode int, errorDescription string, errValue error) {
	var errMsg string
	if errValue == nil {
		errMsg = errorDescription
	} else {
		errMsg = fmt.Sprintf(errorDescription + ": " + errValue.Error())
	}

	error := Error{
		ErrorDescription: errMsg,
	}

	output, err := json.MarshalIndent(error, "", "  ")
	if err != nil {
		slog.Error("Failed to marshal error ", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(output)
}

func SendMessage(w http.ResponseWriter, statusCode int, messageDescription string) {
	message := Message{
		MessageDescription: messageDescription,
	}

	output, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		slog.Error("Failed to marshal error ", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(output)
}
