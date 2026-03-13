# Bota-Fora Completo da Sessão — Eras do Brasil

**Data de consolidação:** 2026-03-13
**Escopo:** Linha do tempo completa desta sessão de trabalho (desde o primeiro contexto trazido até o estado atual).
**Objetivo:** Servir como handover integral para continuidade em novo chat/novo ambiente (WSL), sem perda de contexto.

---

## 1) Contexto de Entrada (início da sessão)

### O que o usuário trouxe no começo

1. Projeto grande e com histórico longo de decisões feitas por IA (ChatGPT e Gemini).
2. Sensação de desorganização, excesso de escopo e dificuldade de execução prática.
3. Perfil pessoal e restrições reais:
   - dev fullstack sênior
   - pouca disponibilidade (CLT + família)
   - intenção de avançar com poucas horas semanais
4. Dúvidas estratégicas importantes:
   - organização de repositório (mono vs multi)
   - como transformar documentação em execução
   - como usar Vibe Flow para dar governança
5. Desejo de manter visão de MMORPG no longo prazo, com pragmatismo no curto prazo.

### O que o Copilot respondeu no começo

1. Executou leitura ampla de artefatos de projeto e histórico.
2. Criou plano de organização e documentação de base para orientar execução.
3. Estruturou backlog e plano de desenvolvimento com recorte de escopo e cadência realista.

### Decisão macro consolidada no começo

1. Priorizar organização e execução incremental.
2. Formalizar decisões em ADRs e documentos de referência.
3. Transformar o projeto em um fluxo por fases com entregas verificáveis.

---

## 2) Reorganização Inicial e Padronização

### Principais ações realizadas

1. Criação e revisão de documentos estruturantes de produto/planejamento.
2. Revisão de estado, roadmap e backlog para reduzir ambiguidade.
3. Definição de fluxo de trabalho com Vibe Flow e separação clara de papéis (product, tech, vibe).

### Decisões relevantes da etapa

1. Manter governança por documentação viva e iterativa.
2. Priorizar rastreabilidade por sessões em `vibe/sessions/`.

---

## 3) Auditoria de Documentação e Limpeza de Contradições

### Pedido do usuário na etapa

1. Auditar contradições, duplicações e lacunas no conjunto documental.
2. Ajustar o que fosse necessário para reduzir ruído e conflito de fontes.

### Ações executadas

1. Auditoria ampla em roadmap, backlog, guia de retomada e GDD.
2. Correções de inconsistências pontuais (inclusive fórmulas/referências em regras).
3. Consolidação de documentos sobrepostos para reduzir duplicidade.
4. Limpeza de artefatos considerados legados/temporários na época da auditoria.

### Resultado prático

1. Base documental ficou mais enxuta e mais previsível para execução.
2. Projeto avançou para estado de “pronto para iniciar código”.

---

## 4) Pivot Estratégico (momento crítico da sessão)

### O que o usuário explicitou

1. Frustração por anos de documentação sem código efetivo jogável.
2. Necessidade de simplificar caminho de execução.
3. Vontade de retomar direção MMORPG com abordagem mais realista de implementação.

### Decisão arquitetural principal

1. Pivot de Unity Co-op/P2P para **MUD Moderno server-authoritative**:
   - servidor em Go
   - cliente em web (HTML/CSS/JS)
   - comunicação via WebSocket

### Formalização técnica da decisão

1. ADR-004 consolidou o pivot.
2. ADR-005 consolidou arquitetura de processo único + goroutines + EventBus.
3. ADR-006 consolidou persistência RAM-First com snapshots em SQLite (WAL).

### Impacto do pivot

1. Reescrita/atualização ampla de documentação para alinhar produto, técnica e execução.
2. Roadmap reorganizado em 7 fases sequenciais (Heartbeat -> Multiplayer).

---

## 5) Reorganização de Issues e Backlog Operacional

### O que foi feito

1. Arquivamento/fechamento de issues antigas que não refletiam mais a direção pós-pivot.
2. Criação de novo conjunto de issues por fases do roadmap.
3. Análise de reuso de trabalho passado.
4. Reabertura e ajuste de issues reaproveitáveis (com novo contexto Go/MMORPG).

### Decisão operacional

1. Preservar histórico, mas manter backlog ativo estritamente alinhado às fases atuais.
2. Usar issues por fase como trilha de execução incremental.

---

## 6) Pedido de Timeline e Divisão de Responsabilidades

### O que o usuário pediu

1. Plano de execução prático para começar a codar imediatamente.
2. Sem geração de código pronto pelo Copilot, porque o usuário quer aprender implementando tudo.
3. Clareza do que é responsabilidade do usuário e do que é responsabilidade do Copilot.

### Resposta e decisão

1. Usuário escreve todo o código.
2. Copilot atua em arquitetura, revisão, depuração guiada, planejamento e documentação.
3. Timeline orientada por fases curtas e objetivos testáveis.

