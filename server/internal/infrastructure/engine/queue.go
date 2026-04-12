package engine

import (
	"sync"

	"github.com/sidinei-silva/eras-do-brasil-game/server/internal/domain/command"
)

type CommandQueue struct {
	mu       sync.Mutex
	Commands []command.PlayerCommand
}

func NewCommandQueue() *CommandQueue {
	return &CommandQueue{
		Commands: make([]command.PlayerCommand, 0),
	}
}

// Enqueue é chamado pelo WebSocket imediatamente quando a mensagem chega.
func (q *CommandQueue) Enqueue(cmd command.PlayerCommand) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.Commands = append(q.Commands, cmd)
}

// Drain é chamado pelo GameLoop 1 vez a cada tick.
// Ele pega tudo, limpa a fila original e retorna a lista para ser processada.
func (q *CommandQueue) Drain() []command.PlayerCommand {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.Commands) == 0 {
		return nil
	}

	// Cria uma cópia da fila atual
	copyCommands := make([]command.PlayerCommand, len(q.Commands))
	copy(copyCommands, q.Commands)

	// Esvazia a fila original (mantendo a capacidade de memória para eficiência)
	q.Commands = q.Commands[:0]

	return copyCommands
}
