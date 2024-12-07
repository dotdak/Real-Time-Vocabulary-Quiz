package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/dotdak/Real-Time-Vocabulary-Quiz/events"
	"github.com/dotdak/Real-Time-Vocabulary-Quiz/services"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	Subprotocols:      []string{"binary"},
	HandshakeTimeout:  30 * time.Second,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections for simplicity. In a production environment, you may want to implement a more secure policy.
		return true
	},
}

type NotificationPayload[T any] struct {
	Status    services.SessionStatus `json:"@"`
	QuizId    string                 `json:"quizId"`
	EventType events.EventType       `json:"eventType"`
	Data      T                      `json:"data"`
}

func GetQuizState(r *services.QuizSession) NotificationPayload[map[string]int] {
	leaderBoard := make(map[string]int)
	r.Users.Range(func(key, value any) bool {
		user := value.(*services.User)
		leaderBoard[user.Name] = user.TotalPoint
		return true
	})
	return NotificationPayload[map[string]int]{
		Status:    r.Status,
		QuizId:    r.Id,
		EventType: events.EventTypeLeaderBoard,
		Data:      leaderBoard,
	}
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

type SessionStatusRequest struct {
	QuizId string                 `param:"quizId"`
	Status services.SessionStatus `form:"status"`
}

func handleQuizStatus(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req SessionStatusRequest
		err := c.Bind(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		session := s.FindQuizSession(req.QuizId)
		if session == nil {
			return c.String(http.StatusBadRequest, "quiz not found")
		}

		session.Status = req.Status

		// room.NotifyUpdate(nil)

		return c.NoContent(http.StatusOK)
	}
}

type QuizSessionByIdRequest struct {
	QuizId string `param:"quizId"`
}

func handleQuizSessionById(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req QuizSessionByIdRequest
		err := c.Bind(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		quiz := s.FindQuizSession(req.QuizId)
		if quiz == nil {
			return c.String(http.StatusBadRequest, "quiz not found")
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": quiz.Status,
			"config": quiz.Config,
		})
	}
}

type QuizSessionStatusByIdRequest struct {
	QuizId string                 `param:"quizId"`
	Status services.SessionStatus `form:"status"`
}

func handleQuizSessionStatusById(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req QuizSessionStatusByIdRequest
		err := c.Bind(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		quiz := s.FindQuizSession(req.QuizId)
		if quiz == nil {
			return c.String(http.StatusBadRequest, "quiz not found")
		}

		if quiz.Status != services.InProgressStatus && req.Status == services.InProgressStatus {
			quiz.Status = services.StartingStatus

			s.notificationPublisher.Send(NotificationPayload[int]{
				EventType: events.EventTypeCommand,
				QuizId:    quiz.Id,
				Status:    quiz.Status,
				Data:      3,
			})

			go func() {
				time.Sleep(3 * time.Second)

				quiz.Status = services.InProgressStatus
				for count := 0; quiz.Status != services.EndedStatus && count < quiz.Config.TotalQuestion; count++ {
					question := s.GameService.GenerateQuiz()
					s.notificationPublisher.Send(NotificationPayload[services.Question]{
						EventType: events.EventTypeQuestion,
						QuizId:    quiz.Id,
						Status:    quiz.Status,
						Data:      question,
					})
					time.Sleep(quiz.Config.QuestionTimer)
				}

				quiz.Status = services.EndedStatus
				s.notificationPublisher.Send(NotificationPayload[struct{}]{
					EventType: events.EventTypeCommand,
					QuizId:    quiz.Id,
					Status:    quiz.Status,
				})
			}()
			return c.NoContent(http.StatusOK)
		} else if quiz.Status == services.InProgressStatus && req.Status == services.EndedStatus {
			quiz.Status = services.EndedStatus
			return c.NoContent(http.StatusOK)
		}

		return c.String(http.StatusBadRequest, "invalid status")
	}
}

type QuizSessionRequest struct {
	QuizId string              `param:"quizId"`
	Config services.QuizConfig `form:"config"`
}

