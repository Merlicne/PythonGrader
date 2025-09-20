package model
// INSERT INTO senior_project.student_testcases_v2
// (std_test_v2_id, student_question_file_v2_id, testcase_id, score, status, test_output_text, checked_user_id, checked_at, created_at, updated_at)
// VALUES(0, 0, 0, 0, 'N', '', 0, '', '', '');

import "time"

type TestcaseResult struct {
	TestcaseResultId       int       `json:"testcase_result_id" db:"std_test_v2_id"`
	StudentQuestionFileV2Id int      `json:"student_question_file_v2_id" db:"student_question_file_v2_id"`
	TestcaseId             int       `json:"testcase_id" db:"testcase_id"`
	Score                  int       `json:"score" db:"score"`
	Status                 string    `json:"status" db:"status"`
	TestOutputText        string    `json:"test_output_text" db:"test_output_text"`
	CheckedUserId          int       `json:"checked_user_id" db:"checked_user_id"`
	CheckedAt              time.Time `json:"checked_at" db:"checked_at"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
}

