package config

import "os"

// LogConfig controla quais grupos de logs de debug estão habilitados.
// Cada flag permite ativar/desativar um grupo de forma independente,
// sem alterar o nível global do logger.
//
// Configure via variáveis de ambiente antes de iniciar o servidor.
// Exemplo: LOG_NPC_NEEDS=true LOG_LEVEL=debug ./server
type LogConfig struct {
	// Level define o nível mínimo do logger.
	// Valores aceitos: "debug", "info", "warn", "error"
	// Env: LOG_LEVEL (padrão: "info")
	Level string

	// NPCNeeds exibe o estado de needs (fome/fadiga/solidão) antes e após o decay
	// de cada NPC a cada tick. Muito verboso: 2 linhas × nº de NPCs por tick.
	// Env: LOG_NPC_NEEDS (padrão: false)
	NPCNeeds bool

	// NPCSchedule exibe o schedule carregado por cada NPC na inicialização.
	// Env: LOG_NPC_SCHEDULE (padrão: false)
	NPCSchedule bool

	// CommandRouting exibe o enfileiramento de comandos no GameLoop (move, attack, etc.).
	// Env: LOG_COMMAND_ROUTING (padrão: false)
	CommandRouting bool

	// GameLoopTicks exibe a duração de processamento de cada tick do game loop.
	// Env: LOG_GAME_LOOP_TICKS (padrão: false)
	GameLoopTicks bool

	// WorldLoading exibe detalhes do carregamento de blocos e NPCs na inicialização.
	// Env: LOG_WORLD_LOADING (padrão: false)
	WorldLoading bool

	// NPCBehavior exibe decisões de comportamento dos NPCs (transições de atividade, etc.).
	// Env: LOG_NPC_BEHAVIOR (padrão: false)
	NPCBehavior bool
}

// Log é a configuração global de log. Inicializada uma vez na startup do processo.
var Log = &LogConfig{
	Level:          getEnv("LOG_LEVEL", "info"),
	NPCNeeds:       getEnvBool("LOG_NPC_NEEDS", false),
	NPCSchedule:    getEnvBool("LOG_NPC_SCHEDULE", false),
	CommandRouting: getEnvBool("LOG_COMMAND_ROUTING", false),
	GameLoopTicks:  getEnvBool("LOG_GAME_LOOP_TICKS", false),
	WorldLoading:   getEnvBool("LOG_WORLD_LOADING", false),
	NPCBehavior:    getEnvBool("LOG_NPC_BEHAVIOR", false),
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v == "1" || v == "true" || v == "yes"
}
