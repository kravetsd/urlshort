package handler

import (
	"log"
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
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, pathsToUrls[r.URL.Path], http.StatusFound)
			log.Default().Println("GOT it ! Redirecting to", pathsToUrls[r.URL.Path])
			return
		}
		log.Default().Println("Redirecting to fallback")
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
	pathUrls, err := urlUnmarshal(yml)
	pathMap := pathMapUrls(pathUrls)
	return MapHandler(pathMap, fallback), err
}

func urlUnmarshal(yml []byte) ([]pathUrl, error) {
	pathUrls := []pathUrl{}
	err := yaml.Unmarshal(yml, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

func pathMapUrls(pathUrls []pathUrl) map[string]string {
	pathMap := make(map[string]string)

	for _, path := range pathUrls {
		pathMap[path.Path] = path.Url
	}
	return pathMap
}

type pathUrl struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}
