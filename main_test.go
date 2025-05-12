package main

import (
	"docsncode/buildcache"
	"docsncode/pathsignorer"
	"os"
	"path/filepath"
	"testing"

	compare "github.com/kilianpaquier/compare/pkg"
	"github.com/stretchr/testify/require"
)

// TODO: tests
// * test file with several one line comment blocks and code
// * test file with multiline comment block
// * test file with several multiline line comment blocks and code
// * test file with several one line comment blocks, several multiline comment blocks and code
// * test many files in project
// * test allowed extensions
// * tests for cache
// * tests for .docsncodeignore
// * tests for errors

type testCase struct {
	name                        string
	expectedError               error
	createResultDirInTestFolder bool
}

func runTests(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pathToProjectRoot := filepath.Join("tests", tc.name, "project")
			pathToExpectedResultDir := filepath.Join("tests", tc.name, "expected_result")

			var resultDir string
			if tc.createResultDirInTestFolder {
				resultDir = filepath.Join("tests", tc.name, "actual_result")
				os.Mkdir(resultDir, 0755)
				defer os.RemoveAll(resultDir)
			} else {
				resultDir = t.TempDir()
			}

			// TODO: поддержать кэш в тестах
			err := buildDocsncode(pathToProjectRoot, resultDir, buildcache.NewAlwaysEmptyBuildCache(), pathsignorer.NewAlwaysNotIgnoringPathsIgnorer())

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
			name:          "c_style_comments/file_with_multiline_comment_block_and_code",
			expectedError: nil,
		},
	}

	runTests(t, testCases)
}

func TestPythonStyleComments(t *testing.T) {
	testCases := []testCase{
		{
			name:          "c_style_comments/file_with_code",
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
			name:                        "links/link_with_rel_path_to_file_in_project_dir_without_result_file",
			expectedError:               nil,
			createResultDirInTestFolder: true,
		},
		{
			name:                        "links/link_with_rel_path_to_file_outside_project_dir_without_result_file",
			expectedError:               nil,
			createResultDirInTestFolder: true,
		},
		{
			name:          "links/link_to_website",
			expectedError: nil,
		},
	}

	runTests(t, testCases)
}

func TestImages(t *testing.T) {
	testCases := []testCase{
		{
			name:                        "images/image_in_project_dir",
			expectedError:               nil,
			createResultDirInTestFolder: true,
		},
		{
			name:                        "images/image_outside_project_dir",
			expectedError:               nil,
			createResultDirInTestFolder: true,
		},
		{
			name:          "images/image_from_website",
			expectedError: nil,
		},
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

func TestCodeBlocks(t *testing.T) {
	testCases := []testCase{
		{
			name:          "code_blocks/empty_code_blocks",
			expectedError: nil,
		},
	}

	runTests(t, testCases)
}
