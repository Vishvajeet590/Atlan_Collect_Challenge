package entity

type response struct {
	QuestionId   int    `json:"question_id"`
	Response     string `json:"response"`
	ResponseType string `json:"response_type"`
}

type Response struct {
	ResponseId int        `json:"response_id"`
	Responses  []response `json:"responses"`
}
