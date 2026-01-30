package static

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const INDEX = "index.html"

// ServeFileSystem is an interface that combines http.FileSystem with a check for file existence.
type ServeFileSystem interface {
	http.FileSystem
	// Exists checks if the given filepath exists under the specified prefix.
	Exists(prefix string, path string) bool
}

type localFileSystem struct {
	http.FileSystem
	root    string
	indexes bool
}

// LocalFile creates a new localFileSystem that serves files from the specified root directory.
//
// Parameters:
//   - root: The root directory path on the local filesystem.
//   - indexes: Whether to enable directory indexing (listing files).
//
// Returns:
//   - *localFileSystem: A filesystem that serves local files.
func LocalFile(root string, indexes bool) *localFileSystem {
	return &localFileSystem{
		FileSystem: gin.Dir(root, indexes),
		root:       root,
		indexes:    indexes,
	}
}

func (l *localFileSystem) Exists(prefix string, filepathStr string) bool {
	if p := strings.TrimPrefix(filepathStr, prefix); len(p) < len(filepathStr) {
		// Join the requested path with the root and ensure it stays within the root directory.
		name := path.Join(l.root, p)

		absRoot, err := filepath.Abs(l.root)
		if err != nil {
			return false
		}

		absName, err := filepath.Abs(name)
		if err != nil {
			return false
		}

		rel, err := filepath.Rel(absRoot, absName)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			return false
		}

		stats, err := os.Stat(absName)
		if err != nil {
			return false
		}
		if stats.IsDir() {
			if !l.indexes {
				index := filepath.Join(absName, INDEX)
				_, err := os.Stat(index)
				if err != nil {
					return false
				}
			}
		}
		return true
	}
	return false
}

// ServeRoot creates a Gin middleware to serve static files from a root directory.
// It simplifies the usage of Serve by automatically creating a LocalFile system.
//
// Parameters:
//   - urlPrefix: The URL prefix to strip from the request path.
//   - root: The local directory to serve files from.
//
// Returns:
//   - gin.HandlerFunc: The middleware handler.
func ServeRoot(urlPrefix, root string) gin.HandlerFunc {
	return Serve(urlPrefix, LocalFile(root, false))
}

// Serve returns a middleware handler that serves static files in the given directory.
// It checks if the file exists using the ServeFileSystem interface before serving.
//
// Parameters:
//   - urlPrefix: The URL prefix to strip from the request path.
//   - fs: The ServeFileSystem implementation to use for file serving.
//
// Returns:
//   - gin.HandlerFunc: The middleware handler.
func Serve(urlPrefix string, fs ServeFileSystem) gin.HandlerFunc {
	fileserver := http.FileServer(fs)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(c *gin.Context) {
		if fs.Exists(urlPrefix, c.Request.URL.Path) {
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}
