# Sistema de Vilas Dinâmico

> Visão futura. Não está no MVP nem nas próximas fases. Documentado aqui para preservar a ideia e estruturar o pensamento de design.

## Status

**Não implementar agora.** Este documento existe para:

- Registrar a visão antes que ela se perca
- Estruturar as decisões de design que precisam ser tomadas
- Listar pré-requisitos técnicos
- Servir de referência quando chegar a fase de implementação

**Fase prevista:** Pós-MVP. Provavelmente Temporada 2 ou 3 (1-2 anos após o lançamento do MVP).

## Visão

Vilas no mundo de Eras do Brasil são entidades vivas com ciclo de vida próprio. Elas:

- Nascem espontaneamente em locais predefinidos quando há população
- Crescem em tier (1 → 2 → 3...) através de quests realizadas por jogadores
- Decaem em tier automaticamente se ficam sem ocupantes
- Podem ser destruídas por eventos (ataques de facções, desastres)
- Podem ser reconstruídas via decay reverso ou ação coordenada de jogadores

Esse sistema concretiza o pilar fundamental do jogo: o mundo central, vivendo com ou sem participação direta dos jogadores. Jogadores aceleram processos, mas não os causam exclusivamente.

## Como conecta com o North Star

Cada mecanismo do sistema reforça a identidade emergente do jogo:

**Vilas evoluem sem jogador → mundo é vivo, não esperando o herói.**

**Vilas decaem sem manutenção → consequências reais existem, jogador não está protegido pela narrativa.**

**Reconstrução acelerada por quests → jogador importa, mas é catalisador, não salvador único.**

**NPCs migrantes carregam fofocas → história se espalha pelo mundo, não fica restrita ao jogador presente.**

**Sobreviventes de vilas destruídas reaparecem em outras vilas → continuidade narrativa sem precisar simular tudo.**

## Os 4 mecanismos auto-balanceadores

### Mecanismo 1: Tier define vagas

Cada vila tem um tier (1 a N). Tier define quantas e quais vagas profissionais existem:

- **Tier 1**: vagas básicas (Ferreiro, Taverneiro, Curandeiro, Guarda)
- **Tier 2**: adiciona vagas especializadas (Mercador, Padre, Escriba)
- **Tier 3**: adiciona vagas avançadas (Alquimista, Capitão de Guarda, Líder Comunitário)
- **Tier N**: vagas únicas e raras

Vagas são tipadas por `Role`. Cada vaga aceita apenas NPCs daquele papel.

### Mecanismo 2: NPCs migrantes preenchem vagas

Quando há vaga aberta em uma vila e há NPCs vagantes (sem vila) com o papel correto, o NPC migra para a vila e ocupa a vaga.

NPCs vagantes existem por:

- Sobreviventes de vilas destruídas
- Geração espontânea quando população cai abaixo de threshold
- Trânsito normal (NPCs eventualmente se deslocam)

Migração não é instantânea. NPC vagante precisa de tempo de jogo para "viajar" até a vila. Durante esse tempo, ele leva fofocas e conhecimento entre regiões.

### Mecanismo 3: Quests dos jogadores aceleram crescimento

Para uma vila subir de tier, é necessário acumular **pontos de prosperidade**.

Pontos vêm de:

- Quests de jogadores (fonte principal e mais rápida)
- Eventos positivos (boas colheitas, festivais bem-sucedidos)
- Tempo natural com vagas todas preenchidas (lento, mas existe)

Pontos são perdidos por:

- Vagas não preenchidas por muito tempo
- Eventos negativos (ataques sobreviventes, desastres)
- Decay temporal natural se sem progresso

Quando atinge threshold do próximo tier, vila evolui. Quando cai abaixo de threshold do tier atual, decai.

### Mecanismo 4: Decay temporal sem player

Vilas perdem tier automaticamente em duas situações:

- **Vagas vazias por longo tempo**: vila sem ferreiro por 3 estações perde a "produtividade" daquele papel. Continua perdendo tier até voltar a ser apenas tier que comporta vagas atuais.

- **Tempo sem progresso**: mesmo com vagas preenchidas, se nenhuma quest é feita, sociedade estagna e decai lentamente.

Esse mecanismo garante que o mundo não fica congelado. Vilas sem atenção definham. Isso cria pressão narrativa para jogadores apoiarem vilas que querem ver crescer.

## Decisões em aberto (precisam ser resolvidas antes da implementação)

### Decisão A: Vagas predefinidas no mapa ou dinâmicas?

**Opções:**

- A1) Slots de vila são definidos no JSON da região (ex: Mata Costeira tem 9 zonas, 5 podem hospedar vila)
- A2) Vilas podem nascer em qualquer ponto fértil, com regras dinâmicas

**Recomendação preliminar:** A1. Mais simples, mais previsível, melhor para narrativa pré-escrita do GDD. Migração para A2 fica para temporadas futuras.

### Decisão B: Como se calcula "longo tempo" para decay?

**Opções:**

- B1) Tempo absoluto (ex: 30 dias de jogo)
- B2) Por estações (ex: 1 inverno completo sem ferreiro)
- B3) Por temporada de jogo (alinhado com sistema de Temporadas)

**Recomendação preliminar:** B2. Cria narrativa natural ("a vila perdeu o ferreiro no inverno passado"). Alinha com ciclo natural do mundo.

