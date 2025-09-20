package model

import "time"

type Testcase struct {
	TestcaseId    int     `json:"testcase_id" db:"testcase_id"`
	QuestionId    int  `json:"question_id" db:"question_id"`
	TestcaseTitle string  `json:"testcase_title" db:"testcase_title"`
	TestcaseInput string  `json:"testcase_input" db:"testcase_input"`
	TestcaseOutput string `json:"testcase_output" db:"testcase_output"`
	Score         float64 `json:"score" db:"score"`
	RegexMatch    string    `json:"regex_match" db:"regex_match"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}



// tc.testcase_id 
// 		, tc.question_id 
// 		, tc.testcase_title 
// 		, tc.testcase_input 
// 		, tc.testcase_output 
// 		, tc.score
// 		, tc.regex_match 
// 		, tc.created_at 
// 		, tc.updated_at