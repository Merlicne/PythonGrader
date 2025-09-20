package model

import "time"

type SourceCode struct {
	StudentQuestionFileV2Id int       `json:"student_question_file_v2_id" db:"student_question_file_v2_id"`
	StudentQuestionFileId   int       `json:"student_question_file_id" db:"student_question_file_id"`
	UserId                  int       `json:"user_id" db:"user_id"`
	QuestionId              int       `json:"question_id" db:"question_id"`
	SourceCode              string    `json:"sourcecode" db:"sourcecode"`
	Version                 int       `json:"version" db:"version"`
	Score                   float32   `json:"score" db:"score"`
	Status                  string    `json:"status" db:"status"`
	CreatedAt               time.Time `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time `json:"updated_at" db:"updated_at"`
}