---

## 7) Reestruturação Física do Repositório (etapa recente)

### O que o usuário pediu

1. Avaliar exclusão de `pocs/` e `spinoffs/`.
2. Estruturar projeto em `game/server/` e `game/textClient/`.
3. Atualizar site em `web/` para refletir o novo posicionamento.

### O que foi executado

1. Remoção de `pocs/` e `spinoffs/` (conteúdo histórico/planejamento antigo).
2. Criação de:
   - `game/server/README.md`
   - `game/textClient/README.md`
3. Atualização de `web/index.html` para linguagem e narrativa pós-pivot.

### Resultado

1. Estrutura do repositório ficou coerente com a execução da Fase 0.
2. Comunicação pública do projeto ficou alinhada com MUD Moderno.

---

## 8) Discussão sobre Admin Commands (RCON como exemplo)

### O que o usuário pediu

1. Planejar comandos administrativos para facilitar dev/testes.
2. Cobrir casos como listar NPCs/inimigos, localizar entidades, inspecionar estado, resetar/mutar mundo.
3. Entender arquitetura e comunicação do subsistema (sem código), incluindo decisão de goroutine dedicada.

### Decisões tomadas

1. Começar com interface local de dev (mais barata e rápida).
2. Prever evolução para endpoint administrativo interno.
3. Deixar remoto (RCON-like) como opcional posterior, não bloqueante.
4. Definir `AdminCommandManager` em goroutine dedicada.
5. Integrar comandos e respostas ao EventBus.
6. Exigir trilha de auditoria e proteção de comandos destrutivos.

### Onde ficou documentado

1. `ROADMAP.md` (trilha transversal de comandos admin)
2. `vibe/backlog.md` (workstream transversal com entregas por fase)
3. `docs/GUIA_RETOMADA.md` (arquitetura e entregas admin por fase)
4. `vibe/project-status.md` (status do workstream admin)
5. `game/server/README.md` (encaixe do subsistema admin no servidor)

---

## 9) Discussão sobre migração de sessão para WSL

### O que o usuário trouxe

1. Migração de ambiente de ZorinOS para WSL.
2. Teste de cópia de workspaceStorage para reaproveitar chat.
3. Ideia de editar state database local para corrigir paths.

### Recomendação consolidada

1. Não depender de edição manual de banco interno da sessão.
2. Usar handover documental versionado como mecanismo principal.
3. Entrar no WSL com contexto carregado por:
   - ADRs
   - status do projeto
   - roadmap
   - logs de sessão

### Decisão

1. Criar e manter bota-fora detalhado para continuidade segura entre ambientes.

---

## 10) Correção de entendimento nesta reta final

### O que aconteceu

1. O usuário apontou corretamente que o bota-fora anterior cobria apenas a parte mais recente.
2. O usuário reforçou que queria o histórico da sessão inteira desde o começo.
3. O usuário colou conteúdo extenso do histórico do chat para viabilizar reconstrução completa.

### Ação executada

1. Refeito o bota-fora em escopo completo (este documento).
2. Expansão da visibilidade dos comandos admin em mais locais de planejamento e operação.

---

## 11) Estado Atual Consolidado

### Situação do projeto

1. Direção técnica consolidada: Go server-authoritative + cliente web.
2. Fase atual: Heartbeat (início de código).
3. Base documental alinhada ao pivot e com rastreabilidade por ADR + sessões.
4. Trilha de comandos admin explicitada em múltiplos documentos.
5. Estrutura física `game/` pronta para início efetivo da implementação pelo usuário.

### Responsabilidades acordadas

1. Usuário implementa todo o código para aprendizado.
2. Copilot apoia com:
   - arquitetura
   - revisão técnica
   - troubleshooting guiado
   - planejamento e documentação

---

## 12) Arquivos-Chave Tocadas/Consolidados nesta janela final

1. `ROADMAP.md`
2. `vibe/backlog.md`
3. `docs/GUIA_RETOMADA.md`
4. `vibe/project-status.md`
5. `game/server/README.md`
6. `game/textClient/README.md`
7. `web/index.html`
8. `vibe/sessions/2026-03-13-admin-comandos-e-migracao-wsl.md`
9. `vibe/sessions/2026-03-13-bota-fora-sessao-completa.md` (este)

---

## 13) Próximo Passo no novo chat (WSL)

1. Abrir o repositório no WSL e validar os arquivos de handover acima.
2. Iniciar Fase 0 com foco em heartbeat e conectividade.
3. Subir primeiro pacote admin read-only (status/tick/entidades) junto da base do servidor.
4. Evoluir comandos admin por fase, sem travar o core loop de desenvolvimento.

---

**Nota final:** Este documento foi escrito para preservar a continuidade completa da sessão em mudança de ambiente, com decisões, contexto, execução e próximos passos.
