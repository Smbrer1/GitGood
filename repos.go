package main

import "path/filepath"

type Repo struct {
	Folder    string   `json:"folder"`
	Name      string   `json:"name"`
	Tags      []string `json:"tags"`
	Favourite bool     `json:"favourite"`
}

// Path returns the path <folder>/<repoName>
func (r Repo) Path() string {
	return filepath.Join(r.Folder, r.Name)
}

type Repos struct {
	repos []Repo
}
