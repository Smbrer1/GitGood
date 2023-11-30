package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var helpText = strings.TrimSpace(`
Nap is a code snippet manager for your terminal.
github.com/Smbrer1/GitGood

Usage:
  nap           - for interactive mode
  nap list      - list all snippets
  nap <snippet> - print snippet to stdout

Create:
  nap < main.go                 - save snippet from stdin
  nap example/main.go < main.go - save snippet with name`)

func main() {
	runCLI(os.Args[1:])
}

func runCLI(args []string) {
	config := readConfig()
	repos := readRepos(config)
	_ = repos
	_ = config

	if len(args) > 0 {
		switch args[0] {
		case "list":
			listSnippets(snippets)
		case "-h", "--help":
			fmt.Println(helpText)
		default:
			snippet := findSnippet(args[0], snippets)
			fmt.Print(snippet.Content(isatty.IsTerminal(os.Stdout.Fd())))
		}
		return
	}
}

// readStdin returns the stdin that is piped in to the command line interface.
func readStdin() string {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return ""
	}

	if stat.Mode()&os.ModeCharDevice != 0 {
		return ""
	}

	reader := bufio.NewReader(os.Stdin)
	var b strings.Builder

	for {
		r, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		_, err = b.WriteRune(r)
		if err != nil {
			return ""
		}
	}

	return b.String()
}

func readRepos(config Config) []Repo {
	var repos []Repo
	file := filepath.Join(config.Home, config.File)
	dir, err := os.ReadFile(file)
	if err != nil {
		// File does not exist, create one.
		err := os.MkdirAll(config.Home, os.ModePerm)
		if err != nil {
			fmt.Printf("Unable to create directory %s, %+v", config.Home, err)
		}
		f, err := os.Create(file)
		if err != nil {
			fmt.Printf("Unable to create file %s, %+v", file, err)
		}
		defer f.Close()
		dir = []byte("[]")
		_, _ = f.Write(dir)
	}
	err = json.Unmarshal(dir, &repos)
	if err != nil {
		fmt.Printf("Unable to unmarshal %s file, %+v\n", file, err)
		return repos
	}
	return repos
}

func saveRepo(content string, args []string, config Config, repos []Repo) {
	// folder, name, language := parseName(name)
	// file := fmt.Sprintf("%s.%s", name, language)
	// filePath := filepath.Join(config.Home, folder, file)
	// if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
	// 	fmt.Println("unable to create folder")
	// 	return
	// }
	// err := os.WriteFile(filePath, []byte(content), 0o644)
	// if err != nil {
	// 	fmt.Println("unable to create snippet")
	// 	return
	// }
	//
	// // Add snippet metadata
	// snippet := Snippet{
	// 	Folder:   folder,
	// 	Date:     time.Now(),
	// 	Name:     name,
	// 	File:     file,
	// 	Language: language,
	// }
	//
	// snippets = append([]Snippet{snippet}, snippets...)
	// writeSnippets(config, snippets)
}
