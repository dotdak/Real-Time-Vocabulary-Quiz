package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dotdak/Real-Time-Vocabulary-Quiz/controllers"
	"github.com/dotdak/Real-Time-Vocabulary-Quiz/services"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	s := controllers.NewServer(ctx)
	s.RegisterHandler()

	httpPortEnv := os.Getenv("HTTP_PORT")
	if httpPortEnv == "" {
		panic("$HTTP_PORT must be set")
	}

	httpPort, err := strconv.Atoi(httpPortEnv)
	if err != nil {
		panic(err)
	}

	httpsPortEnv := os.Getenv("HTTPS_PORT")
	if httpsPortEnv == "" {
		panic("$HTTPS_PORT must be set")
	}

	httpsPort, err := strconv.Atoi(httpsPortEnv)
	if err != nil {
		panic(err)
	}

	_ = s.CreateQuizSession("123", &services.QuizConfig{
		LifeTime:      15 * time.Minute,
		QuestionTimer: 3 * time.Second,
		TotalQuestion: 10,
	})

	go func(ctx context.Context) {
		defer cancel()
		if err := s.Echo.StartTLS(fmt.Sprintf(":%d", httpsPort), "localhost.crt", "localhost.key"); err != http.ErrServerClosed {
			panic(err)
		}
	}(ctx)

	if err := s.Echo.Start(fmt.Sprintf(":%d", httpPort)); err != http.ErrServerClosed {
		panic(err)
	}
}
