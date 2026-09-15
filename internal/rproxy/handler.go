package rproxy

import (
	"fmt"
	"io"
	"net/http"

	"github.com/aditya-sutar-45/rproxie/internal/utils"
)

func (rp *ReverseProxy) handler(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s%s", rp.backendAPI, r.URL)
	method := r.Method

	request, err := http.NewRequest(method, url, r.Body)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "backend unavailable")
		return
	}

	for key, values := range r.Header {
		for _, value := range values {
			request.Header.Set(key, value)
		}
	}

	client := http.Client{}
	resp, err := client.Do(request)
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
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
	}
}
