package sse

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

func ServerMain(ctx context.Context, addr string) error {
	router := httprouter.New()
	router.GET("/sse/:id", handleSse)
	return http.ListenAndServe(addr, router)
}

func handleSse(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	iterStr := r.URL.Query().Get("iter")
	iter, _ := strconv.Atoi(iterStr)
	if iter == 0 {
		iter = 10
	}

	id := ps.ByName("id")
	log.Println("Received a request from", r.RemoteAddr, "id=", id, "iter=", iter)

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ctx := r.Context() // Get the request context

	for i := 0; i < iter; i++ {
		msgStr := fmt.Sprintf("data: %s|%s|message=%d\n\n",
			time.Now().Format(time.RFC3339), id, i)

		_, err := w.Write([]byte(msgStr))
		if err != nil {
			log.Println("Connection is broken: err=", err)
			return // Stop if the connection is already broken
		}
		flusher.Flush()

		// Idiomatic way to handle cancellation during wait
		waitDuration := time.Duration(rand.Int31n(200)) * time.Millisecond

		select {
		case <-ctx.Done():
			// The client disconnected or the request timed out
			log.Println("Client disconnected:", r.RemoteAddr)
			return
		case <-time.After(waitDuration):
			// The timer finished, continue to the next loop iteration
			continue
		}
	}
}
