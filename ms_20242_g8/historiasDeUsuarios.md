## **História de Usuário 1**

### **ID:**  
HU-001

### **Título:**  
Cadastro de Usuários

### **Descrição:**  
Como usuário, eu quero poder me **cadastrar** no sistema com meu nome, e-mail e senha para que eu possa participar do marketplace de trocas de livros.

### **Critérios de Aceitação:**

#### Cenário 1: Cadastro com sucesso
- **Dado** que o usuário está na página de cadastro,  
- **Quando** o usuário insere seu nome, e-mail e senha válidos,  
- **Então** o sistema deve criar uma nova conta e permitir que o usuário acesse o sistema.

#### Cenário 2: Falha no cadastro por e-mail duplicado
- **Dado** que o usuário está tentando se cadastrar,  
- **Quando** o e-mail informado já estiver em uso por outra conta,  
- **Então** o sistema deve exibir uma mensagem de erro informando que o e-mail já está cadastrado.

### **Prioridade:**  
Alta

### **Estimativa de Esforço:**  
5 Story Points

### **Dependências:**  
- Integração com o banco de dados de usuários.  
- Desenvolvimento da interface de cadastro.

### **Notas/Comentários Adicionais:**  
- O sistema deve validar a formatação do e-mail e a força da senha.

__
__

## **História de Usuário 2**

### **ID:**
HU - 002

### **Título:**
Login/Logout

### **Descrição:**
Como usuário cadastrado, eu quero fazer login e logout para que eu possa acessar e gerenciar meu perfil no marketplace.

### **Critérios de aceitação:**

#### Cenário 1: Login com sucesso
- **Dado** que o usuário está na página de login,
- **Quando** o usuário insere seu e-mail e senha corretos,
- **Então** o sistema deve autenticar o usuário e direcioná-lo para a página inicial.

#### Cenário 2: Logout com sucesso:
- **Dado** que o usuário está logado,
- **Quando** o usuário seleciona a opção de logout,
- **Então** o sistema deve deslogar o usuário e redirecioná-lo para página de login.

### **Prioridade**
Alta

### **Estimativa de Esforço:**
5 Story Points

### **Dependências:**
- Integração com o banco de dados de autenticação.
- Desenvolvimento da interface de login/logout.

### **Notas/Comentários adicionais:**
- O sistema deve validar o e-mail e a senha durante o login.

__
__

## **História de Usuário 3**

### **ID:**  
HU-003

### **Título:**  
Adicionar Livros ao Marketplace

### **Descrição:**  
Como **usuário**, eu quero **adicionar meus livros ao marketplace** para que **outros usuários possam vê-los e escolham para trocas**.

### **Critérios de Aceitação:**

#### Cenário 1: Adicionar livro com sucesso
- **Dado** que o usuário está logado e acessa a página de adicionar livros,  
- **Quando** o usuário insere as informações do livro (nome, autor, edição, foto),  
- **Então** o livro deve ser adicionado ao marketplace e exibido na lista de livros disponíveis.

#### Cenário 2: Falha ao adicionar livro por informações incompletas
- **Dado** que o usuário está logado e acessa a página de adicionar livros,  
- **Quando** o usuário não insere todas as informações obrigatórias (nome, autor, edição),  
- **Então** o sistema deve exibir uma mensagem de erro informando que os campos obrigatórios precisam ser preenchidos.

#### Cenário 3: Falha ao adicionar livro por erro no upload da foto
- **Dado** que o usuário está logado e tenta fazer o upload de uma foto para o livro,  
- **Quando** o upload falhar por formato ou tamanho da imagem inválidos,  
- **Então** o sistema deve exibir uma mensagem de erro explicando a falha e solicitando que o usuário tente novamente com uma imagem válida.

#### Cenário 4: Verificação de duplicidade de livro
- **Dado** que o usuário está adicionando um livro ao marketplace,  
- **Quando** um livro com as mesmas informações (nome, autor, edição) já estiver listado no marketplace,  
- **Então** o sistema deve alertar o usuário de que o livro já está disponível e permitir ou solicitar que ele revise suas entradas.

### **Prioridade:**  
Média

### **Estimativa de Esforço:**  
8 Story Points

### **Dependências:**  
- Integração com o banco de dados para salvar as informações do livro.  
- Desenvolvimento da interface de upload de livros.

### **Notas/Comentários Adicionais:**  
- O sistema deve permitir que o usuário faça upload de uma foto do livro.  
- Verificar o tamanho máximo permitido para a imagem e os formatos aceitos (JPEG, PNG).

__
__

## **História de Usuário 4**

### **ID:**  
HU-004

### **Título:**  
Explorar Livros Disponíveis

