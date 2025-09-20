package service

import (
	"context"
	"fmt"
	"os"
	"python-runner/executer"
	"python-runner/model"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func GradeFileByOldId(ctx context.Context, file string) error {
	if file == "" {
		return fmt.Errorf("--file must be provided")
	}

	// extract file name from path
	filename := strings.Split(file, "/")[len(strings.Split(file, "/"))-1]
	oldIdStr := strings.Split(strings.Split(filename, ".")[0], "_")[0]
	var versionId int
	var err error
	if len(strings.Split(strings.Split(filename, ".")[0], "_")) > 1 {
		versionId, err = strconv.Atoi(strings.Split(strings.Split(filename, ".")[0], "_")[1])
		if err != nil {
			return fmt.Errorf("failed to parse versionId from filename: %v", err.Error())
		}
	}

	var sourceCode string
	sourceCode, err = ReadSourceCodeFromFile(file)
	if err != nil {
		return fmt.Errorf("failed to read source code from file: %v", err.Error())
	}

	oldId, err := strconv.Atoi(oldIdStr)
	if err != nil {
		return fmt.Errorf("failed to parse oldId from filename: %v", err.Error())
	}
	return Grade(ctx, oldId, versionId, sourceCode)
}

func Grade(ctx context.Context, oldId int, versionId int, sourceCode string) error {
	// Add overall timeout for the entire grading process
	gradeCtx, gradeCancel := context.WithTimeout(ctx, time.Minute*2)
	defer gradeCancel()

	mysqlExecuter := executer.NewMySQLExecuter()
	python := &executer.PythonExecutor{}

	// Add timeout for database operations
	dbCtx, dbCancel := context.WithTimeout(gradeCtx, time.Second*30)
	codeInfo, err := mysqlExecuter.GetSourceCodeInfoWithContext(dbCtx, oldId)
	dbCancel()
	if err != nil {
		return fmt.Errorf("failed to get source code info for old ID %d: %v", oldId, err.Error())
	}

	dbCtx2, dbCancel2 := context.WithTimeout(gradeCtx, time.Second*30)
	testcases, err := mysqlExecuter.GetTestCasesWithContext(dbCtx2, codeInfo.QuestionId)
	dbCancel2()
	if err != nil {
		return fmt.Errorf("failed to get test cases for question ID %d: %v", codeInfo.QuestionId, err.Error())
	}

	if versionId == 0 {
		versionId = codeInfo.Version
	}
	newSourceCodeInfo := model.SourceCode{
		StudentQuestionFileId: codeInfo.StudentQuestionFileId,
		UserId:                codeInfo.UserId,
		QuestionId:            codeInfo.QuestionId,
		SourceCode:            sourceCode,
		Version:               versionId,
		Score:                 0,
		Status:                "N",
	}
	newSourceCodeInfoId, err := mysqlExecuter.InsertSourceCodeAtV2(newSourceCodeInfo)
	if err != nil {
		return fmt.Errorf("failed to insert source code: %v", err.Error())
	}

	newSourceCodeInfo.StudentQuestionFileV2Id = newSourceCodeInfoId

	// Process test cases with timeout protection
	for _, tc := range testcases {
		// Check if main context is cancelled
		select {
		case <-gradeCtx.Done():
			return fmt.Errorf("grading cancelled for old ID %d: %v", oldId, gradeCtx.Err())
		default:
		}

		testResult := model.TestcaseResult{
			StudentQuestionFileV2Id: newSourceCodeInfoId,
			TestcaseId:              tc.TestcaseId,
			Score:                   int(tc.Score),
			Status:                  "N",
			TestOutputText:          "",
		}

		// Create separate timeout for each test case execution
		testCtx, testCancel := context.WithTimeout(gradeCtx, time.Second*10)
		output, err := python.Execute(testCtx, sourceCode, tc.TestcaseInput)
		testCancel() // Always cancel to free resources

		var match bool = false
		var similarity float32 = 0
		if err != nil && output == "" {
			testResult.TestOutputText = err.Error()
			match = false
			similarity = 0
		} else {
			testResult.TestOutputText = output
			output = strings.TrimSpace(output)
			expected := strings.TrimSpace(tc.TestcaseOutput)
			match, similarity = compareResult(output, expected)
		}
		if match {
			testResult.Status = "P"
		} else {
			testResult.Status = "F"
			testResult.Score = int(float32(testResult.Score) * similarity)
		}

		err = mysqlExecuter.InsertTestRunResultV2(testResult)
		if err != nil {
			fmt.Printf("Error inserting test result: testcase %d, ErrorMessage: %v\n", tc.TestcaseId, err.Error())
			continue
		}
	}

	finalScore, err := mysqlExecuter.CalculateSourceCodeScoreV2(newSourceCodeInfoId, codeInfo.QuestionId)
	if err != nil {
		return fmt.Errorf("failed to calculate final score: %v", err.Error())
	}
	// update sourceCode info v2 with final score
	newSourceCodeInfo.Score = finalScore
	err = mysqlExecuter.UpdateSourceCodeAtV2(newSourceCodeInfo)
	if err != nil {
		return fmt.Errorf("failed to update source code with final score: %v", err.Error())
	}
	return nil
}

func compareResult(got string, want string) (bool, float32) {

	// replace all \r\n with \n
	var similarity float32 = 0
	got = strings.ReplaceAll(got, "\r\n", "\n")
	want = strings.ReplaceAll(want, "\r\n", "\n")

	gotLines := strings.Split(strings.TrimSpace(got), "\n")
	wantLines := strings.Split(strings.TrimSpace(want), "\n")

	min_len := len(gotLines)
	if len(wantLines) < min_len {
		min_len = len(wantLines)
	}

	for i := 0; i < min_len; i++ {
		gotLine := normalizeLine(gotLines[i])
		wantLine := normalizeLine(wantLines[i])
		// compare
		if gotLine == wantLine {
			similarity += 1.0
		}
	}
	return similarity/float32(len(wantLines)) == 1.0, similarity / float32(len(wantLines))
}

func normalizeLine(s string) string {
	s = strings.ReplaceAll(s, ":", "")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	return s
}

func ReadSourceCodeFromFile(file string) (string, error) {
	if file != "" {
		sourceCodeBytes, err := os.ReadFile(file)
		if err != nil {
			return "", fmt.Errorf("failed to read file: %v", err.Error())
		}
		return string(sourceCodeBytes), nil
	} else {
		return "", fmt.Errorf("file must be provided")
	}
}

func GradeFilesFromIdsCSV(csvfile string, latestVersionDir string, olderVersionDir string) error {
	return GradeFilesFromIdsCSVWithWorkers(csvfile, latestVersionDir, olderVersionDir, 4) // Default to 4 workers
}

func GradeFilesFromIdsCSVWithWorkers(csvfile string, latestVersionDir string, olderVersionDir string, maxWorkers int) error {
	if csvfile == "" {
		return fmt.Errorf("--csvfile must be provided")
	}
	if latestVersionDir == "" {
		return fmt.Errorf("--latestVersionDir must be provided")
	}
	if olderVersionDir == "" {
		return fmt.Errorf("--olderVersionDir must be provided")
	}
	if maxWorkers <= 0 {
		maxWorkers = 4 // Default fallback
	}

	csvBytes, err := os.ReadFile(csvfile)
	if err != nil {
		return fmt.Errorf("failed to read csv file: %v", err.Error())
	}
	csvContent := string(csvBytes)
	lines := strings.Split(strings.TrimSpace(csvContent), "\n")

	// Parse all valid IDs first
	var validIds []int
	for _, id := range lines {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		oldId, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("Skipping invalid ID '%s': %v\n", id, err)
			continue
		}
		validIds = append(validIds, oldId)
	}

	fmt.Printf("Processing %d valid IDs with %d workers\n", len(validIds), maxWorkers)

	// Create job channels
	latestVersionJobs := make(chan int, len(validIds))
	olderVersionJobs := make(chan int, len(validIds))

	// WaitGroups to track completion
	var latestWg sync.WaitGroup
	var olderWg sync.WaitGroup

	// Progress counters
	var completedCount int64
	totalFiles := int64(len(validIds))

	// Start workers for latest version files
	for i := 0; i < maxWorkers; i++ {
		latestWg.Add(1)
		go func(workerID int) {
			defer latestWg.Done()
			for oldId := range latestVersionJobs {
				ctx, cancelFunc := context.WithTimeout(context.Background(), time.Minute*3) // Increase timeout
				processLatestVersionFile(ctx, oldId, latestVersionDir)
				cancelFunc()

				// Update progress
				completed := atomic.AddInt64(&completedCount, 1)
				if completed%100 == 0 {
					fmt.Printf("Progress: %d/%d files processed (%.1f%%)\n", completed, totalFiles, float64(completed)/float64(totalFiles)*100)
				}
			}
		}(i)
	}

	// Start workers for older version files
	for i := 0; i < maxWorkers; i++ {
		olderWg.Add(1)
		go func(workerID int) {
			defer olderWg.Done()
			for oldId := range olderVersionJobs {
				ctx, cancelFunc := context.WithTimeout(context.Background(), time.Minute*5) // Longer timeout for multiple files
				processOlderVersionFiles(ctx, oldId, olderVersionDir)
				cancelFunc()
			}
		}(i)
	}

	// Send jobs to workers
	for _, oldId := range validIds {
		latestVersionJobs <- oldId
		olderVersionJobs <- oldId
	}

	// Close job channels to signal no more work
	close(latestVersionJobs)
	close(olderVersionJobs)

	// Wait for all workers to complete
	latestWg.Wait()
	olderWg.Wait()

	fmt.Println("All processing completed!")
	return nil
}

