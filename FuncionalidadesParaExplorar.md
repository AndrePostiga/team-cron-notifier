
O domínio deste servico é revisar PRs, então a ideia é que ele possa ser configurado para revisar PRs de diferentes repositórios e enviar uma mensagem para o responsável por revisar o PR.
Portanto, a lógica de domínio está centralizada no pacote `vcs`.

### Funcionalidades
- Levar em consideração quando não há PRs para serem revisados e enviar uma mensagem de parabens (Também possibilitar configurar como não enviar caso não tenha PRs)
- Tentar ordenar os PRs enviados por Dias em que ele está aberto
- Adicionar o nome do repositório na mensagem enviada para indicar
- Adicionar regras de "Aguardando Merge" configuráveis