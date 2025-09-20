WITH rank_rows AS (
    SELECT 
        stv2.student_question_file_v2_id, 
        stv2.testcase_id, 
        q.question_id,
        stv2.score / NULLIF(t.score, 0) AS normalized_score,
        q.total_score,
        ROW_NUMBER() OVER (
            PARTITION BY stv2.student_question_file_v2_id, q.question_id, stv2.testcase_id 
            ORDER BY stv2.std_test_v2_id DESC
        ) AS row_no
    FROM senior_project.student_testcases_v2 stv2
    INNER JOIN senior_project.testcases t 
        ON stv2.testcase_id = t.testcase_id
    INNER JOIN senior_project.questions q 
        ON t.question_id = q.question_id
    WHERE stv2.student_question_file_v2_id = ?
      AND q.question_id = ?
      AND stv2.status != 'N'
)
SELECT 
    AVG(normalized_score) * total_score AS score
FROM rank_rows
WHERE row_no = 1
GROUP BY student_question_file_v2_id, question_id, total_score
;