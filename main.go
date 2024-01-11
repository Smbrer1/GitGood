package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

var helpText = strings.TrimSpace(`
Nap is a code snippet manager for your terminal.
github.com/Smbrer1/GitGood

Usage:
  ggd           - for interactive mode
  ggd list      - list all snippets
	ggd repo		  - change cwd to repo path

Create:
  ggd < main.go                 - save snippet from stdin
  nap example/main.go < main.go - save snippet with name`)

func main() {
	runCLI(os.Args[1:])
}

func runCLI(args []string) {
	config := readConfig()
	repos := readRepos(config)
	_ = config
	if len(args) > 0 {
		switch args[0] {
		case "list":
			listRepos(repos)
		case "-h", "--help":
			fmt.Println(helpText)
		case "add":
			if len(args) == 1 {
				fmt.Println("Please add path to repo")
				return
			}
			saveRepo(args[1], config, repos)
		default:
			repoArgs := args[1:]
			repoName := args[0]
			for _, repo := range repos {
				if repoName == repo.Name {
					if len(repoArgs) == 0 {
						fmt.Println("trying to cd")
						err := syscall.Chdir(repo.Folder)
						// move := exec.Command("cd", repo.Folder)
						// err := move.Run()
						if err != nil {
							fmt.Println("error", err)
							return
						}

						return
					}
					repoActions(repoName, repoArgs)
					return
				}
			}
			fmt.Println("You didn't save a repo with that name")
		}
		return
	}
}

func repoActions(repo string, args []string) {
	fmt.Println("Repo Action")
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

func listRepos(repos []Repo) {
	for _, repo := range repos {
		fmt.Printf("Path: %s\nName: %s, Tags: %v, Favourite: %v\n", repo.Folder, repo.Name, repo.Tags, repo.Favourite)
	}
}

func saveRepo(repoPath string, config Config, repos []Repo) {
	fullPath, name, err := parseRepo(repoPath)
	if err != nil {
		fmt.Printf("Unable to save repo, %+v", err)
		return
	}

	for _, repo := range repos {
		if fullPath == repo.Folder {
			fmt.Println("You already saved that repo")
			return
		}
	}

	// Add snippet metadata
	repo := Repo{
		Folder: fullPath,
		Name:   name,
	}

	repos = append([]Repo{repo}, repos...)
	writeRepo(config, repos)
}

func parseRepo(path string) (string, string, error) {
	switch path {
	case ".":
		repoPath, err := os.Getwd()
		if err != nil {
			return "", "", err
		}
		ok, err := checkIfRepo(repoPath)
		if err != nil {
			return "", "", err
		}
		sPath := strings.Split(repoPath, "/")
		if ok {
			return repoPath, sPath[len(sPath)-1], nil
		} else {
			return "", "", errors.New("not a git repo")
		}
	default:
		return "haha", "hehe", nil
	}
}

func checkIfRepo(path string) (bool, error) {
	gitCommmand := exec.Command("git", "-C", path, "rev-parse")

	_, err := gitCommmand.Output()
	if err != nil {
		print(err.Error())
		if err.Error() == "exit status 128" {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func writeRepo(config Config, repos []Repo) {
	b, err := json.Marshal(repos)
	if err != nil {
		fmt.Println("Could not marshal latest repo data.", err)
		return
	}

	if err := os.WriteFile(filepath.Join(config.Home, config.File), b, os.ModePerm); err != nil {
		fmt.Println("Could not save repo file.", err)
	}
}
