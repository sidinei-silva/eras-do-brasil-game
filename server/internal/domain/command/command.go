package command

import "encoding/json"

// Constantes para não usar strings soltas pelo código
const (
	TypeSendChat = "send_chat"
	TypeMove     = "move"
	TypeAttack   = "attack"
)

// A casca genérica que chega da rede
// ClientMessage representa o JSON cru enviado pelo cliente.
// O json.RawMessage guarda o resto do JSON como bytes para decodificar depois.
type ClientMessage struct {
	Type    string          `json:"type"`    // Ex: "send_chat"
	Payload json.RawMessage `json:"payload"` // O resto dos dados do comando
}

// PlayerCommand junta quem enviou com o que foi enviado.
type PlayerCommand struct {
	PlayerID string
	Message  ClientMessage
}

// ==========================================
// PAYLOADS ESPECÍFICOS (Crie um para cada ação)
// ==========================================

// Payload exclusivo para o Chat
type ChatPayload struct {
	Message string `json:"message"`
	Channel string `json:"channel"` // Ex: "global", "guild", "local"
}

// Payload exclusivo para Mover
type MovePayload struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
