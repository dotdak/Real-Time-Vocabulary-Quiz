package services

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type SessionStatus string

var (
	WaitingStatus    SessionStatus = "waiting"
	StartingStatus   SessionStatus = "starting"
	InProgressStatus SessionStatus = "inProgress"
	EndedStatus      SessionStatus = "ended"
)

type QuizConfig struct {
	QuestionTimer time.Duration `json:"questionTimer"`
	LifeTime      time.Duration `json:"lifeTime"`
	TotalQuestion int           `json:"totalQuestion"`
}

type QuizSession struct {
	Id            string
	LifeTime      time.Duration
	Users         sync.Map `json:"-"`
	timer         int      `json:"-"`
	Status        SessionStatus
	Config        *QuizConfig
	QuestionQueue chan Question
	GameService   *GameService
}

func NewQuizSession(id string, config *QuizConfig, gameService *GameService) *QuizSession {
	quiz := &QuizSession{
		Id:          id,
		Users:       sync.Map{},
		timer:       int(config.LifeTime.Milliseconds()),
		Status:      WaitingStatus,
		Config:      config,
		GameService: gameService,
	}

	return quiz
}

func (r *QuizSession) NotifyUsers(message interface{}) {
	r.Users.Range(func(key, value any) bool {
		user := value.(*User)
		log.Info().Msgf("notifying user %s in quiz %s", user.Name, user.quiz.Id)
		user.nofication <- message
		return true
	})
}

func (r *QuizSession) GetLeaderBoard() map[string]int {
	leaderBoard := make(map[string]int)
	r.Users.Range(func(key, value any) bool {
		user := value.(*User)
		leaderBoard[user.Name] = user.TotalPoint
		return true
	})

	return leaderBoard
}
func (r *QuizSession) SetStatus(status SessionStatus) {
	if r == nil {
		return
	}

	r.Status = status
}

func (r *QuizSession) Join(user *User) {
	if r == nil {
		return
	}

	r.Users.Store(user.Name, user)
}

func (r *QuizSession) FindUser(username string) *User {
	if r == nil {
		return nil
	}

	value, ok := r.Users.Load(username)
	if !ok {
		return nil
	}

	return value.(*User)
}
