package dto

type MessageStatus string

const (
	OkStatus    MessageStatus = "success"
	NotOkStatus MessageStatus = "failed"
)

type Message struct {
	SendedAt        string        `bson:"sendedAt"`
	Status          MessageStatus `bson:"status"`
	DisctributionId string        `bson:"distributionId"`
	ClientId        string        `bson:"clientId"`
}

func (m *Message) IsValid() bool {
	return IsDateTimeValid(m.SendedAt)
}
