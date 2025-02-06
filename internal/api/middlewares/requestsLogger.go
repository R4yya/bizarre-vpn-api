package middlewares

import (
	"bizarre-vpn-api/internal/lib/formatting"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestsLogger(log *slog.Logger) gin.HandlerFunc {
	const op = "internal.api.middlewares.RequestsLogger"

	entry := log.With(
		slog.String("op", op),
	)

	entry.Info("logger middleware enabled")

	return func(c *gin.Context) {
		r := c.Request

		t1 := time.Now()

		c.Next()

		var builder strings.Builder

		builder.WriteString("|")
		builder.WriteString(formatting.PadString(fmt.Sprintf("%v", c.Writer.Status()), 3))
		builder.WriteString("|")
		builder.WriteString(formatting.PadString(time.Since(t1).String(), 14))
		builder.WriteString("|")
		builder.WriteString(formatting.PadString(r.RemoteAddr, 15))
		builder.WriteString("|")
		builder.WriteString(formatting.PadString(r.Method, 6))
		builder.WriteString("| ")
		builder.WriteString(r.URL.Path)

		log.Info(builder.String())
	}
}
