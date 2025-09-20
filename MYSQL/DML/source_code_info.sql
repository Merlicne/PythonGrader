SELECT 
    sqf.student_question_file_id,
    sqf.user_id,
    sqf.question_id,
    COALESCE(sqf.sourcecode, '') AS sourcecode,
    sqf.version,
    sqf.score,
    sqf.status,
    sqf.created_at,
    sqf.updated_at
    FROM student_question_files sqf
    WHERE sqf.student_question_file_id = ?
