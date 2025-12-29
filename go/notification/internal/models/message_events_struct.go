package models

type MessageEvent struct {
	EventType string      `bson:"event_type" json:"event_type"`
	Payload   interface{} `bson:"payload" json:"payload"`
	Timestamp int64       `bson:"timestamp" json:"timestamp"`
}

type UserMessageEvent struct {
	EventType string `bson:"event_type" json:"event_type"`
	Payload   User   `bson:"payload" json:"payload"`
	Timestamp int64  `bson:"timestamp" json:"timestamp"`
}
