package main

import (
	"testing"

	compare "github.com/kilianpaquier/compare/pkg"
	"github.com/stretchr/testify/require"
)

// TODO: tests
// 1. test file with code
// 2. test file with one line comment block
// 3. test file with one line comment block and code
// 4. test file with several one line comment blocks and code
// 5. test file with multiline comment block
// 6. test file with multiline comment block and code
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

func TestMain(t *testing.T) {
	for _, tc := range []struct {
		name                    string
		pathToProjectRoot       string
		pathToExpectedResultDir string
		expectedError           error
	}{{
		name:                    "c_style_comments/one_file_with_code",
		pathToProjectRoot:       "tests/c_style_comments/one_file_with_code/project",
		pathToExpectedResultDir: "tests/c_style_comments/one_file_with_code/expected_result",
		expectedError:           nil,
	}} {
		t.Run(tc.name, func(t *testing.T) {
			resultDir := t.TempDir()

			err := buildDocsncode(tc.pathToProjectRoot, resultDir)

			require.Equal(t, err, tc.expectedError)

			err = compare.Dirs(tc.pathToExpectedResultDir, resultDir)
			require.NoError(t, err)
		})
	}
}
