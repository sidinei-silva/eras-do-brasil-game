# game/

Este diretório conterá o projeto Unity do Eras do Brasil.

Será criado na **Sprint 2** com o setup do projeto Unity conforme [ADR-001](../vibe/decisions/ADR-001-organizacao-projeto-unity.md).

## Estrutura Planejada

```
game/
├── Assets/
│   ├── Scripts/
│   │   ├── Core/          # ErasDoBrasil.Core (dados, D20, atributos)
│   │   ├── Combat/        # ErasDoBrasil.Combat (turnos, status, skills)
│   │   ├── World/         # ErasDoBrasil.World (ticks, navegação, eventos)
│   │   ├── NPC/           # ErasDoBrasil.NPC (rotinas, IA, fofoca)
│   │   ├── Economy/       # ErasDoBrasil.Economy (inventário, crafting)
│   │   ├── UI/            # ErasDoBrasil.UI (HUD, menus, telas)
│   │   ├── Persistence/   # ErasDoBrasil.Persistence (save/load)
│   │   ├── Network/       # ErasDoBrasil.Network (P2P, sincronia)
│   │   └── Narrative/     # ErasDoBrasil.Narrative (quests, diálogos)
│   ├── Scenes/
│   │   ├── POC/           # Cenas isoladas para cada POC
│   │   └── Game/          # Cenas do jogo final
│   ├── Resources/         # Assets carregados em runtime
│   ├── Prefabs/
│   ├── Art/               # Sprites, tilesets, UI art
│   ├── Audio/             # Música, SFX
│   └── Tests/
│       ├── EditMode/
│       └── PlayMode/
├── Packages/
└── ProjectSettings/
```
