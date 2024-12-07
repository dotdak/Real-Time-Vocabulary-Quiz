package services

import (
	"errors"
	"math/rand"

	"github.com/dotdak/Real-Time-Vocabulary-Quiz/utils"
	"github.com/google/uuid"
)

type GameService struct {
	QuestionBank map[string]Question
}

type Question struct {
	Id            string            `json:"id"`
	Question      string            `json:"question"`
	Answers       map[string]string `json:"answer"`
	correctAnswer string            `json:"-"`
}

func NewGameService() *GameService {
	return &GameService{
		QuestionBank: map[string]Question{},
	}
}

var questions = []string{
	"relating to a very formal social occasion where men wear a black bow tie, or to the clothes worn for this occasion",
	"existing, happening, or done outside, rather than inside a building",
	"having a surface or consisting of a substance that is perfectly regular and has no holes, lumps, or areas that rise or fall suddenly",
	"not even or smooth, often because of being in bad condition",
	"of or from a long time ago, having lasted for a very long time",
	"having only a short distance from the top to the bottom",
}

var answers = []map[string]string{{
	"A": "hat",
	"B": "suit",
	"C": "pants",
	"D": "black tie",
}, {
	"A": "exception",
	"B": "upfront",
	"C": "indoor",
	"D": "outdoor",
}, {
	"A": "ugly",
	"B": "rough",
	"C": "ground",
	"D": "smooth",
}, {
	"A": "flawless",
	"B": "seamless",
	"C": "hard",
	"D": "rough",
}, {
	"A": "modern",
	"B": "history",
	"C": "nowaday",
	"D": "ancient",
}, {
	"A": "deep",
	"B": "range",
	"C": "ground",
	"D": "shallow",
},
}

var correctAnswers = []string{
	"D", "D", "D", "D", "D", "D",
}

var visited = utils.NewSet[int]()

func (g *GameService) GenerateQuiz() Question {
	questionId := uuid.NewString()
	for _, ok := g.QuestionBank[questionId]; ok; {
		questionId = uuid.NewString()
	}

	// TODO: remove this
	if visited.Len() == len(questions) {
		visited = utils.NewSet[int]()
	}
	index := rand.Intn(len(questions))
	for visited.Has(index) {
		index = rand.Intn(6)
	}
	visited.Add(index)

	q := Question{
		Id:            questionId,
		Question:      questions[index],
		Answers:       answers[index],
		correctAnswer: correctAnswers[index],
	}

	g.QuestionBank[q.Id] = q
	return q
}
func (g *GameService) AnswerAndGetPoint(questionId string, answer string) int {
	question, ok := g.QuestionBank[questionId]
	if !ok {
		return 0
	}

	if question.correctAnswer != answer {
		return 0
	}

	return 1
}

func (g *GameService) CheckAnswer(questionId string, answer string) (bool, error) {
	question, ok := g.QuestionBank[questionId]
	if !ok {
		return false, errors.New("question does not exist")
	}

	return question.correctAnswer == answer, nil
}
