package poster

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var posterExtensions = "jpg"

type Fetcher struct {
	fsBasePath string
	baseUrl    string
}

func NewFetcher(fsBasePath string, baseUrl string) *Fetcher {
	return &Fetcher{
		fsBasePath: strings.TrimRight(fsBasePath, "/"),
		baseUrl:    strings.TrimRight(baseUrl, "/"),
	}
}

func (f *Fetcher) Download(srcPath string, mediaId string) error {
	if err := os.MkdirAll(f.fsBasePath, 0777); err != nil {
		return err
	}

	extension := path.Ext(srcPath)

	posterFilePath := fmt.Sprintf("%s/%s%s", f.fsBasePath, mediaId, extension)

	posterFile, err := os.Create(posterFilePath)

	if err != nil {
		return err
	}

	defer func(posterFile *os.File) {
		_ = posterFile.Close()
	}(posterFile)

	posterUrl := fmt.Sprintf("%s/%s", f.baseUrl, strings.TrimLeft(srcPath, "/"))

	resp, err := http.Get(posterUrl)

	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	_, err = io.Copy(posterFile, resp.Body)

	return err
}

func GetPosterPath(mediaId string) (string, error) {
	fsBasePath := strings.TrimRight(viper.GetString("media_providers.tmdb.poster_fs_base_path"), "/")

	files, err := filepath.Glob(fmt.Sprintf("%s/%s.%s", fsBasePath, mediaId, posterExtensions))

	if err != nil {
		return "", err
	}

	if len(files) < 1 {
		return "", errors.New("poster not found")
	}

	return files[0], nil
}
