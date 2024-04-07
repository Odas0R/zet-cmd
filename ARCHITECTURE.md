# Architecture

Documenting the `zet-cmd` architecture.

## Models

#### Zettel Model

`Zettel` is the central entity , representing individual notes or "zettels".
The model includes essential fields for identifying and managing zettels, as
well as auxiliary fields to enhance functionality, like Lines for content
manipulation and Links for inter-zettel connections.

#### Link Model

The Link model facilitates the connections between zettels, enabling the
creation of a network of related notes. Each link associates two zettels,
indicating a directional or bidirectional relationship.

#### Search Model

The Search model manages the logic and data involved in searching the
Zettelkasten. This model does not correspond directly to a database table but
encapsulates the criteria and algorithms for searching through zettels based on
various parameters.

#### History Model

The History model tracks the user's interaction history within the system, such
as recently viewed or edited zettels. This model helps in providing users with
a personalized experience by enabling quick access to previously interacted
zettels.

#### Stats Model

The Stats model collects and stores statistics related to various operations
within the system, such as the number of zettels created, links made, searches
performed, and historical interactions. It serves as a basis for analyzing the
usage patterns and efficiency of the Zettelkasten.

## Controllers

- `ZettelController`: Endpoint for Zettel management operations to process user
  requests like creation, update, and deletion of zettels.

- `LinkController`: Facilitates operations related to linking zettels,
  including creating and removing links, and querying for backlinks or all
  links related to a zettel.

- `SearchController`: Provides endpoints for searching the Zettelkasten, return
  search results.

- `HistoryController`: Offers access to a user's history within the
  Zettelkasten, such as recently edited or viewed zettels

- `ViewController`: Handles server-side rendering (SSR) of pages, presenting
  the user interface for interacting with the Zettelkasten. This includes
  rendering views for search results, zettel content, and navigation.

## Views

Views are the SSR templates or components that render the Zettelkasten's user
interface:

1. `List View`: Displays search results or collections of zettels, such as
   backlogs or history lists. This view is essential for users to browse
   through multiple zettels efficiently.

2. `Detail View`: Shows individual zettels in detail, including content, metadata
   (like creation and update times), and links/backlinks. This view is crucial
   for reading and understanding the content of a zettel.

3. `Edit View`: Provides a form or interface for creating new zettels or updating
   existing ones. This view is fundamental for content creation and editing
   within the Zettelkasten.

4. `Navigation and Search Interface`: Offers a comprehensive navigation bar or
   search box, enabling users to easily move between different parts of the
   Zettelkasten or to find specific zettels. This interface enhances the
   usability and accessibility of the system.

## Interaction Diagrams

### Flux Sequence Diagram E.g: Creating, Linking & Stats

This sequence diagram illustrates the process flow for creating zettels,
linking them, and viewing stats. It highlights the interactions between the
user, controllers, models, and views, showcasing the dynamic nature of the
Zettelkasten's operations

```mermaid
sequenceDiagram
    participant User
    participant VC as ViewController
    participant ZC as ZettelController
    participant Z as Zettel (Model)
    participant SM as Stats (Model)

    User->>+ZC: Request (create Zettel)
    ZC->>+Z: Create Zettel
    Z-->>-ZC: Zettel Created
    ZC-->>-User: Response (Zettel created)

    User->>+ZC: Request (view stats)
    ZC->>+SM: Get Stats
    SM-->>-ZC: Stats Data
    ZC-->>-User: Response (Stats)

    User->>+ZC: Request (link Zettels)
    ZC->>+Z: Create Link
    Z-->>-ZC: Link Created
    ZC-->>-User: Response (Zettels linked)

    alt Create Zettel
        User->>+VC: Request Create Zettel Page
        VC-->>-User: Show Create Zettel Form
        User->>+VC: Submit Create Zettel
        VC->>+ZC: Forward Create Zettel
        deactivate VC
        ZC->>+Z: Create Zettel
        Z->>SM: Increment Creation Stats
        loop Update Stats
            SM->>SM: Update Stats
        end
        Z-->>-ZC: Zettel Created
        activate VC
        ZC-->>VC: Update View
        VC-->>-User: Display Zettel Created
    end

    alt View Stats
        User->>+VC: Request Stats Page
        VC-->>-User: Show Stats Request Form
        User->>+VC: Submit Stats Request
        VC->>+ZC: Fetch Stats
        deactivate VC
        ZC->>+SM: Get Stats
        SM-->>-ZC: Stats Data
        activate VC
        ZC-->>VC: Update View with Stats
        VC-->>-User: Display Stats
    end
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
