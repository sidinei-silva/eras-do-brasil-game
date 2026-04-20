# Template de Refinamento — [NOME DO ITEM]

> Preencher **antes** de começar a codar. Se não consegue preencher, o item não está pronto pra implementação — busca no GDD ou traz pro designer.

## 1. O que é

Uma ou duas frases descrevendo o que o item entrega, em voz de produto (não de implementação).

> Exemplo: "NPCs conseguem decidir o que fazer baseado em várias necessidades competindo, não mais em thresholds isolados."

## 2. Referência no GDD

Liste os arquivos do GDD que descrevem isso. Se a referência é ambígua ou contradiz, aponte o conflito.

- `gdd/01_Livro_de_Regras/08_Mundo_Vivo_e_NPCs.md` §8.3
- ...

Conflitos encontrados:

- (preencher se tiver)

## 3. Decisões em aberto

Liste os pontos onde o GDD não decide e **você precisa decidir antes do código**. Para cada um, registre a opção escolhida e a justificativa.

**Decisão 3.1 — [tema]**

- Opções: A, B, C
- Escolhido: B
- Justificativa: [por quê]
- Consequência: [o que essa escolha implica que você NÃO vai fazer agora]

## 4. Escopo — o que ENTRA

Lista bullet do que o item inclui. Seja específico.

## 5. Escopo — o que NÃO entra

Lista bullet do que fica pra depois. Isso é tão importante quanto o que entra — é o que te impede de expandir escopo no meio do código.

## 6. Critérios de aceite

Lista checkável do que precisa estar funcionando pra considerar o item pronto.

- [ ] ...
- [ ] ...

## 7. Invariantes e edge cases

O que não pode acontecer nunca? Quais os casos limites?

- Ex: "Nenhum NPC pode ter score negativo de meta"
- Ex: "Se todas as necessidades estão em 0, o NPC escolhe Idle"

## 8. Como validar

Passo-a-passo pra verificar que funciona. Não precisa ser teste automatizado — pode ser "subo o servidor, deixo rodar 5 minutos, vou no admin e confiro X".

## 9. O que isso DESTRAVA

Quando isso estiver pronto, o que fica possível que não era antes?

## 10. O que isso AINDA não resolve

O débito que você aceita conscientemente.
