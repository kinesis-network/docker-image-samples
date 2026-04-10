package sse

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/julienschmidt/httprouter"
)

type (
	// Minimal OpenAI types for the stream response
	ChatCompletionChunk struct {
		ID      string   `json:"id"`
		Object  string   `json:"object"`
		Created int64    `json:"created"`
		Model   string   `json:"model"`
		Choices []Choice `json:"choices"`
		Usage   *Usage   `json:"usage,omitempty"`
	}
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	}
	Choice struct {
		Delta ChoiceDelta `json:"delta"`
	}
	ChoiceDelta struct {
		Content string `json:"content,omitempty"`
		Role    string `json:"role,omitempty"`
	}
	ChatRequest struct {
		Model      string    `json:"model"`
		MaxTokens  int       `json:"max_tokens"`
		MaxCTokens int       `json:"max_completion_tokens"` // guidellm uses this too
		Messages   []Message `json:"messages"`
	}
	Message struct {
		Role    string      `json:"role"`
		Content interface{} `json:"content"` // Use interface{} to handle string OR array
	}
)

var activeChatSession int32

func ServerMain(ctx context.Context, addr string) error {
	router := httprouter.New()
	router.GET("/sse/:id", handleSse)
	// 1. The specific health check guidellm is looking for
	router.GET("/health", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		log.Println("Received a request for /health: ", r.RemoteAddr)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	// 3. Update your /v1/models to match the model name you pass to guidellm
	// Since you used "mistralai/Mistral-Nemo-Instruct-2407" in your command,
	// the server should ideally claim to be that model.
	router.GET("/v1/models", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		log.Println("Received a request for /v1/models: ", r.RemoteAddr)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"object":"list","data":[{"id":"mistralai/Mistral-Nemo-Instruct-2407","object":"model"}]}`)
	})
	router.POST("/v1/chat/completions", handleChatCompletion)
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

func handleChatCompletion(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	n := atomic.AddInt32(&activeChatSession, 1)
	defer atomic.AddInt32(&activeChatSession, -1)
	log.Println("Received a chat request from", r.RemoteAddr, "n=", n)

	var chatReq ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&chatReq); err != nil {
		log.Printf("Decode error: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	defer log.Println("Finishing the request", r.RemoteAddr)

	// guidellm uses max_completion_tokens; check both fields
	limit := chatReq.MaxTokens
	if chatReq.MaxCTokens > 0 {
		limit = chatReq.MaxCTokens
	}
	if limit <= 0 {
		limit = 600
	}

	// Extract content using our helper
	userMsg := ""
	if len(chatReq.Messages) > 0 {
		userMsg = getMessageContent(chatReq.Messages[len(chatReq.Messages)-1].Content)
	}

	// 1. Headers MUST be set before any writing
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, _ := w.(http.Flusher)

	// 2. SEND FIRST CHUNK IMMEDIATELY (Reduces TTFT, establishes role)
	firstChunk := ChatCompletionChunk{
		ID:      "mock-123",
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   chatReq.Model,
		Choices: []Choice{{Delta: ChoiceDelta{Role: "assistant", Content: ""}}},
	}
	initialData, _ := json.Marshal(firstChunk)
	fmt.Fprintf(w, "data: %s\n\n", initialData)
	flusher.Flush()

	var words []string
	if strings.Contains(userMsg, "get me ") {
		url := strings.TrimPrefix(userMsg, "get me ")
		words = strings.Split(fetchWikiText(r.Context(), url), " ")
	} else {
		// FALLBACK: Generate dummy words for guidellm benchmark
		words = make([]string, limit)
		for i := 0; i < limit; i++ {
			words[i] = "token"
		}
	}

	// 4. Stream loop
	for i, word := range words {
		if i >= limit {
			break
		}

		chunk := ChatCompletionChunk{
			ID:      "mock-123",
			Object:  "chat.completion.chunk",
			Created: time.Now().Unix(),
			Model:   chatReq.Model,
			Choices: []Choice{{Delta: ChoiceDelta{Content: word + " "}}},
		}

		jsonData, _ := json.Marshal(chunk)
		fmt.Fprintf(w, "data: %s\n\n", jsonData)
		flusher.Flush()

		time.Sleep(5 * time.Millisecond) // Faster for 96 conc
		if r.Context().Err() != nil {
			return
		}
	}

	// 5. Final Usage Chunk (CRITICAL for guidellm metrics)
	finalUsage := ChatCompletionChunk{
		ID:      "mock-123",
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   chatReq.Model,
		Choices: []Choice{},
		Usage: &Usage{
			PromptTokens:     200,
			CompletionTokens: limit,
			TotalTokens:      200 + limit,
		},
	}
	usageData, _ := json.Marshal(finalUsage)
	fmt.Fprintf(w, "data: %s\n\n", usageData)
	flusher.Flush()

	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}

// Helper to extract text from the interface{} content
func getMessageContent(content interface{}) string {
	switch v := content.(type) {
	case string:
		return v
	case []interface{}:
		// Handle the array of objects guidellm is sending
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				if m["type"] == "text" {
					return fmt.Sprintf("%v", m["text"])
				}
			}
		}
	}
	return ""
}

func fetchWikiText(ctx context.Context, url string) string {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "Error creating a request: " + err.Error()
	}
	req.Header.Add("user-agent", "kinesis-test-client")
	req.Header.Add("accept-encoding", "gzip")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "Error fetching URL: " + err.Error()
	}
	defer resp.Body.Close()

	var reader io.ReadCloser
	// Check if the server actually sent back Gzip
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return "Gzip Error: " + err.Error()
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	body, _ := io.ReadAll(reader)
	html := string(body)

	// Extract only the content inside <p> tags for a cleaner "LLM-like" look
	re := regexp.MustCompile(`(?s)<p>(.*?)</p>`)
	matches := re.FindAllStringSubmatch(html, -1)

	var sb strings.Builder
	tagRemover := regexp.MustCompile(`<[^>]*>`)

	for _, m := range matches {
		cleanText := tagRemover.ReplaceAllString(m[1], "")
		sb.WriteString(cleanText + "\n\n")
		if sb.Len() > 5000 {
			break
		} // Cap it so we don't stream for hours
	}

	if sb.Len() == 0 {
		return "No text content found at that URL."
	}
	return sb.String()
}
