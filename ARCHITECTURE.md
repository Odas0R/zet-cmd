# Architecture

Documenting the `zet-cmd` architecture.

## Models

`Zettel` is the central entity , representing individual notes or "zettels".
The model includes essential fields for identifying and managing zettels, as
well as auxiliary fields to enhance functionality, like Lines for content
manipulation and Links for inter-zettel connections.

```go
type Zettel struct {
	ID        string `db:"id"`
	Title     string `db:"title"`
	Content   string `db:"content"`
	Path      string `db:"path"`
	Type      string `db:"type"`
	CreatedAt Time   `db:"created_at"`
	UpdatedAt Time   `db:"updated_at"`

	// Auxiliary fields (not stored in the database)
	Lines []string
	Links []*Zettel
}
```

## Controllers

- `ZettelController`: Endpoint for Zettel management operations, interfacing
  with `ZettelService` to process user requests like creation, update, and
  deletion of zettels.

- `LinkController`: Facilitates operations related to linking zettels,
  including creating and removing links, and querying for backlinks or all
  links related to a zettel.

- `SearchController`: Provides endpoints for searching the Zettelkasten,
  interfacing with `SearchService` to execute and return search results.

- `HistoryController`: Offers access to a user's history within the
  Zettelkasten, such as recently edited or viewed zettels, leveraging the
  `HistoryService`.

- `ViewController`: Handles server-side rendering (SSR) of pages, presenting
  the user interface for interacting with the Zettelkasten. This includes
  rendering views for search results, zettel content, and navigation.

## Views

Views are the SSR templates or components that render the Zettelkasten's user
interface:

1. List view for search results or zettel collections (e.g., backlog, history).
2. Detail view for individual zettels, displaying content, metadata, and
   links/backlinks.
3. Edit view for creating or updating zettels.
4. Navigation and search interfaces, providing easy access to various parts of
   the Zettelkasten.

## Flux Sequence Diagram E.g: Creating, Linking & Stats

```mermaid
sequenceDiagram
    participant User
    participant ZC as ZettelController
    participant ZS as ZettelService
    participant ZKS as ZettelkastenService
    participant LS as LinkService
    participant Z as Zettel
    participant DB as Database
    participant FS as Filesystem

    User->>+ZC: Request (create Zettel)
    ZC->>+ZS: Create Zettel
    ZS->>+Z: Initialize
    Z->>-ZS: Zettel Initialized
    ZS->>DB: Save Zettel
    ZS->>FS: Save File
    ZS-->>-ZC: Zettel Created
    ZC-->>-User: Response (Zettel created)

    User->>+ZC: Request (view stats)
    ZC->>+ZKS: Get Stats
    ZKS->>+DB: Query Stats
    DB-->>-ZKS: Stats Result
    ZKS-->>-ZC: Stats Data
    ZC-->>-User: Response (Stats)

    User->>+ZC: Request (link Zettels)
    ZC->>+LS: Create Link
    LS->>DB: Update Links
    LS-->>-ZC: Link Created
    ZC-->>-User: Response (Zettels linked)
```

## Flux Diagram MVC E.g: Creating, Linking & Stats

```mermaid
flowchart TD
    User("User") --> View{"View\n(User Interface)"}
    View --> Controller{"Controller\n(ZettelController, etc.)"}
    Controller -->|Create/Update/Delete| Service(("Services\n(ZettelService, LinkService, etc.)"))
    Service --> Model(("Model\n(Zettel)"))
    Service -->|Query| DB[(Database)]
    Service -->|Read/Write| FS[(Filesystem)]
    Model --> DB
    Model --> FS
    DB --> Model
    FS --> Model
    Model -->|Data Response| Service
    Service -->|Logic Processed| Controller
    Controller -->|Render| View
    View -->|Display| User

    Controller -->|Get Stats| ZKS(("ZettelkastenService"))
    ZKS -->|Query| DB
    DB -->|Stats Result| ZKS
    ZKS -->|Processed Stats| Controller
    Controller -->|Render Stats| View

    Controller -.->|Sync| ZS(("ZettelService"))
    ZS -.->|Filesystem Check| FS
    FS -.->|Update| DB
    DB -.->|Sync Complete| ZS
    ZS -.->|Notify| Controller
    Controller -.->|Update View| View
```
