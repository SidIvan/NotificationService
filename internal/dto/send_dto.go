package dto

type SendInfo struct {
	Id    int64  `json:"id"`
	Phone int    `json:"phone"`
	Text  string `json:"text"`
}
