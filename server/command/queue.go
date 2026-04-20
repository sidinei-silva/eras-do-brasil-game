package command

import (
	"sync"
)

type CommandQueue struct {
	mu       sync.Mutex
	Commands []PlayerCommand
}

func NewCommandQueue() *CommandQueue {
	return &CommandQueue{
		Commands: make([]PlayerCommand, 0),
	}
}

// Enqueue é chamado pelo WebSocket imediatamente quando a mensagem chega.
func (q *CommandQueue) Enqueue(cmd PlayerCommand) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.Commands = append(q.Commands, cmd)
}

// Drain é chamado pelo GameLoop 1 vez a cada tick.
// Ele pega tudo, limpa a fila original e retorna a lista para ser processada.
func (q *CommandQueue) Drain() []PlayerCommand {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.Commands) == 0 {
		return nil
	}

	// Cria uma cópia da fila atual
	copyCommands := make([]PlayerCommand, len(q.Commands))
	copy(copyCommands, q.Commands)

	// Esvazia a fila original (mantendo a capacidade de memória para eficiência)
	q.Commands = q.Commands[:0]

	return copyCommands
}
