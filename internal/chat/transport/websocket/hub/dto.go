package chat_transport_websocket_hub

type MessageInput struct {
	Text string `json:"text"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
