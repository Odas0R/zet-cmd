# Proposta Equipa Cmd-Zet

## Lista de Inputs e Outputs de Formas Aleatórias

|              | Inputs | Outputs |
| ------------ | ------ | ------- |
| De controlo  |        |         |
| Algoritmicos |        |         |

### Flux Sequence Diagram E.g: Creating, Linking & Stats

This sequence diagram illustrates the process flow for creating zettels,
linking them, and viewing stats. It highlights the interactions between the
user, controllers, models, and views, showcasing the dynamic nature of the
Zettelkasten's operations

```mermaid
sequenceDiagram
    participant Utilizador
    participant VC as ControladorView
    participant ZC as ControladorZettel
    participant Z as Zettel (Modelo)
    participant L as Link (Modelo)
    participant SM as Estatísticas (Modelo)
    participant H as Histórico (Modelo)
    participant P as Pesquisa (Modelo)

    # Criar Zettel
    Utilizador->>+ZC: Solicitar (criar Zettel)
    activate VC
    VC->>+Utilizador: Mostrar Formulário Criar Zettel
    deactivate VC
    Utilizador->>+VC: Submeter Criar Zettel com Título & Conteúdo
    activate VC
    VC->>+ZC: Encaminhar Criar Zettel
    deactivate VC
    ZC->>+Z: Criar Zettel
    Z->>SM: Incrementar Estatísticas Criação
    Z-->>-ZC: Zettel Criado
    activate VC
    ZC-->>VC: Atualizar View com Zettel
    VC-->>-Utilizador: Mostrar Zettel Criado
    deactivate VC

    # Visualizar Estatísticas
    Utilizador->>+ZC: Solicitar (ver estatísticas)
    ZC->>+SM: Obter Estatísticas
    SM-->>-ZC: Dados Estatísticas
    ZC-->>-Utilizador: Resposta (Estatísticas)

    # Criar Link
    Utilizador->>+ZC: Solicitar (ligar Zettels)
    activate VC
    VC->>+Utilizador: Selecionar Zettels para Ligar
    deactivate VC
    Utilizador->>+VC: Submeter Seleção Link (IDs Zettel)
    activate VC
    VC->>+ZC: Encaminhar Ligar Zettels
    deactivate VC
    ZC->>+Z: Criar Link
    Z-->>-ZC: Link Criado
    ZC-->>-Utilizador: Resposta (Zettels ligados)

    # Visualizar Histórico
    Utilizador->>+ZC: Solicitar (ver histórico)
    ZC->>+H: Obter Histórico
    H-->>-ZC: Dados Histórico
    ZC-->>-Utilizador: Resposta (Histórico)

    # Pesquisar
    Utilizador->>+ZC: Solicitar (pesquisar Zettels)
    activate VC
    VC->>+Utilizador: Mostrar Barra de Pesquisa
    deactivate VC
    Utilizador->>+VC: Submeter Consulta de Pesquisa
    activate VC
    VC->>+ZC: Encaminhar Consulta de Pesquisa
    deactivate VC
    ZC->>+P: Realizar Pesquisa
    P-->>-ZC: Resultados da Pesquisa
    ZC-->>-Utilizador: Resposta (Resultados da Pesquisa)

    # Editar Zettel
    Utilizador->>+ZC: Solicitar (editar Zettel)
    activate VC
    VC->>+Utilizador: Mostrar Formulário Editar Zettel
    deactivate VC
    Utilizador->>+VC: Submeter Editar Zettel com Alterações
    activate VC
    VC->>+ZC: Encaminhar Editar Zettel
    deactivate VC
    ZC->>+Z: Editar Zettel
    Z->>SM: Atualizar Estatísticas Modificação
    Z-->>-ZC: Zettel Editado
    activate VC
    ZC-->>VC: Atualizar View com Zettel
    VC-->>-Utilizador: Mostrar Zettel Editado
    deactivate VC

    # Remover Zettel
    Utilizador->>+ZC: Solicitar (remover Zettel)
    activate VC
    VC->>+Utilizador: Confirmar Remoção Zettel
    deactivate VC
    Utilizador->>+VC: Submeter Confirmação Remoção
    activate VC
    VC->>+ZC: Encaminhar Remover Zettel
    deactivate VC
    ZC->>+Z: Remover Zettel
    Z->>SM: Atualizar Estatísticas Remoção
    Z-->>-ZC: Zettel Removido
    activate VC
    ZC-->>VC: Atualizar View
    VC-->>-Utilizador: Mostrar Zettel Removido
    deactivate VC



```

### Flux Diagram MVC

Showcases the zet-cmd system's architecture, detailing the flow from user
interactions through views—like List, Detail, and Edit Views—to controllers and
models for operations such as zettel creation, linking, and stats viewing. It
highlights the system's dynamic interaction between user requests, data
processing, and the subsequent update of views, illustrating the seamless
integration of components for efficient data management and responsive user
experience.

```mermaid
flowchart TD
    User("User") -->|Interact| V{"Views"}

    subgraph ViewLayer
    V --> LV("List View")
    V --> DV("Detail View")
    V --> EV("Edit View")
    V --> SV("Search Interface")
    V --> HV("History View")
    end

    subgraph Controllers
    C["Controllers\n(ZettelController, LinkController, \nSearchController, HistoryController)"] --> ZC("ZettelController")
    C --> LC("LinkController")
    C --> SC("SearchController")
    C --> HC("HistoryController")
    end

    LV -->|Request List| ZC
    DV -->|Request Detail| ZC
    EV -->|Create/Update Request| ZC
    SV -->|Search Request| SC
    HV -->|History Request| HC

    subgraph Models
    direction LR
    Z("Zettel\nModel") -->|Notify| ZC
    L("Link\nModel") -->|Notify| LC
    S("Search\nModel") -->|Notify| SC
    H("History\nModel") -->|Notify| HC
    SM("Stats\nModel") -.->|Notify Stats Update| C
    end

    ZC -.->|Render List Response| LV
    ZC -.->|Render Detail Response| DV
    ZC -.->|Render Edit Response| EV
    SC -.->|Render Search Results| SV
    HC -.->|Render History Data| HV

    LV -->|Show Zettels List| User
    DV -->|Show Zettel Details| User
    EV -->|Show Zettel Edit Form| User
    SV -->|Show Search Results| User
    HV -->|Show Viewed History| User

    ZC -->|Create/Update/Delete Zettel| Z
    LC -->|Manage Links| L
    SC -->|Search Operation| S
    HC -->|Access History| H
    ZC -->|Increment Stats on Create| SM
    LC -->|Update Stats on Link| SM
    HC -->|Update Stats on History Access| SM
    SC -->|Update Stats on Search| SM


```

### Notes

It's important to note that the `Flux Sequence Diagram E.g: Creating, Linking &
Stats` is an example, and the models are not 100% complete since in development
implementations might change. This is an initial sketch and overview of the
`zet-cmd` architecture.

### Repository

1. <https://github.com/Odas0R/zet-cmd>
2. <https://github.com/Odas0R/zet-cmd/blob/main/ARCHITECTURE.md>
