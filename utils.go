package goster

import (
	"fmt"
	"mime"
	"path/filepath"
	"strings"
)

// Adds basic headers
func DefaultHeader(c *Ctx) {
	c.Response.Header().Set("Access-Control-Allow-Origin", "*")
	c.Response.Header().Set("Connection", "Keep-Alive")
	c.Response.Header().Set("Keep-Alive", "timeout=5, max=997")
}

// cleanPath sanatizes a URL path. It removes suffix '/' if any and adds prefix '/' if missing. If the URL contains Query Parameters or Anchors,
// they will be removed as well.
func cleanPath(path *string) {
	if len(*path) == 0 {
		return
	}

	if (*path)[0] != '/' {
		*path = "/" + *path
	}

	*path = strings.TrimSuffix(*path, "/")
}

func cleanEmptyBytes(b *[]byte) {
	cleaned := []byte{}

	for _, v := range *b {
		if v == 0 {
			break
		}
		cleaned = append(cleaned, v)
	}
	*b = cleaned
}

func matchDynPathValue(dynPath, url string) (dp []DynamicPath, err error) {
	dynPathSlice := strings.Split(dynPath, "/")
	urlSice := strings.Split(url, "/")

	if len(dynPathSlice) != len(urlSice) {
		err = fmt.Errorf("request URL path `%s` does not match with dynamic Route path `%s`", url, dynPath)
		return
	}

	dp = make([]DynamicPath, len(dynPathSlice))
	for i, path := range dynPathSlice {
		if strings.ContainsRune(path, ':') {
			dp = append(dp, DynamicPath{
				path:  strings.TrimPrefix(path, ":"),
				value: urlSice[i],
			})
		}
	}

	return
}

func getContentType(filename string) string {
	ext := filepath.Ext(filename)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return contentType
}
