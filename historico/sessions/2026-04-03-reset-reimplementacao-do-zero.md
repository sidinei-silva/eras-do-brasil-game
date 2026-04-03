# 2026-04-03 — Reset: Reimplementação do Zero

## Contexto

O desenvolvedor identificou que perdeu a capacidade de abstração e raciocínio lógico por ter delegado demais para IA. O código do servidor funcionava, mas ele não reconhecia a lógica, não entendia as decisões, e não conseguia se sentir dono do projeto.

## Diagnóstico

- **Atrofia de abstração:** a IA fazia a parte difícil (decompor problemas, criar estruturas), e o músculo de engenharia parou de ser exercitado.
- **Barreira não é sintaxe:** o problema não é Go — é a construção da lógica, o pensamento de engenheiro, a capacidade de abstrair. Isso afeta até JS/TS.
- **Perda de reconhecimento:** código gerado por IA não era reconhecível. Ao ler, não conseguia entender o "porquê" das decisões.

## Decisões

1. **Apagar todo o código Go do servidor** (commit `4d64406` como backup para "pescar")
2. **Reimplementar do zero, sozinho** — pode começar tudo em 1 arquivo só, sem estrutura de pacotes prematura
3. **IA como PM/PD, não como dev** — guia teórico, sem dar código. Ajuda a pensar, não a escrever.
4. **Remover comentários educativos** — o dev vai reescrever com suas próprias palavras ao estudar
5. **Limpeza de documentação** — muitos arquivos redundantes, simplificar para o essencial
6. **Backlog simplificado** — lista O QUE precisa funcionar, não COMO implementar. O dev decide a ordem e a estrutura.

## Abordagem Nova

- O dev decide a ordem de implementação (exercita abstração)
- Quando travar, pesca do commit de backup
- Sem pressão de arquitetura — faz funcionar primeiro, organiza depois
- Sessões de estudo: ler código, comentar com suas palavras, testar

## Decisao: Go (basico) desde o inicio

Cogitou fazer em JS/TS primeiro para isolar o problema de abstracao. Decisao final: ficar com Go usando apenas o subconjunto basico (structs, funcoes, if/for/switch). Sem goroutines, channels, interfaces ou pacotes separados ate precisar. Tudo no main.go.

Motivo: o problema e decomposicao de problemas, nao sintaxe. Trocar de linguagem nao resolve, e refazer em duas linguagens custa o dobro do tempo.

## Ordem de desenvolvimento (metafora da criacao)

1. Mundo (o container)
2. Tempo passando (dias/tick)
3. Locais (mapa/zonas)
4. Recursos (natureza)
5. Animais/mobs
6. NPCs
7. Rotinas e necessidades dos NPCs
8. Conhecimento dos NPCs
9. Mundo funcionando sozinho
10. Player entra no mundo ja vivo

Filosofia: o mundo e o protagonista, nao o jogador. Construir de baixo pra cima ate o mundo rodar sozinho, so entao incluir o player.

## Arquivos Afetados

- Servidor resetado para Hello World (`server/main.go`)
- Código anterior preservado no git (commit `4d64406`)
- Docs limpos: removidos GUIA_RETOMADA, perfil-developer, plano-desenvolvimento
- Backlog reescrito (foco em critérios, não em passos)
- Pasta estudos/ limpa (removidas referências a código apagado)
