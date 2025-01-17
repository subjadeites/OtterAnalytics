package handlers

import (
	"OtterAnalytics/config"
	"OtterAnalytics/pkg/errors"
	"bufio"
	"gorm.io/gorm"
	"io"
	"log"
	"net"
	"net/http"
)

type Handler struct {
	DB     *gorm.DB
	Config *config.Config
}

func NewHandler(db *gorm.DB, cfg *config.Config) *Handler {
	return &Handler{
		DB:     db,
		Config: cfg,
	}
}

func (h *Handler) HandleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			errors.Normal(err, "Error closing connection")
		}
	}()
	bufReader := bufio.NewReader(conn)

	for {
		req, err := http.ReadRequest(bufReader)
		if err != nil {
			if err == io.EOF {
				log.Println("Client closed the connection")
			} else {
				log.Printf("Error reading request: %v", err)
			}
			break
		}
		log.Printf("Received %s request for %s", req.Method, req.URL.Path)

		methodRoutes, ok := Routes[req.Method]
		if !ok {
			log.Printf("Unsupported method: %s", req.Method)
			WriteMethodNotAllowedResponse(conn)
			_ = req.Body.Close()
			continue
		}

		handler, ok := methodRoutes[req.URL.Path]
		if ok {
			handler(h, req, conn)
		} else {
			log.Printf("No handler found for path: %s", req.URL.Path)
			WriteNotFoundResponse(conn)
		}

		_ = req.Body.Close()
	}
}
