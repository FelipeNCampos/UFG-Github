## ms_20242_g8
Repositório definido para a manutenção do controle de versão dos artefatos do projeto do Grupo 8, da Disciplina de Modelagem de Software, no semestre 2024-2.

### Nome do Projeto:
Troca Livro

### Descrição:
Este projeto visa desenvolver uma plataforma online voltada para leitores, promovendo a troca de livros já lidos por outros de interesse de forma colaborativa. A plataforma terá como público-alvo pessoas que buscam expandir suas leituras de maneira prática e econômica, oferecendo recursos como cadastro, login, exploração de livros disponíveis e um mecanismo eficiente de troca. Além disso, incluirá funcionalidades como histórico e avaliações de trocas para garantir uma experiência segura e confiável.
A monetização será baseada em assinaturas e serviços pagos, proporcionando aos usuários opções exclusivas para uma experiência de troca e descoberta literária ainda mais personalizada e enriquecedora. Dessa forma, o projeto busca unir praticidade e interação na troca de livros, enquanto oferece serviços que agregam valor à jornada dos leitores.

### Problema
Atualmente, muitos leitores acumulam livros lidos que, embora tenham cumprido seu papel inicial de entretenimento ou aprendizado, acabam tendo pouco ou nenhum impacto contínuo na evolução dos interesses ou da ideologia de seus donos. Esse acúmulo pode limitar o espaço para novas obras que realmente inspirem reflexões e ampliação de horizontes. Uma plataforma de troca de livros permitiria aos leitores renovar constantemente suas estantes, trocando obras já lidas por outras de interesse, promovendo um fluxo contínuo de conhecimento e de novas perspectivas.

### Objetivos da Solução:
O objetivo é criar um sistema que facilite a rotatividade de livros, permitindo trocas entre leitores e promovendo uma constante renovação das estantes. Além disso, busca-se incentivar o crescimento pessoal e intelectual dos usuários, proporcionando acesso a novos temas e perspectivas. Ao sugerir títulos alinhados às preferências individuais, o sistema aumenta o acesso a conteúdos relevantes, enriquecendo a experiência de leitura. A criação de uma comunidade literária interativa, onde leitores possam compartilhar recomendações e experiências, fortalece os laços e o conhecimento coletivo. Algoritmos personalizados oferecerão sugestões baseadas em gosto, estilo e temas, enquanto funcionalidades exclusivas e modelos de assinatura garantirão uma monetização sustentável. Dessa forma, visa-se resolver o acúmulo de livros e enriquecer a experiência de leitura, criando uma dinâmica de troca e evolução contínua.

### Grupo
Este projeto será desenvolvido pelos componentes do grupo 8:

| Matrícula | Nome                             | Usuário Git                                             |
|-----------|----------------------------------|---------------------------------------------------------|
| 201703749 | ERIC DE JESUS MACHADO            | [ericjmachado](https://github.com/ericjmachado)                     |
| 202403069 | FELIPE NUNES CAMPOS      | [Perfil](https://github.com/ufgfelipencampos)                   |
| 202201702 | LUCAS COSTA DE ALVARENGA   | [Lucas](https://github.com/lucascstlvrng) |
| 202305528 | CLAUDIO CRISTIANO LOUZA FILHO | [ClaudioLouza](https://github.com/ClaudioLouza) |
| 201508741 | MARCUS MITRA MUNIZ INACIO    | [MarcusMitra](https://github.com/MarcusMitra) |

## Backlog do Produto

### Requisitos Funcionais

|Código | Título | Descrição |
|--|--|--|
| **RF001** | Cadastro de Usuários | Usuários podem criar contas fornecendo nome, e-mail e senha. O e-mail é único para cada usuário.|
| **RF002** | Login | Usuários cadastrados podem fazer login com e-mail e senha para gerenciar seus perfis. |
| **RF003** | Logout | Usuários logados podem fazer logout de suas contas. |
| **RF004** | Adicionar Livros ao Marketplace | Usuários podem adicionar livros ao Marketplace para trocas fornecendo os seguintes dados do livro: Título; Autor; Edição(Opcional); Editora; Foto(Opcional); Categoria(Opcional).|
| **RF005** | Explorar Livros Disponíveis | Todos os livros disponíveis no Marketplace podem ser visualizados em uma lista, com filtros para categoria, autor, título e editora. |
| **RF006** | Marcar Interesse em Troca | Usuários podem indicar interesse em trocas, marcando um livro e automaticamente notificando o dono. A troca só ocorrerá após a concordância de ambos. |
| **RF007** | Gerenciar Livros no Perfil | Usuários gerenciam seus livros, podendo editar ou remover informações. As seguintes informações não podem ficar em branco: Título; Autor; Editora;|
| **RF008** | Concluir Troca | Trocas são concluídas após aceitação mútua e o livro é removido do marketplace ao finalizar a troca. |
| **RF009** | Histórico de Trocas | Lista todas as trocas realizadas, em aberto e canceladas. |
| **RF011** | Avaliação de Trocas | Após a troca, usuários podem avaliar sua experiência, deixando uma nota entre 1-5 estrelas e/ou um comentário. |

### Requisitos Não Funcionais

|Código | Categoria | Descrição|
|--|--|--|
| **RNF001** | Usabilidade | O sistema deve fornecer feedback claro e coerente com mensagens de erro. |
| **RNF002** | Segurança | O sistema deve autenticar os usuários antes de permitir acesso às funcionalidades. |
| **RNF003** | Desempenho | O sistema deve responder em menos de 2 segundos para a maioria das solicitações. |
| **RNF004** | Confiabilidade | Disponibilidade do software 24/7. |
| **RNF005** | Manutenibilidade| Desenvolvimento cunhado em boas práticas de programação e documentado. |
| **RNF006** | Portabilidade | Suporte mobile. |

### Regras de Negócio

|Código | Categoria | Descrição|
|--|--|--|
| **RN001** | Segurança | Um usuário não pode acessar a conta de outro usuário. |
| **RN002** | Segurança | Um usuário pode apenas se registrar com uma conta válida, ou seja, o sistema deve verificar se as informações fornecidas(nome, e-mail, senha) são válidas.|

### Modelo Arquitetural
Através do Framework Django, seguimos o modelo arquitetural conhecido como **MVC (Model-View-Controller)**, ou, chamado por ele de **MTV (Model-Template-View)**. 

**Fluxo Resumido:**
1.  O cliente faz uma requisição HTTP.
2.  A **View** processa a requisição e decide qual modelo e template usar.
3.  O **Model** interage com o banco de dados para buscar ou atualizar os dados.
4.  O **Template** renderiza os dados em uma resposta HTML (ou outro formato).
5.  A resposta é enviada ao cliente.

### Modelo de Interfaces Gráficas
Estamos utilizando HTML juntamente com o Django, que oferece alguns templates que permitem gerar HTML dinâmico com base nos dados do servidor. Além disso, utilizamos o framework CSS Bootstrap para estilizar as interfaces.

### Tecnologia de Persistência de Dados
Estamos utilizando o SQLite3 juntamente com o Django para fazer a comunicação com o banco de dados. O Django implementa a tecnologia de persistência de dados por meio do **ORM (Object-Relational Mapping)**, que traduz os objetos Python em instruções SQL para interagir com os bancos de dados relacionais. 

-   **Definição dos Modelos**:
    
    -   Os modelos são definidos como classes Python no arquivo models. Sendo assim, cada classe representa uma tabela do banco de dados, e cada atributo da classe representa uma coluna da tabela.

-   **Migrações**:
    
    -   O Django usa o sistema de migrações para traduzir os modelos Python para instruções SQL, criando ou atualizando tabelas no banco de dados.

### Local do _Deploy_
O nosso planejamento inicial é hospedar nossa aplicação no Heroku ou no Firebase do Google, sendo crucial para nossa escolha a facilidade de configuração e também os recursos ofertados na versão gratuita.

### Configuração e Instalação

Para configurar o ambiente local:

1.  **Pré-requisitos**: Certifique-se de ter [Python,Django,asgiref,sqlparse,pillow].
2.  **Passos de Instalação**:
	- Clone o repositório:
		> git clone https://github.com/ufgfelipencampos/ms_20242_g8.git
		
	- Instale as dependências:
		> pip install django asgiref sqlparse pillow
		
	- Configure o banco de dados:
		> python .\trocalivros\manage.py migrate
		
3. **Execução**:
	- python .\trocalivros\manage.py runserver

### Cronograma de Desenvolvimento

| Iteração | Descrição                                |Data Início| Data Fim   |Responsável|Situação|
|----------|------------------------------------------|---|------------|---|---|
| 1        | Concepção                                |10/09/2024| 17/09/2024 |Grupo|Concluída|
| 2        | Preparação                               |24/09/2024| 01/10/2024 |Grupo|Concluída|
| 3        | Definição Inicial do Backlog e Atribuição de Tarefas |08/10/2024| 15/10/2024 |Grupo| Concluída |
| 4        | Detalhamento do Backlog; Estórias de Usuário; Diagrama de Classes; Modelo Arquitetural; Modelo de Interfaces Gráficas; Tecnologia de Persistência de Dados             |22/10/2024| 29/10/2024 |Grupo|Concluída|
| 5        | Diagrama de Sequência; Diagrama de Estado; Diagrama de Atividade; Prototipagem; Implementação Inicial            |05/11/2024| 12/11/2024 |Grupo|Concluída|
| 6        | Revisão e Implementação Final             |12/11/2024| 18/11/2024 |Grupo|Em progresso|
| 7        | Testes; Correções; Configuração e Instalação; Apresentação do Projeto                  |19/11/2024| 26/11/2024 |Grupo|Programada|

* Implementação se aplicará, se os itens da iteração em andamento, forem eleitos para validação do projeto do trabalho.

Usuário github do Professor Gilmar: gilmarUFG.
