package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
)

// ChatCompletionRequest 描述了向 OpenAI API 发送的请求格式。
type ChatCompletionRequest struct {
	Model    string                  `json:"model"`
	Messages []ChatCompletionMessage `json:"messages"`
}

// ChatCompletionMessage 描述了消息的格式，其中包含角色（如用户、AI）和内容。
type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func main() {

	socks5Addr := "127.0.0.1:10808" //本地代理

	// 创建SOCKS5拨号器
	dialer, err := proxy.SOCKS5("tcp", socks5Addr, nil, proxy.Direct)
	if err != nil {
		fmt.Println("Error creating dialer:", err)
		return
	}

	// 封装拨号器以支持DialContext
	dialerFunc := func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
		return dialer.Dial(network, addr)
	}

	httpTransport := &http.Transport{
		DialContext: dialerFunc,
	}

	httpClient := &http.Client{
		Transport: httpTransport,
	}

	url := "https://api.openai.com/v1/chat/completions" // OpenAI completion endpoint
	apiKey := "your token"

	requestBody := ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []ChatCompletionMessage{
			{
				Role:    "user",
				Content: "Hello!",
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println("Received:", string(body))
}
