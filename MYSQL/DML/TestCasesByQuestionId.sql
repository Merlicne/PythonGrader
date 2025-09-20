SELECT 
		tc.testcase_id 
		, tc.question_id 
		, tc.testcase_title 
		, tc.testcase_input 
		, tc.testcase_output 
		, tc.score
		, tc.regex_match 
		, tc.created_at 
		, tc.updated_at 
	FROM testcases tc
	WHERE tc.question_id  = ?