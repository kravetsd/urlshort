package handler

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...
	type pathUrl struct {
		path string
		url  string
	}
	var pathUrls []pathUrl
	for path, url := range pathsToUrls {
		pathUrls = append(pathUrls, pathUrl{path, url})
	}
	return func(w http.ResponseWriter, r *http.Request) {
		for _, pathUrl := range pathUrls {
			if r.URL.Path == pathUrl.path {
				http.Redirect(w, r, pathUrl.url, http.StatusFound)
				return
			}
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
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// TODO: Implement this...
	type pathUrl struct {
		Path string `yaml:"path"`
		URL  string `yaml:"url"`
	}
	var pathUrls []pathUrl
	err := yaml.Unmarshal(yml, &pathUrls)
	if err != nil {
		return nil, err
	}
	pathMap := make(map[string]string)
	for _, pathUrl := range pathUrls {
		pathMap[pathUrl.Path] = pathUrl.URL
	}
	fmt.Println(pathMap)
	return MapHandler(pathMap, fallback), err
}
