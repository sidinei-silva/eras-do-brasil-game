# Sessão 6 — Estrutura de Pastas, Site, Comandos Admin e Migração para WSL

> ⚠️ Este registro foi substituído pela versão completa do bota-fora da sessão:
> [2026-03-13-bota-fora-sessao-completa.md](2026-03-13-bota-fora-sessao-completa.md)

**Data:** 2026-03-13
**Tipo:** Planejamento técnico + organização de projeto + handover de contexto
**Participantes:** Dev (usuário) + Copilot

---

## Objetivo da Sessão

1. Definir organização física do projeto pós-pivot.
2. Limpar diretórios históricos que não agregam para o início do código.
3. Atualizar o site público em `web/` para o posicionamento MMORPG server-authoritative.
4. Debater estratégia de comandos administrativos para desenvolvimento/testes.
5. Debater estratégia de migração de sessão do Copilot para ambiente WSL sem perder contexto.
6. Produzir registro completo da sessão para continuidade em novo chat.

---

## Linha do Tempo Detalhada

### 1) Pergunta inicial sobre limpeza e estrutura

**Usuário trouxe:**
- Dúvida se deveria apagar `pocs/` e `spinoffs/`.
- Intenção de criar estrutura nova em `game/`, com `server/` e `textClient/`.

**Copilot respondeu:**
- Verificou que `pocs/` era histórico e sem código executável.
- Confirmou que `spinoffs/` não havia sido encontrado na primeira checagem.
- Recomendou estrutura `game/server/` + `game/textClient/` e pediu confirmação para executar.

**Decisão parcial:**
- Prosseguir com limpeza e criação de estrutura após confirmação explícita.

---

### 2) Confirmação do usuário + correção de contexto sobre `spinoffs/`

**Usuário trouxe:**
- Confirmou que `spinoffs/` existia.
- Autorizou apagar `pocs/`.
- Pediu criação da estrutura nova com READMEs.
- Pediu atualização do site em `web/` para refletir o novo momento do projeto.

**Copilot executou:**
- Removeu `pocs/` e `spinoffs/`.
- Criou:
  - `game/server/README.md`
  - `game/textClient/README.md`
- Atualizou `web/index.html` para alinhar com o pivot MMORPG.

**Mudanças no site (`web/index.html`):**
- Título e descrição: de foco em RPG tático para foco em MUD Moderno.
- Hero status: passou a indicar Fase 0 (Heartbeat, servidor Go + WebSocket).
- Card de gameplay: removeu mensagem de Co-op P2P e posicionou multiplayer persistente.
- Notícias: removeu foco em Unity/POCs antigas, incluiu notícia de pivot, ADRs e fases.
- Footer e textos principais: alinhados à nova direção do projeto.

**Decisão:**
- Estrutura física pós-pivot consolidada em `game/`.
- Diretórios legados removidos para reduzir ruído e acelerar início do código.

---

### 3) Discussão técnica sobre comandos admin (RCON como exemplo)

**Usuário trouxe:**
- Necessidade de comandos de admin para desenvolvimento e testes.
- Exemplos práticos solicitados:
  - listar NPCs
  - localizar NPC/inimigos
  - detalhes de entidades
  - estado do mundo
  - reset/mutação de estado
- Disse que "RCON" era apenas exemplo (aceita outra abordagem).

**Copilot respondeu (sem código):**
- Recomendou começar por console local de desenvolvimento (mais simples, menor custo).
- Sugeriu deixar acesso remoto tipo RCON-like para fase posterior, quando houver necessidade real.
- Defendeu auditoria e controles de segurança para comandos destrutivos.

**Decisão funcional:**
- Criar trilha formal de comandos admin no planejamento do projeto.
- Arquitetura orientada por manager dedicado + integração por eventos.

---

### 4) Discussão de migração de sessão para WSL

**Usuário trouxe:**
- Ambiente atual em ZorinOS com VS Code Insiders.
- Plano de migração para WSL sem perder sessão.
- Teste feito copiando conteúdo de `workspaceStorage` com aparente sucesso.
- Considerou editar `state.vscdb` para ajustar caminhos.