// processLatestVersionFile searches for and processes the latest version file ("<id>.py") in the specified directory
func processLatestVersionFile(ctx context.Context, oldId int, latestVersionDir string) {
	latestVersionFile := fmt.Sprintf("%s/%d.py", latestVersionDir, oldId)
	if _, err := os.Stat(latestVersionFile); err == nil {
		err := GradeFileByOldId(ctx, latestVersionFile)
		if err != nil {
			fmt.Printf("Error grading latest version file %s: %v\n", latestVersionFile, err)
		}
	}
}

// processOlderVersionFiles searches for and processes older version files ("<id>_<version>.py") in the specified directory
func processOlderVersionFiles(ctx context.Context, oldId int, olderVersionDir string) {
	files, err := os.ReadDir(olderVersionDir)
	if err != nil {
		fmt.Printf("Error reading older version directory %s: %v\n", olderVersionDir, err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		// Check if file matches pattern "<id>_<version>.py"
		if strings.HasSuffix(filename, ".py") {
			filePrefix := strings.TrimSuffix(filename, ".py")
			parts := strings.Split(filePrefix, "_")
			if len(parts) == 2 {
				fileOldId, err := strconv.Atoi(parts[0])
				if err == nil && fileOldId == oldId {
					olderVersionFile := fmt.Sprintf("%s/%s", olderVersionDir, filename)
					err := GradeFileByOldId(ctx, olderVersionFile)
					if err != nil {
						fmt.Printf("Error grading older version file %s: %v\n", olderVersionFile, err)
					}
				}
			}
		}
	}
}
