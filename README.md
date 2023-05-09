# Zettelkasten under a Terminal

This is just a public repo, you may take some ideas from it since it is not
extensible.

The idea is to be able to:

1. Query a zettel
2. Create links between zettels
3. Intuitive backlog and history
4. Auto generate tags for zettels
5. Export data for static site generation _SSG_

The possibilities here are many. I'll add more specific features in the future
and log what works and doesn't work for me.

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
