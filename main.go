package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func readGitignore(dir string) ([]string, error) {
	var ignorePatterns []string

	gitignorePath := filepath.Join(dir, ".gitignore")
	file, err := os.Open(gitignorePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		ignorePatterns = append(ignorePatterns, line)
	}
	return ignorePatterns, scanner.Err()
}

func shouldIgnore(path string, patterns []string) bool {
	for _, pattern := range patterns {
		matchBase, _ := filepath.Match(pattern, filepath.Base(path))
		matchPath, _ := filepath.Match(pattern, path)
		if matchBase || matchPath {
			return true
		}
	}
	return false
}

func isBinaryFile(filename string) bool {
	data, err := os.ReadFile(filename)
	if err != nil {
		return false
	}
	return bytes.Contains(data, []byte{0x00})
}

func printTreeStructure(path string, ignoredPatterns []string, prefix string, isLast bool, isRoot bool) {
	relPath, _ := filepath.Rel(".", path)

	// Check if this path should be ignored before proceeding further
	if shouldIgnore(relPath, ignoredPatterns) {
		return
	}

	info, err := os.Stat(path)
	if err != nil {
		return
	}

	// Adjust the prefix for current item based on whether it's the last in its directory
	currentPrefix := "├── "
	if isLast {
		currentPrefix = "└── "
	}

	if !isRoot { // Don't print the root itself
		fmt.Println(prefix + currentPrefix + info.Name())
	}

	if info.IsDir() {
		files, _ := os.ReadDir(path)
		for i, file := range files {
			newPrefix := prefix
			if isRoot || isLast {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			printTreeStructure(filepath.Join(path, file.Name()), ignoredPatterns, newPrefix, i == len(files)-1, false)
		}
	}
}

func printFileContents(path string, ignoredPatterns []string, addLineBreaks bool) {
	info, err := os.Stat(path)
	if err != nil || shouldIgnore(path, ignoredPatterns) || isBinaryFile(path) || info.IsDir() {
		return
	}

	if addLineBreaks {
		fmt.Println("\n\n>>>> " + path + " <<<<")
		fmt.Println("```")
	}
	file, _ := os.Open(path)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	fmt.Println("```")

}

func printAllFileContents(path string, ignoredPatterns []string) {
	relPath, _ := filepath.Rel(".", path)
	if shouldIgnore(relPath, ignoredPatterns) {
		return
	}

	info, err := os.Stat(path)
	if err != nil {
		return
	}

	if !info.IsDir() {
		printFileContents(path, ignoredPatterns, true)
	} else {
		files, _ := os.ReadDir(path)
		for _, file := range files {
			printAllFileContents(filepath.Join(path, file.Name()), ignoredPatterns)
		}
	}
}
func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for _, v := range elements {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}
	return result
}

func printHelp() {
	fmt.Println(`
Usage: proj-to-prompt [OPTIONS]

This tool provides a tree view of your project directory, including the content of non-binary files. It also respects .gitignore rules and provides additional ignore functionalities.
If a file named "project-requirements.md" exists as a sibling of the current directory, its content will be printed first.

OPTIONS:
    -h, --help 
        Show this help message and exit.

    --ignore PATTERN1,PATTERN2,...
        A comma-separated list of files or patterns to ignore. This is in addition to any patterns found in .gitignore.

Examples:

    # Print the tree of the current directory
    $ proj-to-prompt

    # Print the tree but ignore .md and .log files
    $ proj-to-prompt --ignore *.md,*.log

Note:
    - Binary files and files/directories matching the patterns in .gitignore or specified by --ignore will be excluded.
    - Always ensure you're in the right directory before running this script.

Visit https://github.com/ajduberstein/proj-to-prompt for more information and updates.`)
}

func main() {
	helpFlag := flag.Bool("help", false, "Show help message")
	ignoreFlag := flag.String("ignore", "", "Comma-separated list of files or patterns to ignore")
	flag.BoolVar(helpFlag, "h", false, "Show help message")
	flag.Parse()

	if *helpFlag {
		printHelp()
		return
	}

	ignoredPatterns := strings.Split(*ignoreFlag, ",")

	if len(ignoredPatterns) == 1 && ignoredPatterns[0] == "" {
		ignoredPatterns = []string{}
	}

	gitignorePatterns, _ := readGitignore(".")
	ignoredPatterns = append(ignoredPatterns, gitignorePatterns...)
	ignoredPatterns = append(ignoredPatterns, ".git/*", ".gitignore", "project-requirements.md")
	ignoredPatterns = removeDuplicates(ignoredPatterns)

	siblingFile := "project-requirements.md"

	if _, err := os.Stat(siblingFile); err == nil {
		fmt.Println("# Project requirements")
		printFileContents(siblingFile, []string{}, false)
	}

	// Print the tree structure
	fmt.Println(".")
	printTreeStructure(".", ignoredPatterns, "", false, true)

	// Print the contents of each file
	printAllFileContents(".", ignoredPatterns)

}