### Decisão C: Quanto tempo um NPC sobrevivente "vagante" permanece?

**Opções:**

- C1) Permanente (NPC fica vagando até achar vaga)
- C2) Tempo limitado (NPC desaparece se não acha vaga em X tempo)
- C3) Geração reverso (NPC sobrevivente eventualmente "volta para casa" virtualmente, é despawned)

**Recomendação preliminar:** C2. Sobreviventes ficam por algumas estações, depois somem se não acharem vaga. Mantém população controlada.

### Decisão D: Como narrar geração espontânea de NPCs?

Quando população do mundo cai abaixo de threshold, novos NPCs aparecem. Como justificar?

**Opções:**

- D1) "Aparecer do nada" (assumir abstração, não justificar)
- D2) "Imigrantes chegam pelas estradas" com histórias rápidas de origem
- D3) Sistema completo de imigração (NPCs vêm de regiões fora do mapa, com narrativa)

**Recomendação preliminar:** D2. Equilibra abstração com imersão. Cada novo NPC tem 1-2 frases de backstory ("João chegou de Olinda fugindo de assalto"). Narrativamente coerente sem custo alto.

### Decisão E: Quanto custa em pontos para evoluir tier?

**Opções:**

- E1) Linear (cada tier custa o mesmo)
- E2) Exponencial (tier 2→3 mais caro que 1→2)
- E3) Variável por região (vilas em zonas hostis crescem mais devagar)

**Recomendação preliminar:** E2 com modificadores regionais leves. Tier alto deve ser conquista significativa, não trivial.

## Pré-requisitos técnicos

Para implementar esse sistema, são necessários (em ordem):

1. **Sistema de combate funcional** (Fase 3+)
   - Para ataques de facções e defesa de vilas

2. **Sistema de quests robusto** (Fase 3+)
   - Para implementar quests de prosperidade

3. **StoryManager com eventos globais** (Fase 5+)
   - Para disparar ataques, desastres, eventos de prosperidade

4. **Persistência avançada** (Fase 6+)
   - Vilas precisam ter estado persistido entre reinicializações
   - Histórico de eventos por vila precisa ser salvo

5. **NPC com role/profissão tipadas** (parcialmente implementado)
   - Vagas dependem de roles específicos

6. **Sistema de migração de NPCs** (Fase 4+)
   - NPCs precisam saber se mover entre vilas
   - Movimento real entre blocos é pré-requisito

7. **Decay temporal calibrado**
   - Sistema de tempo do mundo precisa estar consolidado
   - Estações/temporadas precisam estar implementadas

## Entrega mínima viável (versão "simples" futura)

Mesmo na versão futura, recomenda-se entregar em fases:

**Versão 0.1 (mais simples possível)**

- 3-5 vilas pré-definidas no mapa
- Tier fixo no JSON, sem evolução automática
- NPCs respawnam manualmente via comando admin
- Sem migração, sem decay, sem ataques

**Versão 1.0 (sistema básico)**

- Vagas tipadas por role
- NPCs sobreviventes vagam e ocupam vagas
- Tier evolui via pontos de prosperidade (lento, sem quests)
- Decay simples (perde tier se vagas vazias)

**Versão 2.0 (sistema completo)**

- Quests de jogadores aceleram crescimento
- Eventos de ataque podem destruir vilas
- Múltiplas vilas em diferentes regiões
- Migração rica com fofoca

## Riscos identificados

**Risco 1: Calibração extremamente difícil**

Os 4 mecanismos interagem entre si. Calibrar pontos de prosperidade, decay rate, custo de quests, frequência de NPCs vagantes — tudo precisa ser balanceado. Provavelmente vai exigir muitos meses de tuning após implementação inicial.

**Risco 2: Pode ficar "morto" em servidores com poucos jogadores**

Se há poucos jogadores fazendo quests, vilas vão decair mais que crescer. Pode ser desmotivador. Solução: ajustar para que decay seja lento o suficiente para servidores pequenos não definharem.

**Risco 3: Pode ser invisível para o jogador casual**

Mudanças graduais de tier podem passar despercebidas. Solução: HUD ou interfaces que mostram explicitamente o estado atual da vila, histórico de mudanças, quem contribuiu.

**Risco 4: Complexidade de UI/UX**

Mostrar tier, vagas, prosperidade, decay para o jogador de forma compreensível é desafio significativo. Pode exigir várias iterações de design de interface.

## Como conecta com outros sistemas previstos

- **Sistema de facções**: facções podem atacar vilas adversárias
- **Sistema de inimigos evolutivos**: mobs que matam jogadores ganham notoriedade que pode levar a ataques contra vilas
- **Sistema de temporadas**: cada temporada pode ter eventos que afetam todas as vilas simultaneamente
- **Sistema de fofoca/knowledgeBase**: NPCs migrantes carregam histórias entre vilas, espalhando lendas

## Quando reavaliar este documento

Reabrir este documento quando:

1. Fase 3 (Player) estiver concluída
2. Combate básico estiver funcional
3. Sistema de quests estiver maduro
4. Houver pelo menos 3 meses de jogadores reais usando o servidor

Aí sim, com mais experiência, decidir se vale perseguir esse sistema completo, fazer versão simplificada, ou descartar.

## Histórico de revisões

| Data | Mudança |
| ------ | --------- |
| 2026-04-29 | Criado durante conversa de design sobre sistemas de longo prazo |
