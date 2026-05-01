# North Star — Eras do Brasil

> Documento de visão fundamental do jogo. Quando estiver dividido entre opções de design, volte aqui.

## Identidade do jogo

**Eras do Brasil** é um MMORPG onde o mundo é o centro, não o jogador.

### O pilar fundamental

> "O player não é o ponto central, e sim o mundo. O player é apenas mais um no mundo. Ele organicamente vira herói porque todos o veem como herói."

Isso significa três coisas práticas:

1. **O mundo continua vivendo independente do jogador.** NPCs têm rotinas próprias. Vilas evoluem ou decaem. Eventos acontecem com ou sem participação dele.

2. **O jogador conquista importância pelos atos.** Ninguém o chama de "Escolhido" ao chegar. Sua reputação é construída através de ações observadas e narradas pelos NPCs e pela comunidade.

3. **A história é emergente, não roteirizada.** As lendas do servidor surgem do que os jogadores fazem, não de cenas pré-escritas. "A vila de São Tomé foi defendida pela aventureira Maria no inverno de 1503" é história real do servidor, não enredo do GDD.

## Player Fantasy emergente vs declarada

Esta distinção guia todas as decisões de design:

**RPGs tradicionais usam fantasy declarada:**

- "Você é o Escolhido das profecias"
- "Apenas você pode salvar o reino"
- "Os deuses te selecionaram"

**Eras do Brasil usa fantasy emergente:**

- "Você ajudou a vila a sobreviver ao ataque dos Bandeirantes"
- "Foi você quem trouxe o Ferreiro de volta após o exílio"
- "Os NPCs falam de você porque suas ações tiveram consequências reais"

Quando estiver dividido entre duas opções de design, pergunte:

> Qual dessas alternativas reforça a fantasy emergente? Qual transforma o jogador em mais um personagem dentro do mundo, em vez de protagonista predestinado?

A resposta certa quase sempre é óbvia uma vez que essa pergunta é feita.

## Implicações de design

### Sobre NPCs

NPCs não esperam o jogador. Eles têm rotinas, necessidades, relacionamentos próprios. Quando o jogador interage, é uma interrupção bem-vinda mas não esperada. Quando o jogador não está, a vida segue.

### Sobre quests

Quests não são "tarefas para você fazer". São **necessidades do mundo** que o jogador pode ajudar a resolver. Se o jogador não fizer, ou outro jogador faz, ou o mundo se vira sem ele (com consequências), ou a necessidade muda de forma.

### Sobre progressão

O jogador não progride por "level up tradicional acompanhando linha de quests". Ele progride pela acumulação de:

- Habilidades praticadas
- Reputação construída
- Itens conquistados em eventos
- Memória dos NPCs sobre suas ações
- Lugar nos rankings de temporada

Isso é progressão **observável pelo mundo**, não progressão **interna ao jogador**.

### Sobre mortes e consequências

Eventos no mundo têm consequências reais. NPCs morrem permanentemente. Vilas podem cair. Recursos podem se esgotar. Estações podem ser perdidas. O jogador não tem garantia de que o mundo está esperando ele estar pronto. Isso cria peso para suas decisões.

### Sobre multiplayer

Múltiplos jogadores no mesmo servidor compartilham um mundo único. Diferentes jogadores podem se especializar em diferentes vilas, regiões ou facções. Não existe "main quest linear" que todos seguem. Existe um **mundo persistente** onde cada um encontra seu papel.

## O que isso NÃO é

Para reforçar pelo contraste:

- **Não é Sandbox puro.** Há narrativa, há StoryManager, há temporadas com arcos pensados. O mundo tem direção, não é caos puro.
- **Não é Theme Park MMO.** Não há fila de quests. Não há "vai falar com o NPC pra começar a saga". O jogador descobre, observa, decide.
- **Não é jogo single-player coop.** O servidor não escala para um jogador. Vilas são compartilhadas. Eventos afetam todos. Trocar é mais importante que competir.
- **Não é simulação científica de mundo.** É um jogo. Há simplificações. NPCs nascem prontos como adultos (geração espontânea via "imigração"). Não há genética. Não há infância. Mas as simplificações respeitam o pilar fundamental.

## Decisões já tomadas alinhadas com esse North Star

- **NPCs com rotinas próprias** (Fase 1 — Mundo Vivo): cada NPC vive seu dia independente do jogador
- **Schedule + Necessidades emergentes**: NPC interrompe trabalho quando tem fome, sem precisar do jogador
- **Sleep emerge da fadiga**, não é horário fixo: comportamento autônomo
- **Tick contínuo do servidor**: mundo nunca pausa
- **Servidor único persistente** (sem instâncias): todos compartilham o mesmo mundo

## Decisões previstas alinhadas com esse North Star (visão futura)

- **Sistema de vilas dinâmico** (ver `sistema-vilas-dinamico.md`): vilas evoluem e decaem independente de jogadores, com aceleração via quest
- **Inimigos evolutivos** (ADR-009): mobs ganham XP por matar jogadores, viram lendas do servidor
- **Temporadas com legados**: ações em uma temporada deixam marcas observáveis nas seguintes
- **Knowledge base nos NPCs**: NPCs lembram do que fizeram, viram fofocas que se espalham

## Quando este documento se aplica

Sempre que você estiver:

- Decidindo entre duas alternativas de mecânica
- Avaliando proposta de feature nova
- Sentindo que o jogo está "perdendo identidade"
- Sendo tentado por feature genérica de MMORPG (auto-quests, mini-mapa com setinhas, etc.)
- Comunicando o jogo para outras pessoas

Volte aqui. Pergunte: isso reforça ou enfraquece o pilar fundamental?

Se enfraquece, recuse mesmo que pareça boa ideia.
Se reforça, vale considerar mesmo que pareça difícil de implementar.

## Histórico de revisões

| Data | Mudança |
| ------ | --------- |
| 2026-04-29 | Criado. Pilar fundamental articulado durante conversa de design sobre sistema de vilas dinâmico. |
