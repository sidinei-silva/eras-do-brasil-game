# 🌐 Text Client — Eras do Brasil

> Cliente web (HTML/CSS/JS vanilla) que se conecta ao servidor via WebSocket.

## Stack

- **HTML / CSS / JS** (vanilla, sem frameworks)
- **WebSocket** nativo do browser (`new WebSocket(...)`)

## Propósito

Interface de texto para o MUD Moderno. Exibe o estado do mundo, ações do jogador e eventos em tempo real via log textual.

## Estrutura Planejada

```
textClient/
├── index.html        # Página principal
├── css/
│   └── style.css     # Estilos
└── js/
    └── main.js       # Conexão WebSocket + renderização
```

## Como usar

1. Suba o servidor: `cd ../server && go run main.go`
2. Abra `index.html` no browser (ou acesse via HTTP servido pelo próprio servidor Go)

## Fase Atual

**Fase 0: Heartbeat** — Exibir "Tick: N" atualizado em tempo real via WebSocket.
