package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/quic-go/quic-go/http3"
)

func LogFile(server *http3.Server, level string) {

	// Create or open log.txt file
	file, f_err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if f_err != nil {
		log.Fatalln(f_err.Error())
	}

	// Handle logging
	switch level {
	case "Debug":
		server.Logger = slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "Error":
		server.Logger = slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelError}))
	case "Info":
		server.Logger = slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case "Warn":
		server.Logger = slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelWarn}))
	default:
		log.Fatalln("Invalid Logging level")
	}
}
