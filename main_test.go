package main

import (
	"path/filepath"
	"testing"

	compare "github.com/kilianpaquier/compare/pkg"
	"github.com/stretchr/testify/require"
)

// TODO: tests
// 1. test file with code (done)
// 2. test file with one line comment block (done)
// 3. test file with one line comment block and code (done)
// 4. test file with several one line comment blocks and code
// 5. test file with multiline comment block
// 6. test file with multiline comment block and code (done)
// 7. test file with several multiline line comment blocks and code
// 8. test file with several one line comment blocks, several multiline comment blocks and code
// 9. test link with absolute path to file that will have result file
// 10. test link with relative path to file that will have result file
// 11. test link with absolute path to file that won't have result file
// 12. test link with relative path to file that won't have result file
// 13. test link to a website
// 14. test image with relative path to file that will have result file
// 15. test image with absolute path to file that won't have result file
// 16. test image with relative path to file that won't have result file
// 17. test image from website
// 18. test many files in project
// 19. test allowed extensions
// tests 1-8 should be duplicated for every programming language with different comments syntax
// tests for errors

func TestMain(t *testing.T) {
	for _, tc := range []struct {
		name          string
		expectedError error
	}{
		{
			name:          "c_style_comments/file_with_code",
			expectedError: nil,
		},
		{
			name:          "c_style_comments/file_with_single_line_comment_block",
			expectedError: nil,
		},
		{
			name:          "c_style_comments/file_with_single_line_comment_block_and_code",
			expectedError: nil,
		},
		{
			name:          "c_style_comments/file_with_single_line_comment_block_and_code",
			expectedError: nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pathToProjectRoot := filepath.Join("tests", tc.name, "project")
			pathToExpectedResultDir := filepath.Join("tests", tc.name, "expected_result")

			resultDir := t.TempDir()

			err := buildDocsncode(pathToProjectRoot, resultDir)

			require.Equal(t, err, tc.expectedError)

			err = compare.Dirs(pathToExpectedResultDir, resultDir)
			require.NoError(t, err)
		})
	}
}
