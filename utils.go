package goster

import (
	"fmt"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Adds basic headers
func DefaultHeader(c *Ctx) {
	// TODO: refactor
	c.Response.Header().Set("Access-Control-Allow-Origin", "*")
	c.Response.Header().Set("Connection", "Keep-Alive")
	c.Response.Header().Set("Keep-Alive", "timeout=5, max=997")
}

// cleanPath sanatizes a URL path. It removes suffix '/' if any and adds prefix '/' if missing
func cleanPath(path *string) {
	if len(*path) == 0 {
		*path = "/"
		return
	}

	if (*path)[0] != '/' {
		*path = "/" + *path
	}

	if len(*path) != 1 {
		*path = strings.TrimSuffix(*path, "/")
	}
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

func hexEscapeNonASCII(s string) string {
	newLen := 0
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			newLen += 3
		} else {
			newLen++
		}
	}
	if newLen == len(s) {
		return s
	}
	b := make([]byte, 0, newLen)
	var pos int
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			if pos < i {
				b = append(b, s[pos:i]...)
			}
			b = append(b, '%')
			b = strconv.AppendInt(b, int64(s[i]), 16)
			pos = i + 1
		}
	}
	if pos < len(s) {
		b = append(b, s[pos:]...)
	}
	return string(b)
}

var htmlReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	// "&#34;" is shorter than "&quot;".
	`"`, "&#34;",
	// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
	"'", "&#39;",
)

func htmlEscape(s string) string {
	return htmlReplacer.Replace(s)
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
