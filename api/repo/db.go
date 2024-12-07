package repo

type UserRoomKey struct {
	UserId string
	QuizId string
}

var QuizUserDB = make(map[UserRoomKey]string)
var QuizSessionDB = make(map[string]string)
