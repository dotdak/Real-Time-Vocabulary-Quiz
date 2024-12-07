package controllers

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// go:embed
// var statics embed.FS

func (s *Server) RegisterHandler() {
	allowOrigins := make([]string, 0, 1)
	if os.Getenv("ENV") == "dev" {
		allowOrigins = append(allowOrigins, "http://localhost:3000")
	}
	s.Echo.Use(middleware.Logger())
	s.Echo.Use(middleware.Recover())
	s.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	s.Echo.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5: true,
		Root:  "build",
		// Filesystem: http.FS(statics),
	}))
	s.Echo.GET("/api/health", healthCheck)
	// s.Echo.POST("/api/feedback", handleFeedback(s))
	s.Echo.GET("/api/quiz/:quizId", handleConnection(s))
	s.Echo.POST("/api/quiz", handleQuizSession(s))
	s.Echo.PUT("/api/quiz/:quizId", handleQuizSessionStatusById(s))
	s.Echo.GET("/api/quiz/:quizId/overview", handleQuizSessionById(s))
	s.Echo.GET("/api/quiz/:quizId/new-question", handleNewQuestions(s))
	s.Echo.GET("/api/quiz/:quizId/users/:username", handleUserState(s))
	s.Echo.GET("/api/quiz/:quizId/leaderboard", handleLeaderBoard(s))
	// s.Echo.GET("/api/quiz/:quizId/questions/:questionId", handleQuestionById(s))

	// s.Echo.GET("/api/quiz-management/:quizId", handleQuizSessionState(s))
	s.Echo.PATCH("/api/quiz-management/:quizId/status", handleQuizStatus(s))

	// s.Echo.GET("/api/user-management/:quizId/:username", handleUserKeepAlive(s))

	// s.Echo.GET("/", handlePage)
}

// func handlePage(c echo.Context) error {
// 	uiBuildDir := echo.MustSubFS(statics, "Real-Time-Vocabulary-Quiz/build")
// 	index, err := uiBuildDir.Open("index.html")
// 	if err != nil {
// 		return c.JSON(http.StatusNotFound, "index.html not found")
// 	}

// 	pusher, ok := c.Response().Writer.(http.Pusher)
// 	if ok {
// 		err = fs.WalkDir(uiBuildDir, ".", func(path string, d fs.DirEntry, err error) error {
// 			if err != nil {
// 				return err
// 			}
// 			if d.IsDir() {
// 				return nil
// 			}

// 			if d.Name() == "manifest.json" ||
// 				d.Name() == "favicon.ico" ||
// 				strings.HasPrefix(d.Name(), "main") && (strings.HasSuffix(d.Name(), ".js") ||
// 					strings.HasSuffix(d.Name(), ".map") ||
// 					strings.HasSuffix(d.Name(), ".css")) {
// 				if err := pusher.Push(path, nil); err != nil {
// 					log.Err(err)
// 				}
// 			}

// 			return nil
// 		})
// 		if err != nil {
// 			log.Err(err)
// 		}
// 	}

// 	filename := c.Param("filename")
// 	file, err := uiBuildDir.Open(filename)
// 	if err != nil {
// 		log.Debug().Interface("file not found", filename)
// 		return ServeFile(c, index)
// 	}

// 	return ServeFile(c, file)
// }

// func ServeFile(c echo.Context, f fs.File) error {
// 	fi, _ := f.Stat()
// 	ff, ok := f.(io.ReadSeeker)
// 	if !ok {
// 		return errors.New("file does not implement io.ReadSeeker")
// 	}

// 	http.ServeContent(c.Response(), c.Request(), fi.Name(), fi.ModTime(), ff)
// 	return nil
// }
