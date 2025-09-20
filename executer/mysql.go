package executer

import (
	"context"
	"python-runner/model"
	"time"

	mysqlLocal "python-runner/MYSQL"

	"github.com/jmoiron/sqlx"
)

type MySQLExecuter struct {
	conn *sqlx.DB
}

func NewMySQLExecuter() *MySQLExecuter {
	// Initialize MySQL connection
	if mysqlLocal.GlobalConnection == nil {
		err := mysqlLocal.InitializeGlobalConnection()
		if err != nil {
			panic("failed to initialize MySQL connection: " + err.Error())
		}
	}
	return &MySQLExecuter{conn: mysqlLocal.GlobalConnection.DB}
}

func (e *MySQLExecuter) GetSourceCodeInfo(sourceCodeId int) (model.SourceCode, error) {
	var sourceCode model.SourceCode
	query := mysqlLocal.SourceCodeInfo
	err := e.conn.Get(&sourceCode, query, sourceCodeId)
	if err != nil {
		return model.SourceCode{}, err
	}
	return sourceCode, nil
}

func (e *MySQLExecuter) GetSourceCodeInfoWithContext(ctx context.Context, sourceCodeId int) (model.SourceCode, error) {
	var sourceCode model.SourceCode
	query := mysqlLocal.SourceCodeInfo
	err := e.conn.GetContext(ctx, &sourceCode, query, sourceCodeId)
	if err != nil {
		return model.SourceCode{}, err
	}
	return sourceCode, nil
}

func (e *MySQLExecuter) GetTestCases(questionId int) ([]model.Testcase, error) {
	var testCases []model.Testcase
	query := mysqlLocal.TestCasesByQuestionId
	err := e.conn.Select(&testCases, query, questionId)
	if err != nil {
		return nil, err
	}
	return testCases, nil
}

func (e *MySQLExecuter) GetTestCasesWithContext(ctx context.Context, questionId int) ([]model.Testcase, error) {
	var testCases []model.Testcase
	query := mysqlLocal.TestCasesByQuestionId
	err := e.conn.SelectContext(ctx, &testCases, query, questionId)
	if err != nil {
		return nil, err
	}
	return testCases, nil
}

func (e *MySQLExecuter) InsertSourceCodeAtV2(newSourceCodeInfo model.SourceCode) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var newSourceCodeId int
	query := mysqlLocal.InsertSourceCodeAtV2
	err := e.conn.QueryRowContext(ctx, query,
		newSourceCodeInfo.StudentQuestionFileId,
		newSourceCodeInfo.UserId,
		newSourceCodeInfo.QuestionId,
		newSourceCodeInfo.SourceCode,
		newSourceCodeInfo.Version,
		newSourceCodeInfo.Score,
		newSourceCodeInfo.StudentQuestionFileId,
		newSourceCodeInfo.Version,
	).Err()
	if err != nil {
		return 0, err
	}
	query = mysqlLocal.GetSourceCodeInfoV2FromOldIdAndVersion
	err = e.conn.QueryRowContext(ctx, query, newSourceCodeInfo.StudentQuestionFileId, newSourceCodeInfo.Version).Scan(&newSourceCodeId)
	if err != nil {
		return 0, err
	}
	return newSourceCodeId, nil
}

func (e *MySQLExecuter) InsertTestRunResultV2(testResult model.TestcaseResult) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := mysqlLocal.InsertTestRunResultV2
	_, err := e.conn.ExecContext(ctx, query, testResult.StudentQuestionFileV2Id, testResult.TestcaseId, testResult.Score, testResult.Status, testResult.TestOutputText)
	return err
}

func (e *MySQLExecuter) CalculateSourceCodeScoreV2(studentQuestionFileV2Id int, questionId int) (float32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var score float32
	query := mysqlLocal.CalculateSourceCodeScoreV2
	err := e.conn.QueryRowContext(ctx, query, studentQuestionFileV2Id, questionId).Scan(&score)
	if err != nil {
		return 0, err
	}
	return score, nil
}

func (e *MySQLExecuter) UpdateSourceCodeAtV2(sourceCodeInfo model.SourceCode) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := mysqlLocal.UpdateSourceCodeAtV2
	_, err := e.conn.ExecContext(ctx, query,
		sourceCodeInfo.StudentQuestionFileId,
		sourceCodeInfo.UserId,
		sourceCodeInfo.QuestionId,
		sourceCodeInfo.SourceCode,
		sourceCodeInfo.Version,
		sourceCodeInfo.Score,
		sourceCodeInfo.StudentQuestionFileV2Id,
	)
	return err
}