func handleQuizSession(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req QuizSessionRequest
		err := c.Bind(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		quiz := s.FindQuizSession(req.QuizId)
		if quiz == nil {
			newQuizId := req.QuizId
			if newQuizId == "" {
				newQuizId = uuid.New().String()
			}
			quiz = s.CreateQuizSession(newQuizId, &req.Config)
			log.Info().Msgf("new session created: %v", req.Config)
		}

		return c.JSON(http.StatusOK, map[string]string{
			"id": quiz.Id,
		})
	}
}

type NewQuestionsRequest struct {
	QuizId string `param:"quizId"`
}

func handleNewQuestions(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req SessionStatusRequest
		err := c.Bind(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		quiz := s.FindQuizSession(req.QuizId)
		if quiz == nil {
			return c.String(http.StatusBadRequest, "quiz session does not exist")
		}

		return c.JSON(http.StatusCreated, s.GameService.GenerateQuiz())
	}
}

type UserStateRequest struct {
	Username string `param:"username"`
	QuizId   string `param:"quizId"`
}

func handleUserState(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req UserStateRequest
		if err := c.Bind(&req); err != nil {
			return c.String(http.StatusBadRequest, "missing parameter")
		}

		quiz := s.FindQuizSession(req.QuizId)
		if quiz == nil {
			return c.JSON(http.StatusBadRequest, "quiz not exist")
		}

		user := quiz.FindUser(req.Username)
		if user == nil {
			return c.JSON(http.StatusBadRequest, "user not exist")
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"name":   user.Name,
			"points": user.TotalPoint,
		})
	}
}

type LeaderBoardRequest struct {
	QuizId string `param:"quizId"`
}

func handleLeaderBoard(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req UserStateRequest
		if err := c.Bind(&req); err != nil {
			return c.String(http.StatusBadRequest, "missing parameter")
		}

		quiz := s.FindQuizSession(req.QuizId)
		if quiz == nil {
			return c.JSON(http.StatusBadRequest, "quiz not exist")
		}

		return c.JSON(http.StatusOK, quiz.GetLeaderBoard())
	}
}

func handleConnection(s *Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		quizId := c.Param("quizId")
		if quizId == "" {
			log.Info().Msg("quiz id is missing")
			return c.String(http.StatusBadRequest, "quiz id is missing")
		}

		if strings.Contains(quizId, "/") || strings.Contains(quizId, ":") {
			log.Info().Msg("quiz id contains prohibited characters")
			return c.String(http.StatusBadRequest, "quiz id contains prohibited characters")
		}

		quiz := s.FindQuizSession(quizId)
		if quiz == nil {
			log.Info().Msg("quiz session is not opened yet")
			return c.String(http.StatusBadRequest, "quiz session is not opened yet")
		}

		username := c.QueryParam("username")
		if username == "" {
			log.Info().Msg("username is missing")
			return c.String(http.StatusBadRequest, "username is missing")
		}

		if strings.Contains(username, ":") {
			log.Info().Msg("username contains prohibited characters")
			return c.String(http.StatusBadRequest, "username contains prohibited characters")
		}

		conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
		if err != nil {
			log.Err(err).Msg("can not create connection with the client")
			return c.String(http.StatusBadRequest, "can not create connection with the client")
		}
		defer conn.Close()

		user := quiz.FindUser(username)
		if user == nil {
			_ = services.NewUser(quiz, conn, username)
		} else {
			user.Conn = conn
		}

		err = conn.WriteJSON(GetQuizState(quiz))
		if err != nil {
			log.Err(err).Msg("can not send message to client")
			return c.String(http.StatusInternalServerError, "unexpected error")
		}

		s.PublishQuizState(quiz)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Warn().AnErr("error reading message", err)
				return c.JSON(http.StatusGone, "connection failed")
			}
			// var userEvent UserEvent
			// if err := conn.ReadJSON(&userEvent); err != nil {
			// 	log.Error().AnErr("unable to read message: ", err)
			// }

			log.Debug().Msg(string(msg))

			var userEvent UserEvent

			err = json.Unmarshal(msg, &userEvent)
			if err != nil {
				log.Err(err).Msg("unable to parse user event")
				continue
			}

			s.eventPublisher.Send(userEvent)
		}
	}
}
