package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Server struct {
	mu      sync.Mutex
	reqChan chan interface{}
}

func (s *Server) handleRequest(w http.ResponseWriter, req *http.Request) {
	log.Println("Try to acquire lock")
	if !s.mu.TryLock() {
		msg := fmt.Sprintf("Service 1 is busy (%s)", getIP())
		http.Error(w, msg, http.StatusServiceUnavailable)
		log.Println(msg)
		return
	}

	log.Println("Request received")
	w.Header().Set("Content-Type", "application/json")

	service1Info := Message{
		IP:       getIP(),
		PS:       getPS(),
		DF:       getDF(),
		LastBoot: getLastBoot(),
	}
	service2Info := getService2Info()

	payload := map[string]interface{}{
		"service1": service1Info,
		"service2": service2Info,
	}

	jsonMessage, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonMessage)
	if err != nil {
		log.Println("Error writing response:", err)
		return
	}

	s.reqChan <- true
	log.Println("Request processed")
}

func (s *Server) handleStop(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	log.Println("Initiating shutdown sequence...")
	_, _ = w.Write([]byte("Server stopping - this will trigger container shutdown"))
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	time.Sleep(100 * time.Millisecond)

	// Exit with code 137 (SIGKILL)
	os.Exit(137)
}

func main() {
	server := &Server{
		mu:      sync.Mutex{},
		reqChan: make(chan interface{}),
	}

	go func() {
		for range server.reqChan {
			log.Println("Waiting for 2 seconds")
			time.Sleep(time.Second * 2)
			server.mu.Unlock()
			log.Println("Unlocking Server")
		}
	}()

	http.HandleFunc("/request", server.handleRequest)
	http.HandleFunc("/stop", server.handleStop)

	log.Println("Starting server on port 8199")
	err := http.ListenAndServe(":8199", nil)
	if err != nil {
		log.Println("Error starting server: ", err)
		return
	}
}
