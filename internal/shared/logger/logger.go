package logger

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"

	"bizarre-vpn-api/internal/shared/config"
	"bizarre-vpn-api/internal/shared/logger/handlers/slogpretty"
)

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	// Настройка ротации логов в файл
	fileLogger := &lumberjack.Logger{
		Filename:   "logs/app.log", // Имя файла
		MaxSize:    10,             // Максимальный размер файла в мегабайтах
		MaxBackups: 5,              // Максимальное количество старых файлов
		MaxAge:     28,             // Максимальное количество дней хранения
		Compress:   true,           // Сжатие старых файлов
	}

	switch env {
	case config.EnvLocal:
		log = setupPrettySlog()
	case config.EnvDev:
		// Создаем мультиплексор для вывода логов в stdout и файл
		multiWriter := io.MultiWriter(os.Stdout, fileLogger)

		log = slog.New(
			slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}),
		)

	case config.EnvProd:
		// Создаем мультиплексор для вывода логов в stdout и файл
		multiWriter := io.MultiWriter(os.Stdout, fileLogger)

		log = slog.New(
			slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
