package events

type EventType string

var (
	EventTypeAnswer      EventType = "answer"
	EventTypeLeaderBoard EventType = "leaderboard"
	EventTypeCommand     EventType = "command"
	EventTypeQuestion    EventType = "question"
)
