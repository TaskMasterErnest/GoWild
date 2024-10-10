package crawler

import (
	"fmt"
	"io"
	"net/http"
)

func GetDataAndResponse(url string) string {
	request, err := http.Get(url)
	if err != nil {
		fmt.Errorf("Error: %v", err)
	}
	defer request.Body.Close()

	data, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Errorf("Error reading response: %v", err)
	}

	return string(data)
}
