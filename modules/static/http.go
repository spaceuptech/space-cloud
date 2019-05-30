package static

import (
	"bufio"
	"net/http"
	"strings"
)

// HandleRequest creates a static request endpoint
func (m *Module) HandleRequest(w http.ResponseWriter, r *http.Request) {
	url := r.URL.EscapedPath()
	host := strings.Split(r.Host, ":")[0]

	route, ok := m.selectRoute(host, url)
	if !ok {
		http.Error(w, "Path not found", http.StatusNotFound)
		return
	}

	path := strings.TrimPrefix(url, route.URLPrefix)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	path = route.Path + path

	// Its a proxy request
	if route.Proxy != "" {
		addr := route.Proxy + path
		req, err := http.NewRequest(r.Method, addr, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Set the http headers
		req.Header = make(http.Header)
		if contentType, p := r.Header["Content-Type"]; p {
			req.Header["Content-Type"] = contentType
		}

		// Make the http client request
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		defer res.Body.Close()

		reader := bufio.NewReader(res.Body)

		w.Header().Set("Content-Type", res.Header.Get("Content-Type"))
		w.WriteHeader(res.StatusCode)
		reader.WriteTo(w)
		return
	}

	http.ServeFile(w, r, path)
}