### **Descrição:**  
Como **usuário**, eu quero **explorar todos os livros disponíveis no marketplace** para que **eu possa escolher um livro de meu interesse para trocar**.

### **Critérios de Aceitação:**

#### Cenário 1: Visualização da lista de livros
- **Dado** que o usuário está na página de livros disponíveis,  
- **Quando** o usuário acessa a página,  
- **Então** todos os livros disponíveis devem ser exibidos na lista.

#### Cenário 2: Filtro de livros por categoria
- **Dado** que o usuário está na página de livros disponíveis,  
- **Quando** o usuário utiliza o filtro de categoria,  
- **Então** os livros correspondentes à categoria selecionada devem ser exibidos.

#### Cenário 3: Filtro de livros por autor
- **Dado** que o usuário está na página de livros disponíveis,  
- **Quando** o usuário utiliza o filtro de autor,  
- **Então** os livros correspondentes ao autor selecionado devem ser exibidos.

#### Cenário 4: Filtro de livros por título
- **Dado** que o usuário está na página de livros disponíveis,  
- **Quando** o usuário utiliza o filtro de título,  
- **Então** os livros correspondentes ao título digitado devem ser exibidos.

#### Cenário 5: Nenhum livro encontrado após aplicação de filtros
- **Dado** que o usuário está na página de livros disponíveis,  
- **Quando** o usuário aplica um filtro que não encontra correspondências (categoria, autor ou título),  
- **Então** o sistema deve exibir uma mensagem informando que não há livros correspondentes aos critérios selecionados.

### **Prioridade:**  
Alta

### **Estimativa de Esforço:**  
8 Story Points

### **Dependências:**  
- Desenvolvimento da interface de busca e filtros.  
- Integração com o banco de dados para recuperar os livros disponíveis.

### **Notas/Comentários Adicionais:**  
- O sistema deve permitir a filtragem por categoria, autor e título.
- Deve haver uma opção para remover os filtros aplicados e visualizar todos os livros novamente.

__
__

## **História de Usuário 5**

### **ID:**  
HU-005

### **Título:**  
Marcar Interesse em Troca

### **Descrição:**  
Como **usuário**, eu quero **marcar interesse em trocar um livro específico** para que **o dono do livro seja notificado sobre meu interesse**.

### **Critérios de Aceitação:**

#### Cenário 1: Marcar interesse com sucesso
- **Dado** que o usuário está logado e visualizando um livro no marketplace,  
- **Quando** o usuário marca interesse em um livro,  
- **Então** o dono do livro deve ser notificado sobre o interesse.

#### Cenário 2: Troca efetivada após acordo
- **Dado** que o usuário e o dono do livro estão interessados em realizar a troca,  
- **Quando** ambos concordam com a troca,  
- **Então** a troca deve ser efetivada, e os livros devem ser removidos do marketplace.

#### Cenário 3: Falha ao marcar interesse sem estar logado
- **Dado** que o usuário está visualizando um livro no marketplace sem estar logado,  
- **Quando** o usuário tenta marcar interesse,  
- **Então** o sistema deve redirecionar o usuário para a página de login e exibir uma mensagem informando que é necessário estar logado para marcar interesse.

#### Cenário 4: Falha ao marcar interesse em livro já trocado
- **Dado** que o usuário está visualizando um livro no marketplace,  
- **Quando** o livro já foi trocado mas ainda aparece listado,  
- **Então** o sistema deve informar que o livro não está mais disponível para troca.

### **Prioridade:**  
Alta

### **Estimativa de Esforço:**  
8 Story Points

### **Dependências:**  
- Sistema de notificação entre usuários.  
- Regras de negócio para efetivação de troca.

### **Notas/Comentários Adicionais:**  
- O sistema deve registrar o interesse e permitir que o dono do livro aceite ou recuse a troca.
- A troca só deve ser efetivada quando ambos os usuários confirmarem o acordo.

__
__

## **História de Usuário 6**

### **ID:**  
HU-006

### **Título:**  
Gerenciar Livros no Perfil

### **Descrição:**  
Como **usuário**, eu quero **gerenciar os livros que disponibilizei no marketplace**, podendo **adicionar, editar ou remover livros** quando desejar.

### **Critérios de Aceitação:**

#### Cenário 1: Editar informações de um livro com sucesso
- **Dado** que o usuário está logado e acessa a página de gerenciamento de livros,  
- **Quando** o usuário edita as informações de um livro (nome, autor, edição, foto),  
- **Então** o livro deve ser atualizado no marketplace.

#### Cenário 2: Remover livro do marketplace com sucesso
- **Dado** que o usuário está logado e acessa a página de gerenciamento de livros,  
- **Quando** o usuário remove um livro,  
- **Então** o livro deve ser removido do marketplace e não deve mais aparecer na lista de livros disponíveis.

