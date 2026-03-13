# 📦 MVP "O Despertar" — Especificação de Produto

**Data:** 2026-03-07
**Status:** Proposto (aguarda aprovação)

---

## Visão

Um loop jogável de **15-30 minutos** que demonstra o core do Eras do Brasil:
- Criar personagem
- Explorar um micro-mundo
- Lutar contra inimigos
- Completar uma quest
- Salvar e carregar progresso

---

## Escopo Exato

### Personagem
- **1 Origem:** Indígena
- **1 Classe:** Guerreiro Tribal (Tier 1)
- **Criação:** Point-buy de atributos (27 pontos, 6 atributos)
- **Habilidades:** 3 ativas (Golpe Forte, Grito de Guerra, Investida) + ataque básico
- **Progressão:** 3 níveis (1 → 2 → 3) com XP funcional

### Mundo
- **3 Nós:** Vila de São Tomé (hub), Floresta do Norte, Ruínas do Eco
- **Navegação:** Mapa de nós com custo em Ticks para viagem
- **Ticks:** Relógio universal rodando (dia/noite visual, sem impacto mecânico no MVP)

### Combate
- **Tipo:** Estático (sem grid), por turnos
- **Mecânica:** D20 + modificadores vs Defesa. HP, dano, status básicos
- **Inimigos:** 3 tipos (Jacaré Corrompido, Saúva Gigante, Guardião Menor)
- **Encontros:** 2-3 por nó (Floresta e Ruínas)

### Quest
- **1 Quest:** "O Caçador que Não Voltou"
  - NPC em Vila pede ajuda → ir à Floresta → encontrar rastros → combater → voltar com prova
  - Recompensa: XP + item
- **Implementação:** Hardcoded (sem sistema de quests genérico)

### UI
- **HUD:** HP bar, XP bar, relógio de Ticks, botão de inventário
- **Mapa:** Tela com 3 nós clicáveis
- **Combate:** UI de ações (ataque, habilidades, item, fugir)
- **Inventário:** Lista simples com equip slot

### Persistência
- **Save/Load:** JSON local (1 slot)

---

## O que NÃO está no MVP

| Feature | Motivo | Quando |
|---------|--------|--------|
| Múltiplas origens/classes | Escopo controlado | Fase 2 |
| Grid tático | Combate estático é suficiente | Fase 4 |
| NPCs com rotinas/IA | Complexidade desnecessária | Fase 2 |
| Crafting/comércio | Economia não é core do MVP | Fase 2 |
| Multiplayer | Muito cedo | Fase 6 |
| Diálogos ramificados | Quest hardcoded | Fase 4 |
| Pixel art final | Texto/web suficiente | Fase 5 |
| Áudio | Sem impacto no gameplay | Futuro |

---

## Stack Técnica do MVP

| Componente | Tecnologia |
|------------|------------|
| Servidor | Go 1.22+ |
| Cliente | HTML / CSS / JS (vanilla) |
| Comunicação | WebSocket (gorilla/websocket) |
| Dados | JSON (templates + saves) |

> O MVP corresponde à **Fase 3 (Player)** do [Roadmap](../../gdd/99_Meta_e_Backlog/02_Roadmap_Desenvolvimento.md). As Fases 0-2 (Heartbeat, Living World, Observer) são pré-requisitos.

---

## Critérios de Sucesso

1. Loop completo: criar → explorar → lutar → quest → salvar → carregar → continuar
2. Combate resolve corretamente (D20, modificadores, dano, morte)
3. Progressão de nível funciona (XP → level up → atributos/habilidades melhoram)
4. Save/Load preserva estado do jogo
5. Servidor Go roda e cliente web conecta via WebSocket
6. Playtest de 3+ runs sem bugs críticos

---

## POCs Necessárias

| POC | Módulo | Descrição |
|-----|--------|-----------|
| 01 | Core | Fundação de Dados |
| 02 | Core | Motor D20 |
| 03 | Core | Atributos e Criação |
| 06 | Core | Matriz de Itens |
| 07 | World | Motor de Ticks |
| 08 | World | Navegação por Blocos |
| 17 | Combat | Combate Estático |
| 18 | Combat | Status e Condições |
| 19 | Combat | Habilidades em Combate |
| 24 | Economy | Inventário |
| 28 | UI | HUD Principal |
| 33 | Persist | Save/Load |

**Total:** 12 POCs (de 38 originais)

---

*Este documento será a referência para o `@product-manager` validar e para o `@developer` implementar.*
