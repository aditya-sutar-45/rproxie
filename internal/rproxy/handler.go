package rproxy

import (
	"io"
	"log"
	"net/http"

	"github.com/aditya-sutar-45/rproxie/internal/utils"
)

func (rp *ReverseProxy) handler(w http.ResponseWriter, r *http.Request) {
	// url := fmt.Sprintf("%s%s", rp.backendAPI, r.URL)
	target := *rp.backends[0].URL
	target.Path = r.URL.Path
	target.RawQuery = r.URL.RawQuery

	request, err := http.NewRequest(r.Method, target.String(), r.Body)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "backend unavailable")
		return
	}
	request.Header = r.Header.Clone()

	resp, err := rp.client.Do(request)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadGateway, "backend unavailable")
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("ERROR streaming response to client: %v", err)
	}
}
