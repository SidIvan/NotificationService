package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Filter struct {
	OpCode string `bson:"opCode"`
	Tag    string `bson:"tag"`
}

type Distribution struct {
	StartAt string `bson:"startAt"`
	Message string `bson:"message"`
	Filter  Filter `bson:"filter"`
	EndAt   string `bson:"endAt"`
}

type DistributionWithId struct {
	StartAt string             `bson:"startAt"`
	Message string             `bson:"message"`
	Filter  Filter             `bson:"filter"`
	EndAt   string             `bson:"endAt"`
	Id      primitive.ObjectID `bson:"_id" json:"id"`
}

func (d *Distribution) IsValid() bool {
	return IsOppCodeValid(d.Filter.OpCode) && IsDateTimeValid(d.StartAt) && IsDateTimeValid(d.EndAt) && IsDateTimeBefore(d.StartAt, d.EndAt)
}

func (d *Distribution) SholdSend() bool {
	startTime, _ := time.Parse(dateTimeFormat, d.StartAt)
	endTime, _ := time.Parse(dateTimeFormat, d.EndAt)
	currentTime := time.Now()
	return startTime.Before(currentTime) && endTime.After(currentTime)
}

func (d *DistributionWithId) SholdSend() bool {
	startTime, _ := time.Parse(dateTimeFormat, d.StartAt)
	endTime, _ := time.Parse(dateTimeFormat, d.EndAt)
	currentTime := time.Now()
	return startTime.Before(currentTime) && endTime.After(currentTime)
}

func (d *Distribution) WithId(hexId string) *DistributionWithId {
	id, _ := primitive.ObjectIDFromHex(hexId)
	dWI := DistributionWithId{
		StartAt: d.StartAt,
		Message: d.Message,
		Filter:  d.Filter,
		EndAt:   d.EndAt,
		Id:      id,
	}
	return &dWI
}

func (d *DistributionWithId) GetId() string {
	return d.Id.Hex()
}

// func (d *Distribution) FormMongoFilter() bson.D {
// 	return bson.D{
// 		{"opCode", d.Filter.OpCode},
// 		{"tag", d.Filter.Tag},
// 	}
// }
