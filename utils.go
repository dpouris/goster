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

// matchesDynamicRoute checks URL path `reqURL` matches a Dynamic Route path
// `dynamicPath`. A Dynamic Route is a path string that has the following format: "path/anotherPath/:variablePathname" where `:variablePathname`
// is a catch-all identifier that matches any route with the same structure up to that point.
//
// Ex:
//
//	var ctx = ...
//	var url = "path/anotherPath/andYetAnotherPath"
//	var dynamicPath = "path/anotherPath/:identifier"
//	if !matchesDynamicRoute(&ctx, url, dynamicPath) {
//			panic(...)
//	}
//
// The above code will not panic as the matchesDynamicRoute will evaluate to `true`
func matchesDynamicRoute(urlPath string, routePath string) (isDynamic bool) {
	cleanPath(&urlPath)
	cleanPath(&routePath)

	_, isDynamic = matchDynamicPath(urlPath, routePath)
	return
}

func matchDynamicPath(urlPath, routePath string) (dp []DynamicPath, isDynamic bool) {
	routePathSlice := strings.Split(routePath, "/")
	urlSlice := strings.Split(urlPath, "/")

	if len(routePathSlice) != len(urlSlice) {
		return nil, false
	}

	hasDynamic := false
	dp = []DynamicPath{}
	for i, seg := range routePathSlice {
		if strings.HasPrefix(seg, ":") {
			hasDynamic = true
			dynamicValue := strings.Split(urlSlice[i], "?")[0]
			dp = append(dp, DynamicPath{
				path:  strings.TrimPrefix(seg, ":"),
				value: dynamicValue,
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
