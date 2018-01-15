package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"os/exec"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"

	"net"

	"net/http"

	"bufio"

	"sync/atomic"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
)

var (
	videoDir   string
	indexDir   string
	playlist   string
	secret     string
	sockProxy  string
	interval   int
	concurrent int32
)

var RootCmd *cobra.Command

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.00",
	})
	RootCmd = &cobra.Command{
		Use: "youtube-downloader",
		Run: cmdRunner,
	}

	RootCmd.Flags().StringVarP(&playlist, "playlist", "p", "", "playlist id ")
	RootCmd.Flags().StringVar(&videoDir, "video-dir", "video", "[video] Dir for downloaded video")
	RootCmd.Flags().StringVar(&indexDir, "index-dir", ".index", "[.index] Dir for index")
	RootCmd.Flags().StringVar(&secret, "secret", "client_secret.json", "secret file")
	RootCmd.Flags().StringVar(&sockProxy, "sock-proxy", "", "HOST:PORT socket proxy")
	RootCmd.Flags().IntVar(&interval, "interval", 10, "interval of playlist check")
	RootCmd.Flags().Int32Var(&concurrent, "concurrent", 1, "concurrency count")
}

func cmdRunner(cmd *cobra.Command, args []string) {
	if playlist == "" {
		cmd.Usage()
		os.Exit(0)
	}

	err := createDirIfNotExist(videoDir)
	fatalErr(err)

	idx, err := NewIndex(playlist, indexDir)
	fatalErr(err)

	b, err := ioutil.ReadFile(secret)
	fatalErr(err)

	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	fatalErr(err)

	ctx := context.Background()
	var client *http.Client

	if sockProxy != "" {
		dialer, err := proxy.SOCKS5(
			"tcp",
			sockProxy,
			nil,
			&net.Dialer{Timeout: 5 * time.Second, KeepAlive: 30 * time.Second},
		)
		fatalErr(err)

		http.DefaultClient.Transport = &http.Transport{
			Proxy: nil, Dial: dialer.Dial,
		}
	}

	client = getClient(ctx, config)

	service, err := youtube.New(client)
	fatalErr(err)

	call := service.PlaylistItems.List("snippet,contentDetails")
	call = call.MaxResults(10)
	call = call.PlaylistId(playlist)

	var ids = make(chan string)
	var curJobs int32
	go func() {
		for {
			videoId := <-ids
			for curJobs >= concurrent {
			}
			atomic.AddInt32(&curJobs, 1)

			go func() {
				defer atomic.AddInt32(&curJobs, -1)
				if idx.VideoIsDownloaded(videoId) {
					return
				}

				log.Info("start download ", videoId)

				var args = []string{
					"https://www.youtube.com/watch?v=" + videoId,
					fmt.Sprintf("-o%s", videoDir),
				}

				if sockProxy != "" {
					args = append(args, fmt.Sprintf("-s%s", sockProxy))
				}

				cmd := exec.Command("you-get", args...)
				cmdReader, err := cmd.StdoutPipe()
				if err != nil {
					log.Error(err)
					return
				}

				scanner := bufio.NewScanner(cmdReader)
				scanner.Split(bufio.ScanRunes)
				go func() {
					for scanner.Scan() {
						if concurrent < 2 {
							fmt.Print(scanner.Text())
						}
					}
				}()

				if err = cmd.Start(); err != nil {
					log.Error(err)
					return
				}

				cmd.Wait()

				if err := idx.SetVideoDownloaded(videoId); err != nil {
					log.Infof("save %s video flag fail\n", videoId)
				}
			}()
		}
	}()

	for {
		if idx.PageData.PageToken != "" {
			call = call.PageToken(idx.PageData.PageToken)
		}

		response, err := call.Do()
		if err != nil {
			log.Error(err)
		}

		for _, item := range response.Items {
			ids <- item.ContentDetails.VideoId
		}

		if response.NextPageToken != "" {
			idx.UpdatePageToken(response.NextPageToken)
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func main() {
	err := RootCmd.Execute()
	fatalErr(err)
}
