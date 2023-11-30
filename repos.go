package main

import "path/filepath"

type Repo struct {
	Folder    string
	Name      string
	Tags      []string
	Favourite bool
}

// Path returns the path <folder>/<repoName>
func (r Repo) Path() string {
	return filepath.Join(r.Folder, r.Name)
}

type Repos struct {
	repos []Repo
}
