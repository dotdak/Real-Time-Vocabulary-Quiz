package controllers

import (
	"context"
	"sync"
	"time"

	"github.com/dotdak/Real-Time-Vocabulary-Quiz/events"
	"github.com/dotdak/Real-Time-Vocabulary-Quiz/repo"
	"github.com/dotdak/Real-Time-Vocabulary-Quiz/services"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type UserEvent struct {
	EventType  events.EventType `json:"eventType"`
	QuizId     string           `json:"quizId"`
	QuestionId string           `json:"questionId"`
	Username   string           `json:"username"`
	Answer     string           `json:"answer"`
}

type Server struct {
	ctx                   context.Context
	Echo                  *echo.Echo
	quizSessions          sync.Map
	GameService           *services.GameService
	eventPublisher        *repo.Producer
	notificationPublisher *repo.Producer
	eventQueue            chan UserEvent
	notificationQueue     chan NotificationPayload[interface{}]
}

func NewServer(ctx context.Context) *Server {
	s := &Server{
		Echo:                  echo.New(),
		ctx:                   ctx,
		quizSessions:          sync.Map{},
		GameService:           services.NewGameService(),
		eventPublisher:        repo.NewProducer("quizEvent"),
		notificationPublisher: repo.NewProducer("notificationChannel"),
		eventQueue:            make(chan UserEvent),
		notificationQueue:     make(chan NotificationPayload[interface{}]),
	}

	repo.NewConsumer("quizEvent", s.eventQueue)
	repo.NewConsumer("notificationChannel", s.notificationQueue)
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case userEvent := <-s.eventQueue:
				quiz := s.FindQuizSession(userEvent.QuizId)
				if quiz == nil {
					return
				}
				user := quiz.FindUser(userEvent.Username)
				if user == nil {
					return
				}
				switch userEvent.EventType {
				case events.EventTypeAnswer:
					point := s.GameService.AnswerAndGetPoint(userEvent.QuestionId, userEvent.Answer)
					user.TotalPoint += point

					s.PublishQuizState(quiz)
				}
			}
		}
	}(ctx)

	go func(ctx context.Context) {
		for {
			select {
			case notification := <-s.notificationQueue:
				log.Info().Msgf("get notification: %v", notification)
				quiz := s.FindQuizSession(notification.QuizId)
				if quiz == nil {
					continue
				}
				quiz.NotifyUsers(notification)
				// if err := conn.WriteJSON(notification); err != nil {
				// 	log.Err(err).Msgf("can not send event %s", events.EventTypeQuestion)
				// }
				time.Sleep(time.Second)
			case <-ctx.Done():
				log.Error().AnErr("request context: ", ctx.Err())
				return
			}
		}
	}(ctx)

	return s
}

func (s* Server) PublishQuizState(quiz *services.QuizSession) {
	s.notificationPublisher.Send(NotificationPayload[map[string]int]{
		EventType: events.EventTypeLeaderBoard,
		QuizId:    quiz.Id,
		Status:    quiz.Status,
		Data:      quiz.GetLeaderBoard(),
	})
}

func (s *Server) FindQuizSession(quizId string) *services.QuizSession {
	if quizId == "" {
		return nil
	}

	rawQuizSession, ok := s.quizSessions.Load(quizId)
	if !ok {
		return nil
	}

	quiz, ok := rawQuizSession.(*services.QuizSession)
	if !ok {
		return nil
	}

	return quiz
}

func (s *Server) CreateQuizSession(quizId string, config *services.QuizConfig) *services.QuizSession {
	quiz := s.FindQuizSession(quizId)
	if quiz != nil {
		return quiz
	}

	quiz = services.NewQuizSession(quizId, config, s.GameService)
	s.quizSessions.Store(quizId, quiz)
	return quiz
}

func (s *Server) CleanQuizSession(quizId string) {
	s.quizSessions.Delete(quizId)
}
