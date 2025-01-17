package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

type RequestDataReader interface {
	ReadData(r *http.Request) ([]byte, string, error)
}

type ReadPostData struct{}

func (g *ReadPostData) ReadData(r *http.Request) ([]byte, string, error) {
	contentLengthString := r.Header.Get("Content-Length")
	if contentLengthString == "" {
		log.Println("No content length provided")
		return nil, "", nil
	}
	contentLength, err := strconv.Atoi(contentLengthString)
	if err != nil {
		log.Printf("Error parsing content length: %v", err)
		return nil, "", nil
	}
	contentType := r.Header.Get("Content-Type")
	contentBody := r.Body

	// check length of body == header content-length
	if contentLength != 0 {
		body := make([]byte, contentLength)
		_, err := io.ReadFull(contentBody, body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			return nil, "", err
		}

		// match case for content type
		switch contentType {
		case "application/json":
			return body, "json", nil
		default:
			return nil, "", nil
		}

	} else {
		return nil, "", nil
	}
}

type ReadGetData struct{}

func (g *ReadGetData) ReadData(r *http.Request) ([]byte, string, error) {
	query := r.URL.Query()
	queryMap := make(map[string]string, len(query))
	for key, values := range query {
		if len(values) > 0 {
			queryMap[key] = values[0]
		}
	}

	queryJSON, err := json.Marshal(queryMap)
	if err != nil {
		return nil, "error", err
	}

	return queryJSON, "query", nil
}
