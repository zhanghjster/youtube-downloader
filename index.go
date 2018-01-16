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
	if err := createDirIfNotExist(dir + "/" + playlist); err != nil {
		return nil, err
	}

	var index = &Index{
		Playlist: playlist,
		Dir:      dir,
	}

	index.PageData.file = dir + "/" + playlist + ".json"

	f, err := os.OpenFile(
		index.PageData.file, os.O_RDWR|os.O_CREATE, 0755,
	)
	if err != nil {
		return nil, err
	}

	json.NewDecoder(f).Decode(&index.PageData)
	defer f.Close()

	return index, nil
}

func (i *Index) videoFlagFile(videoId string) string {
	return i.Dir + "/" + i.Playlist + "/" + videoId

}
func (i *Index) VideoIsDownloaded(videoId string) bool {
	_, err := os.Stat(i.videoFlagFile(videoId))
	return !os.IsNotExist(err)
}
func (i *Index) SetVideoDownloaded(videoId string) error {
	_, err := os.Create(i.videoFlagFile(videoId))
	return err
}

func (i *Index) UpdatePageToken(token string) error {
	i.PageData.PageToken = token
	return i.FlushPageData()
}

func (i *Index) FlushPageData() error {
	b, err := json.Marshal(i.PageData)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(i.PageData.file, b, 0755)
}
