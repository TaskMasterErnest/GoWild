package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

// config constants
const (
	tlsDir        = "/etc/webhook/certs"
	tlsCertFile   = "tls.crt"
	tlsKeyFile    = "tls.key"
	webhookPort   = ":8443"
	readTimeout   = 10 * time.Second
	writeTimeout  = 10 * time.Second
	idleTimeout   = 30 * time.Second
	shutdownGrace = 5 * time.Second
)

// global variables for Kubernetes serialization/deserialization
var (
	scheme       = runtime.NewScheme()
	codecFactory = serializer.NewCodecFactory(scheme)
	deserializer = codecFactory.UniversalDeserializer()
	podResource  = metav1.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}
	allowedContent = "application/json"
)

// Webhook server contains the HTTP server and logger
type WebhookServer struct {
	server *http.Server
	logger *slog.Logger
}

// the admissionHandler defines the signature for mutation/validation functions
type admissionHandler func(*admissionv1.AdmissionRequest) (*admissionv1.AdmissionResponse, error)

func main() {
	// set up structured logging in JSON
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// verify that the TLS files exist before starting the server
	certPath := filepath.Join(tlsDir, tlsCertFile)
	keyPath := filepath.Join(tlsDir, tlsKeyFile)
	// check if the CertFile exists
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		logger.Error("Certificate file not found", "path", certPath)
		os.Exit(1)
	}
	// check if the keyFile exists
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		logger.Error("Private Key file not found", "path", keyPath)
		os.Exit(1)
	}

	// Load the TLS certificate pair
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		logger.Error("Failed to load key pair", "error", err)
		os.Exit(1)
	}

	// create an HTTP server
	mux := http.NewServeMux()
	whs := &WebhookServer{
		logger: logger,
		server: &http.Server{
			Addr:    webhookPort,
			Handler: mux,
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				MinVersion:   tls.VersionTLS13,
			},
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
		},
	}

	// register the handler functions as endpoints to hit
	mux.HandleFunc("/mutate", whs.handleRequest(whs.serveMutatingRequest))
	mux.HandleFunc("/validate", whs.handleRequest(whs.serveValidatingRequest))
	mux.HandleFunc("/healthz", whs.healthCheck)

	// start server in a goroutine so we can handle shutdown signals
	go func() {
		logger.Info("Starting webhook server", "port", webhookPort)
		if err := whs.server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// set up an OS signal channel for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// wait for shutdown signal
	<-sigChan
	logger.Info("Received shutdown signal")

	// create a shutdown context with some grace period
	ctx, cancel := context.WithTimeout(context.Background(), shutdownGrace)
	defer cancel()

	// attempt graceful shutdown
	if err := whs.server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", "error", err)
	}
	logger.Info("Server shutdown completed")
}

/*
**
handleRequest, the middleware that handles common request processing
**
*/
func (whs *WebhookServer) handleRequest(handler admissionHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// validate the HTTP method
		if r.Method != http.MethodPost {
			whs.errorResponse(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// check if content type is JSON
		if contentType := r.Header.Get("content-type"); contentType != allowedContent {
			whs.errorResponse(w, fmt.Sprintf("Unsupported Content Type: %s", contentType), http.StatusBadRequest)
			return
		}

		// read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 {
			whs.errorResponse(w, "Empty or Unreadable Body", http.StatusBadRequest)
			return
		}

		// decode admissionReview request from Kubernetes
		var admissionReview admissionv1.AdmissionReview
		if _, _, err := deserializer.Decode(body, nil, &admissionReview); err != nil {
			whs.errorResponse(w, "Invalid Admission Review Request", http.StatusBadRequest)
			return
		}

		// check if the request exists
		if admissionReview.Request == nil {
			whs.errorResponse(w, "Missing Admission Request", http.StatusBadRequest)
			return
		}

		// call the actual handler to serve the Mutating or Validating Request
		response, err := handler(admissionReview.Request)
		if err != nil {
			whs.logger.Error("Admission review", "error", err)
			// create error response if handler fails
			response = &admissionv1.AdmissionResponse{
				UID:     admissionReview.Request.UID,
				Allowed: false,
				Result: &metav1.Status{
					Message: err.Error(),
					Code:    http.StatusInternalServerError,
				},
			}
		}

		// final response
		admissionReview.Response = response
		admissionReview.Response.UID = admissionReview.Request.UID

		// marshal response into JSON
		res, err := json.Marshal(admissionReview)
		if err != nil {
			whs.errorResponse(w, "Error Encoding Response", http.StatusInternalServerError)
			return
		}

		// send response back to Kubernetes
		w.Header().Set("Content-Type", allowedContent)
		if _, err := w.Write(res); err != nil {
			whs.logger.Error("Failed to write response", "error", err)
		}
	}
}

// errorResponse handles error reporting and logging
func (whs *WebhookServer) errorResponse(w http.ResponseWriter, msg string, code int) {
	whs.logger.Error(msg, "code", code) // Structured error logging
	http.Error(w, msg, code)            // Send HTTP error response
}

// healthCheck handles Kubernetes health check probes
func (whs *WebhookServer) healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK) // Simple 200 OK response
}

// patchOperation defines a single JSON patch operation
type patchOperation struct {
	Op    string      `json:"op"`              // Operation: add, replace, remove
	Path  string      `json:"path"`            // JSON path to modify
	Value interface{} `json:"value,omitempty"` // New value (if needed)
}
