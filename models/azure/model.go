package azure

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/soulteary/amazing-openai-api/internal/define"
)

type RequestConverter interface {
	Name() string
	Convert(req *http.Request, config *define.ModelConfig) (*http.Request, error)
}

type StripPrefixConverter struct {
	Prefix string
}

func (c *StripPrefixConverter) Name() string {
	return "StripPrefix"
}

func (c *StripPrefixConverter) Convert(req *http.Request, config *define.ModelConfig) (*http.Request, error) {
	req.Host = config.URL.Host
	req.URL.Scheme = config.URL.Scheme
	req.URL.Host = config.URL.Host

	deployName := config.Model
	if config.Alias != "" {
		deployName = config.Alias
	}
	req.URL.Path = path.Join(fmt.Sprintf("/openai/deployments/%s", deployName), strings.Replace(req.URL.Path, c.Prefix+"/", "/", 1))
	req.URL.RawPath = req.URL.EscapedPath()

	query := req.URL.Query()
	query.Add(HeaderAPIVer, config.Version)
	req.URL.RawQuery = query.Encode()
	return req, nil
}

func NewStripPrefixConverter(prefix string) *StripPrefixConverter {
	return &StripPrefixConverter{
		Prefix: prefix,
	}
}
