package entity

type Question struct {
	QuestionId   int    `json:"question_id"`
	Question     string `json:"question" binding:"required"`
	QuestionType string `json:"question_type" binding:"required"`
}

func NewQuestion(question, questionType string) Question {
	return Question{
		Question:     question,
		QuestionType: questionType,
	}
}

/*
{
"form_name": "Payment",
"owner_id": 450,
"question": [
{
"question": "What is your querry ?",
"question_type": "text/plain"
},
{
"question": "What do you want to do with it",
"question_type": "text/plain"
},
{
"question": "What is max ammount you need",
"question_type": "text/plain"
}
]
}
*/
