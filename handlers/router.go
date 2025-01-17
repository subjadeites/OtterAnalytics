package handlers

import (
	"fmt"
	"net"
	"net/http"
)

type HandlerFunc func(h *Handler, r *http.Request, conn net.Conn)

var Routes = make(map[string]map[string]HandlerFunc)

func init() {
	initializeRoutes()
}

func initializeRoutes() {
	Routes[http.MethodGet] = map[string]HandlerFunc{"/": handleIndexGet}
	Routes[http.MethodPost] = map[string]HandlerFunc{}
}

type Header struct {
	contentType string
	connection  string
}

func NewHeader(contentTypes ...string) *Header {
	contentType := "text/plain"
	if len(contentTypes) > 0 {
		contentType = contentTypes[0]
	}
	return &Header{
		contentType: contentType,
		connection:  "keep-alive",
	}
}

func writeResponse(conn net.Conn, header *Header, statusCode int, body string) error {
	response := fmt.Sprintf(
		"HTTP/1.1 %d %s\r\nContent-Type: %s\r\nContent-Length: %d\r\nConnection: %s\r\n\r\n%s",
		statusCode, http.StatusText(statusCode),
		header.contentType, len(body),
		header.connection,
		body,
	)
	_, err := conn.Write([]byte(response))
	return err
}

func WriteResponse(conn net.Conn, header *Header, body string) error {
	return writeResponse(conn, header, 200, body)
}

// 403 Forbidden
func WriteForbiddenResponse(conn net.Conn) {
	responseHeader := NewHeader()
	_ = writeResponse(conn, responseHeader, 403, "403 Forbidden")
}

// 404 Not Found
func WriteNotFoundResponse(conn net.Conn) {
	responseHeader := NewHeader()
	_ = writeResponse(conn, responseHeader, 404, "404 Not Found")
}

// 405 Method Not Allowed
func WriteMethodNotAllowedResponse(conn net.Conn) {
	responseHeader := NewHeader()
	_ = writeResponse(conn, responseHeader, 405, "405 Method Not Allowed")
}

// 500 Internal Server Error
func WriteServerErrorResponse(conn net.Conn) {
	responseHeader := NewHeader()
	_ = writeResponse(conn, responseHeader, 500, "500 Internal Server Error")
}

func handleIndexGet(h *Handler, r *http.Request, conn net.Conn) {
	responseHeader := NewHeader()
	err := WriteResponse(conn, responseHeader, "Hello, World!")
	if err != nil {
		return
	}
}
