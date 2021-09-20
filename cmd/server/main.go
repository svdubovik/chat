package main

import (
	"svdubovik.com/chat/internal/config"
	"svdubovik.com/chat/internal/logger"
	"svdubovik.com/chat/internal/server/tcp"
)

func main() {
	const service = "ChatServer"

	cfg := config.NewConfig(service)
	logger := logger.NewLogger(cfg.LogLevel, cfg.LogFormat, cfg.Service)

	server := tcp.NewServer(cfg, logger)

	if msg, err := server.Run(); err != nil {
		logger.Fatal().Err(err).Msg(msg)
	}
}