#### Cenário 3: Falha ao editar informações por dados incompletos
- **Dado** que o usuário está logado e acessa a página de gerenciamento de livros,  
- **Quando** o usuário tenta editar um livro sem fornecer todas as informações obrigatórias,  
- **Então** o sistema deve exibir uma mensagem de erro informando que os campos obrigatórios precisam ser preenchidos.

#### Cenário 4: Visualizar lista de livros adicionados no perfil
- **Dado** que o usuário está logado e acessa a página de gerenciamento de livros,  
- **Quando** o usuário visualiza sua lista de livros,  
- **Então** todos os livros que ele adicionou ao marketplace devem ser exibidos, com a opção de editar ou remover cada um.

### **Prioridade:**  
Média

### **Estimativa de Esforço:**  
5 Story Points

### **Dependências:**  
- Integração com o banco de dados para edição e remoção de livros.

### **Notas/Comentários Adicionais:**  
- O sistema deve permitir ao usuário visualizar a lista de livros adicionados e gerenciar cada um deles individualmente.

__
__

## **História de Usuário 7**

### **ID:**  
HU-007

### **Título:**  
Concluir Troca

### **Descrição:**  
Como **usuário**, eu quero **concluir uma troca após aceitar o interesse de outro usuário no meu livro** para que **a troca seja efetivada no sistema**.

### **Critérios de Aceitação:**

#### Cenário 1: Concluir troca com sucesso
- **Dado** que o usuário e o outro participante da troca chegaram a um acordo,  
- **Quando** ambos os usuários confirmam a troca,  
- **Então** a troca deve ser concluída, e o livro deve ser removido do marketplace.

#### Cenário 2: Falha na conclusão da troca por falta de confirmação de uma das partes
- **Dado** que o usuário está aguardando a confirmação da troca,  
- **Quando** uma das partes ainda não confirmou a troca,  
- **Então** o sistema deve manter o livro no marketplace e notificar o usuário que a troca está pendente.

#### Cenário 3: Troca cancelada por uma das partes
- **Dado** que o usuário e o outro participante ainda não concluíram a troca,  
- **Quando** uma das partes decide cancelar o processo de troca,  
- **Então** o sistema deve cancelar a troca e notificar ambas as partes.

#### Cenário 4: Notificação de troca concluída
- **Dado** que a troca foi confirmada por ambas as partes,  
- **Quando** a troca for efetivada,  
- **Então** o sistema deve enviar uma notificação para ambos os usuários confirmando que a troca foi concluída.

### **Prioridade:**  
Alta

### **Estimativa de Esforço:**  
8 Story Points

### **Dependências:**  
- Sistema de gerenciamento de trocas.  
- Interface de confirmação da troca por ambas as partes.

### **Notas/Comentários Adicionais:**  
- O sistema deve enviar uma notificação aos usuários confirmando que a troca foi concluída.
- As trocas só devem ser finalizadas quando ambas as partes confirmarem o acordo.

__
__

## **História de Usuário 8**

### **ID:**  
HU-008

### **Título:**  
Histórico de Trocas

### **Descrição:**  
Como **usuário**, eu quero **visualizar um histórico das minhas trocas realizadas** para que **eu possa acompanhar todas as transações que já concluí no marketplace**.

### **Critérios de Aceitação:**

#### Cenário 1: Exibir histórico de trocas com sucesso
- **Dado** que o usuário está logado e acessa a página de histórico de trocas,  
- **Quando** o sistema carrega a página,  
- **Então** todas as trocas concluídas pelo usuário devem ser exibidas com detalhes (livro trocado, data da troca).

#### Cenário 2: Histórico vazio
- **Dado** que o usuário está logado e acessa a página de histórico de trocas,  
- **Quando** o usuário ainda não tiver concluído nenhuma troca,  
- **Então** o sistema deve exibir uma mensagem informando que não há trocas realizadas.

#### Cenário 3: Filtro de histórico por data
- **Dado** que o usuário está na página de histórico de trocas,  
- **Quando** o usuário utiliza um filtro de data,  
- **Então** o sistema deve exibir apenas as trocas realizadas no período selecionado.

### **Prioridade:**  
Média

### **Estimativa de Esforço:**  
5 Story Points

### **Dependências:**  
- Integração com o banco de dados para recuperar o histórico de trocas.

### **Notas/Comentários Adicionais:**  
- O histórico deve exibir detalhes como o nome do livro, data da troca, e o nome do usuário com quem a troca foi realizada.
- O sistema deve permitir ao usuário filtrar o histórico por data.





