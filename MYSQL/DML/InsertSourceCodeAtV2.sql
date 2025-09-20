INSERT INTO student_question_files_v2 (
    student_question_file_id,
    user_id,
    question_id,
    sourcecode,
    version,
    score,
    created_at,
    updated_at
)
SELECT ?, ?, ?, ?, ?, ?, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (
    SELECT 1 FROM student_question_files_v2
    WHERE student_question_file_id = ? AND version = ?
);
