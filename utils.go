package goster

import (
	"fmt"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Adds basic headers
func DefaultHeader(c *Ctx) {
	c.Response.Header().Set("Access-Control-Allow-Origin", "*")
	c.Response.Header().Set("Connection", "Keep-Alive")
	c.Response.Header().Set("Keep-Alive", "timeout=5, max=997")
}

// cleanURLPath sanatizes a URL path. It removes suffix '/' if any and adds prefix '/' if missing. If the URL contains Query Parameters or Anchors,
// they will be removed as well.
func cleanURLPath(path *string) {
	if len(*path) == 0 {
		return
	}

	if (*path)[0] != '/' {
		*path = "/" + *path
	}

	*path = strings.TrimSuffix(*path, "/")
}

func matchDynamicPath(dynamicPath, url string) (dp []DynamicPath, isDynamic bool) {
	dynamicPathSlice := strings.Split(dynamicPath, "/")
	urlSlice := strings.Split(url, "/")

	if len(dynamicPathSlice) != len(urlSlice) {
		return nil, false
	}

	hasDynamic := false
	dp = []DynamicPath{}
	for i, seg := range dynamicPathSlice {
		if strings.HasPrefix(seg, ":") {
			hasDynamic = true
			dp = append(dp, DynamicPath{
				path:  strings.TrimPrefix(seg, ":"),
				value: urlSlice[i],
			})
		} else if seg != urlSlice[i] {
			// static segment doesn't match
			return nil, false
		}
	}

	return dp, hasDynamic
}

func getContentType(filename string) string {
	ext := filepath.Ext(filename)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return contentType
}

func resolveAppPath(dir string) (string, error) {
	fileDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot determine working directory for static dir %s\n", dir)
		return dir, err
	}

	// construct the full path to the static directory
	return path.Join(fileDir, dir), nil
}

func pathExists(path string) (exists bool) {
	_, err := os.Stat(path)
	return err == nil
}

// func cleanEmptyBytes(b *[]byte) {
// 	cleaned := []byte{}

// 	for _, v := range *b {
// 		if v == 0 {
// 			break
// 		}
// 		cleaned = append(cleaned, v)
// 	}
// 	*b = cleaned
// }
