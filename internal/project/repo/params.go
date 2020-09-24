package repo

import (
	"path"
)

type Commit struct {
	Sha string `json:"sha"`
	Url string `json:"url"`
}

type Tag struct {
	Name       string `json:"name"`
	ZipBallUrl string `json:"zipball_url"`
	TarBallUrl string `json:"tarball_url"`
	Commit     Commit `json:"commit"`
	NodeID     string `json:"node_id"`
}

type TagsResponse []Tag

type Repository struct {
	ID       uint32 `json:"id"`
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	HtmlUrl  string `json:"html_url"`
}

func (r Repository) GetPackageName() string {
	return path.Join("github.com", r.FullName)
}

type RepoResponse []Repository
