# Zettelkasten under a Terminal

```text
/myapp
  /cmd
    /myapp
      main.go
  /pkg
    mypackage.go
  /api
    api_definitions.go
  /internal
    /mypackage
      mypackage_internal.go
  /scripts
  /web 
  /assets
  /test
  /vendor
  .gitignore
  README.md
  go.mod
  go.sum
```

**Here's what each directory does:**

- `cmd`: This is where the application's main() functions live. Every application
  in your project should have a subdirectory under cmd.
- `pkg`: Libraries and packages that are okay to be used by applications outside
  of your project. Other projects will import these libraries expecting them to
  work.
- `api`: Definitions of services, protocols, and declarations of the API types.
- `internal`: Private application and library code. This code cannot be imported
  from outside your project.
- `scripts`: Scripts to perform various build, install, analysis, etc operations.
- `web`: Web-related components: static web assets, SPAs, HTML templates, etc.
- `ui`: UI-related components.
- `assets`: Other assets needed for your project to work.
- `test`: Additional external test apps and test data.
- `vendor`: Application dependencies (manually vendored, if needed).

# TODOS

- [ ] Migrate to sqlite so that you can take advantage of sql to get specific
      data or connections & refactor the codebase (idiomatic golang project
      structure)

- [ ] Add better integration with neovim by creating a plugin on top of this
      binary

  - [ ] Add Image support c-v and c-p
  - [ ] Add `nvim-cmp` plugin with sqlite capabilities (maybe some crazy
        syntax is possible here)
  - [ ] Add all basic capabilities to keymaps & commands (query, delete, ...)

- [ ] Create a pwa for the thing
