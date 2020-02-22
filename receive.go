package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chneau/limiter"

	"github.com/urfave/cli"
)

func receiveAction(c *cli.Context) error {
	ip := c.Args().First()
	if !strings.Contains(ip, ":") {
		ip += ":8888"
	}
	files := map[string]int64{}
	{
		resp, err := http.Get("http://" + ip + "/files")
		if err != nil {
			return err
		}
		err = json.NewDecoder(resp.Body).Decode(&files)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
	}
	limit := limiter.New(c.Int("concurrence"))
	for file := range files {
		file := file
		limit.Execute(func() {
			req, err := http.NewRequest("GET", "http://"+ip+"/files", strings.NewReader(file))
			if err != nil {
				log.Println(err)
				return
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Println(err)
				return
			}
			defer resp.Body.Close()
			_ = os.MkdirAll(filepath.Dir(file), 0755)
			f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0755)
			if err != nil {
				log.Println(err)
				return
			}
			defer f.Close()
			_, err = io.Copy(f, resp.Body)
			if err != nil {
				log.Println(err)
				return
			}
		})
	}
	limit.Wait()
	_, _ = (&http.Client{Timeout: time.Microsecond}).Get("http://" + ip + "/end")
	return nil
}