**Copilot respondeu:**
- Não recomendou editar `state.vscdb` como estratégia principal (frágil e propensa a quebra futura).
- Recomendou priorizar continuidade por artefatos de projeto (ADRs, roadmap, status e log de sessão), que são estáveis e versionáveis.
- Concluiu que, dado o contexto já consolidado no repositório, o risco de perda de contexto real é baixo.

**Decisão operacional:**
- Usar handover documentado como mecanismo principal de continuidade.
- Evitar depender de hacks em banco interno do editor.

---

### 5) Solicitação final do usuário nesta sessão

**Usuário pediu explicitamente:**
- "Colocar os comandos nos nossos planos".
- Incluir arquitetura/estrutura dos comandos admin, comunicação entre componentes e decisão sobre goroutine dedicada.
- Criar registro completo desta sessão em formato de linha do tempo, sem poupar detalhes, para abrir novo chat no WSL.

**Copilot executou (sem código de implementação):**
- Atualizou os planos com a trilha de comandos admin.
- Criou este log detalhado para handover.

---

## Decisões Consolidadas da Sessão

1. **Limpeza de legado:** `pocs/` e `spinoffs/` foram removidos.
2. **Nova organização-base:** `game/server/` e `game/textClient/` criados com READMEs.
3. **Comunicação pública alinhada:** `web/index.html` atualizado para o posicionamento MMORPG server-authoritative.
4. **Comandos admin entram no plano oficial:** abordagem incremental por fase, sem travar início do desenvolvimento.
5. **Arquitetura admin definida em alto nível:** manager dedicado, fluxo por eventos e auditoria obrigatória.
6. **Migração de sessão:** preferir handover documental e evitar editar `state.vscdb` manualmente.

---

## Arquivos Alterados Nesta Sessão

### Removidos
- `pocs/` (diretório completo)
- `spinoffs/` (diretório completo)

### Criados
- `game/server/README.md`
- `game/textClient/README.md`
- `vibe/sessions/2026-03-13-admin-comandos-e-migracao-wsl.md` (este arquivo)

### Atualizados
- `web/index.html`
- `ROADMAP.md`
- `vibe/backlog.md`

---

## Arquitetura Planejada — Comandos Admin (Sem Código)

### Componentes

- **AdminCommandManager (goroutine dedicada):**
  - recebe comandos administrativos
  - valida contexto e permissões
  - roteia ações para managers de domínio (mundo, NPC, combate, economia etc.)
  - consolida resposta para observabilidade

- **Fontes de comando por maturidade do projeto:**
  - Fase inicial: console local de desenvolvimento
  - Fases intermediárias: endpoint interno administrativo
  - Fase avançada: remoto opcional (RCON-like), se necessário

- **Integração de comunicação:**
  - entrada do comando em tópico administrativo
  - execução por manager de domínio
  - retorno em tópico de resultado
  - emissão de evento de auditoria

### Categorias de Comandos

- **Inspeção:** listagens e snapshots de estado
- **Diagnóstico:** localização e detalhes de entidades
- **Tempo/simulação:** pause/resume/advance tick
- **Mutação controlada:** set/move/spawn/despawn/reset
- **Suporte a QA:** preparação de cenários específicos para teste

### Segurança e Governança

- ambiente dev/homolog com recursos destrutivos habilitáveis
- produção com comandos sensíveis desabilitados por padrão
- trilha de auditoria obrigatória para qualquer mutação
- controles de escopo para evitar mudanças irreversíveis acidentais

---

## Próximos Passos Recomendados (Para o novo chat no WSL)

1. Validar que o workspace abriu com os arquivos novos em `game/` e planos atualizados.
2. Iniciar Fase 0 do servidor Go.
3. Antes de implementar comandos admin completos, começar com pacote mínimo read-only (status/tick/entidades).
4. Evoluir o subsistema admin junto das fases, sem bloquear o coding loop principal.

---

**Gerenciado por:** Copilot
**Observação:** Registro criado para handover entre ambientes (ZorinOS -> WSL) com rastreabilidade completa desta sessão.
