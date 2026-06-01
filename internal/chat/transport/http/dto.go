package chat_transport_http

import (
	chat_transport "github.com/egorkto/Chat-go/internal/chat/transport"
	transport_http "github.com/egorkto/Chat-go/internal/transport/http"
)

type ErrorResponse = transport_http.ErrorResponse
type GetHistoryResponse = []chat_transport.MessageDTO
