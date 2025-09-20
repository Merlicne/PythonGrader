package mysql


import (
	_ "embed"
)

//go:embed DML/source_code_info.sql
var SourceCodeInfo string

//go:embed DML/TestCasesByQuestionId.sql
var TestCasesByQuestionId string

//go:embed DML/InsertSourceCodeAtV2.sql
var InsertSourceCodeAtV2 string

//go:embed DML/UpdateSourceCodeAtV2.sql
var UpdateSourceCodeAtV2 string

//go:embed DML/InsertTestRunResultV2.sql
var InsertTestRunResultV2 string

//go:embed DML/GetSourceCodeInfoV2FromOldIdAndVersion.sql
var GetSourceCodeInfoV2FromOldIdAndVersion string

//go:embed DML/CalculateSourceCodeScore.sql
var CalculateSourceCodeScoreV2 string