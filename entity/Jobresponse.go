package entity

/*type response struct {
	QuestionId int    `json:"question_id"`
	Response   string `json:"response"`
}
*/
type JobResponses struct {
	JobId         int    `json:"job_id"`
	JobStatusCode int    `json:"status"`
	JobStatus     string `json:"status_code"`
}
