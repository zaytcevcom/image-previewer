package app

import (
	"context"
	"fmt"
	"github.com/zaytcevcom/image-previewer/internal/cacher"
	"github.com/zaytcevcom/image-previewer/internal/utils"
	"net/http"
	"strconv"
)

type App struct {
	logger  Logger
	fetcher Fetcher
	cache   Cache
	resizer Resizer
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Fetcher interface {
	Fetch(ctx context.Context, header http.Header, url string) ([]byte, error)
}

type Cache interface {
	Get(key cacher.Key) (cacher.Value, bool)
	Set(key cacher.Key, value cacher.Value) bool
	Clear()
}

type Resizer interface {
	Fill(bytes []byte, width uint, height uint) ([]byte, error)
}

func New(logger Logger, fetcher Fetcher, cache Cache, resizer Resizer) *App {
	return &App{
		logger:  logger,
		fetcher: fetcher,
		cache:   cache,
		resizer: resizer,
	}
}

func (a *App) Fill(ctx context.Context, header http.Header, url string, width uint, height uint) ([]byte, error) {

	key := cacher.Key(
		utils.GetMD5Hash(strconv.Itoa(int(width)) + strconv.Itoa(int(height)) + url),
	)

	bytesValue, ok := a.cache.Get(key)
	if ok {
		bytes, ok := bytesValue.([]byte)
		if !ok {
			return nil, fmt.Errorf("failed convert to bytes")
		}

		return bytes, nil
	}

	bytes, err := a.fetcher.Fetch(ctx, header, url)
	if err != nil {
		return nil, err
	}

	bytes, err = a.resizer.Fill(bytes, width, height)
	if err != nil {
		return nil, err
	}

	a.cache.Set(key, bytes)

	return bytes, nil
}
