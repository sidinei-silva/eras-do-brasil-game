# ADR-003: Estratégia de Repositórios e Organização do Monorepo

## Status

Aceito (atualizado em 2026-03-11 após ADR-004)

## Contexto

O projeto Eras do Brasil cresceu organicamente com contribuições de 3 IAs diferentes (GPT, Gemini, Copilot) ao longo de meses. Atualmente tudo está concentrado em um monorepo (`eras-do-brasil`), enquanto o histórico de conversas está separado em `ErasDoBrasil-Historico-Legado`.

O criador questiona se deveria:
- Separar GDD, código e POCs em repos distintos
- Manter monorepo
- Usar submodules para conectar tudo

**Restrições:**
- Solo developer com tempo limitado (3 noites/semana)
- IA (GitHub Copilot) é a ferramenta primária de desenvolvimento
- GDD referencia POCs e vice-versa constantemente
- Futuro projeto Unity ficará pesado com assets

## Decisão

**Manter monorepo organizado** com fronteiras claras entre áreas, usando a seguinte estrutura:

```
eras-do-brasil/
├── .github/                    # Copilot instructions, agentes, workflows
├── docs/                       # Documentação oficial (product + tech)
│   ├── product/                # Specs de produto formalizadas
│   └── tech/                   # Arquitetura, ADRs exportados, workflows
├── gdd/                        # Game Design Document (livros, schemas, dados)
│   ├── 01_Livro_de_Regras/
│   ├── 02_Livro_de_Classes/
│   ├── 03_Enredo_e_Mundo/
│   ├── 04_Design_Visual/
│   ├── 05_Livros_Auxiliares/
│   ├── 06_Dados_e_Assets/
│   └── 99_Meta_e_Backlog/
├── server/                     # 🆕 Servidor Go (criado na Fase 0)
│   ├── main.go
│   ├── world/                  # Blocos, NPCs, tick loop
│   ├── combat/                 # Motor D20, combate
│   ├── economy/                # Itens, crafting, comércio
│   └── data/                   # JSONs de templates
├── web/                        # Cliente web (HTML/CSS/JS)
├── vibe/                       # Vibe Flow (decisões, sessões, backlog)
├── ROADMAP.md
└── README.md
```

**Mudanças em relação ao estado atual:**
1. Código Unity vai em `game/` (não na raiz)
2. Nova pasta `spinoffs/` para jogos menores de aprendizado
3. POCs specs (`pocs/`) continuam separadas do código (`game/Assets/Scripts/`)
4. `docs/product/` será populado com specs formais de produto

## Consequências

**Positivas:**
- Contexto único para IA — Copilot vê tudo num repo
- Links internos entre GDD e código funcionam
- Sem overhead de submodules ou multi-repo sync
- Pasta `game/` isolável via `.gitignore` se necessário
- Estrutura clara: docs é docs, game é game

**Negativas:**
- Repo pode ficar grande (mitigado por Git LFS para assets grandes)
- Commits de docs e código misturados no histórico (mitigado por prefixos de commit: `docs:`, `game:`, `gdd:`)

## Alternativas Consideradas

- **Multi-repo (GDD + Game + POCs separados):** Rejeitado — overhead demais para solo dev, perda de contexto para IA, links quebrados
- **Git submodules:** Rejeitado — complexidade desnecessária, experiência de DX ruim
- **Repo por spin-off:** Possível no futuro, mas por agora spin-offs são pequenos o suficiente para ficar juntos
