# Notas sobre o Desenvolvimento

## 10/04/2026

Estou refatorando o websocket para que o pacote de hub não processe os comandos do client em readPump.
Preciso usar o padrão de commands, para que as entradas dos clients va para um lista de commands, esses comandos precisam serem processados fora do hub, assim posso remover o `handleInbound` e `InboundEvent` assim como todos os codigos gerados por eles.

Nesse exato momento o que quero:

[x] comentar tudo sobre admin (ainda não vou refatorar o admin mas vou comentar para servir como base)
[] refatorar o client removendo assim as entradas e processamento do client
[] repassar o que o client envia para um channel command, mas não processar apenas repassar e colocar um log para isso
