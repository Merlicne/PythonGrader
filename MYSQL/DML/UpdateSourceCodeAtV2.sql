
UPDATE senior_project.student_question_files_v2
SET student_question_file_id = ?,
    user_id = ?,
    question_id = ?,
    sourcecode = ?,
    version = ?,
    score = ?,
    updated_at = NOW()
WHERE student_question_file_v2_id = ?
    