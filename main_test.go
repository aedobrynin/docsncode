package main

import (
	"docsncode/buildcache"
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
// 9. test link with absolute path to file that will have result file (?)
// 10. test link with relative path to file that will have result file (done)
// 11. test link with absolute path to file that won't have result file (?)
// 12. test link with relative path to file that won't have result file and it is placed inside the project dir
// 13. test link with relative path to file that won't have result file and it is placed outside the project dir
// 13. test link to a website (done)
// 14. test image with relative path inside project dir
// 15. test image with absolute path (?)
// 16. test image with relative path outside project idr
// 17. test image from website
// 18. test many files in project
// 19. test allowed extensions
// 20. tests for cache
// tests 1-8 should be duplicated for every programming language with different comments syntax
// tests for errors

type testCase struct {
	name          string
	expectedError error
}

func runTests(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pathToProjectRoot := filepath.Join("tests", tc.name, "project")
			pathToExpectedResultDir := filepath.Join("tests", tc.name, "expected_result")

			resultDir := t.TempDir()

			// TODO: поддержать кэш в тестах
			err := buildDocsncode(pathToProjectRoot, resultDir, buildcache.NewAlwaysEmptyBuildCache())

			require.Equal(t, err, tc.expectedError)

			err = compare.Dirs(pathToExpectedResultDir, resultDir)
			require.NoError(t, err)
		})
	}
}

func TestCStyleComments(t *testing.T) {
	testCases := []testCase{
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
	}

	runTests(t, testCases)
}

func TestLinks(t *testing.T) {
	testCases := []testCase{
		{
			name:          "links/link_with_rel_path_to_file_with_result_file",
			expectedError: nil,
		},
		{
			name:          "links/link_to_website",
			expectedError: nil,
		},
		// TODO: добавить тесты на абсолютные пути в ссылках
		// TODO: возможно вообще избавиться от возможности задавать абсолютные пути?
	}

	runTests(t, testCases)
}

func TestDiagrams(t *testing.T) {
	testCases := []testCase{
		{
			name:          "diagrams/graph",
			expectedError: nil,
		},
	}

	runTests(t, testCases)
}
