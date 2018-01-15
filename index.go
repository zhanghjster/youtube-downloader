package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

type Index struct {
	Playlist string
	PageData struct {
		file      string
		PageToken string
	}
	Dir string
}

func NewIndex(playlist, dir string) (*Index, error) {
	err := createDirIfNotExist(dir + "/" + playlist)
	if err != nil {
		return nil, err
	}

	var index = &Index{
		Playlist: playlist,
		Dir:      dir,
	}

	index.PageData.file = dir + "/" + playlist + ".json"

	f, err := os.OpenFile(index.PageData.file, os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()

	if err != nil {
		return nil, err
	}

	json.NewDecoder(f).Decode(&index.PageData)

	return index, nil
}

func (i *Index) videoFlag(videoId string) string {
	return i.Dir + "/" + i.Playlist + "/" + videoId

}
func (i *Index) VideoIsDownloaded(videoId string) bool {
	if _, err := os.Stat(i.videoFlag(videoId)); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
func (i *Index) SetVideoDownloaded(videoId string) error {
	_, err := os.Create(i.videoFlag(videoId))
	return err
}

func (i *Index) UpdatePageToken(token string) error {
	i.PageData.PageToken = token
	return i.FlushPageData()
}

func (i *Index) FlushPageData() error {
	b, err := json.Marshal(i.PageData)
	if err != nil {
		log.Fatal(err)
	}
	return ioutil.WriteFile(i.PageData.file, b, 0755)
}
