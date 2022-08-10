package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"

	"github.com/urfave/cli/v2"
)

func receiveAction(c *cli.Context) error {
	ip := c.Args().First()
	if !strings.Contains(ip, ":") {
		ip += ":8888"
	}
	files := map[string]int64{}
	{
		resp := lo.Must(http.Get("http://" + ip + "/files"))
		lo.Must0(json.NewDecoder(resp.Body).Decode(&files))
		defer resp.Body.Close()
	}
	lop.ForEach(lo.Keys(files), func(file string, _ int) {
		req := lo.Must(http.NewRequest("GET", "http://"+ip+"/files", strings.NewReader(file)))
		resp := lo.Must(http.DefaultClient.Do(req))
		defer resp.Body.Close()
		lo.Must0(os.MkdirAll(filepath.Dir(file), 0755))
		f := lo.Must(os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0755))
		defer f.Close()
		lo.Must(io.Copy(f, resp.Body))
	})
	lo.Must(http.Get("http://" + ip + "/end"))
	return nil
}
