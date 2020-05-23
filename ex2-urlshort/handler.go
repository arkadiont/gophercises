package urlshort

import (
	"encoding/json"
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if path, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, path, http.StatusMovedPermanently)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(ymlB []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathURLs, err := parseYml(ymlB)
	if err != nil {
		return nil, err
	}
	return MapHandler(pathURLToMap(pathURLs), fallback), nil
}

// JSONHandler similar to YAML
func JSONHandler(data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathURLs, err := parseJSON(data)
	if err != nil {
		return nil, err
	}
	return MapHandler(pathURLToMap(pathURLs), fallback), nil
}

func pathURLToMap(pathURLs []pathURL) map[string]string {
	pathUrlsMap := make(map[string]string)
	for _, pathURL := range pathURLs {
		pathUrlsMap[pathURL.Path] = pathURL.URL
	}
	return pathUrlsMap
}

func parseYml(data []byte) (r []pathURL, err error) {
	err = yaml.Unmarshal(data, &r)
	return
}

func parseJSON(data []byte) (r []pathURL, err error) {
	err = json.Unmarshal(data, &r)
	return
}

type pathURL struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url"  json:"url"`
}
