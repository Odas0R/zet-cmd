# `zet` - Your Terminal Zettelkasten

`zet` incorporates the simplicity of a
CLI with the organization capabilities of the zettelkasten method, enabling
users to manage knowledge with efficiency.

## Project Objectives

The aim of `zet` is to provide a streamlined, terminal-based approach to the
zettelkasten method, making it accessible and practical for daily use. Whether
you're managing notes, ideas, or extensive research.

We also want `zet` to be independent, you can create a custom client to
integrate or extend the existing functionality.

## Features

- ✅ **Create, Open, and Remove Zettels**: Easily manage your notes from the command line.
- ✅ **Search**: Utilize SQLite's FTS5 extension for powerful full-text search capabilities.
- ✅ **History and Backlog**: Keep track of your most recent and overall zettel landscape.
- ✅ **Sync and Save**: Keep your filesystem and database in harmony, with automatic fixes on the go.

**Note:** It's advised to use marksman LSP to have a better experience with the
zettelkasten method.

## Usage

```text
NAME:
   zet - A zettelkasten under a terminal approach

USAGE:
   A simple way to manage your zettelkasten using neovim (telescope) and fzf

VERSION:
   0.1

AUTHOR:
   odas0r <guilherme@muxit.co>

COMMANDS:
   new          Create a new zettel
   open         Opens the zettel by the given path
   search       Search for zettels using sqlite3 fs5 extension
   remove, rm   Removes the given zettel from the database and from the filesystem
   history      Retrieves the last 50 opened zettel
   backlog      Retrieves all the fleet of zettels
   brokenlinks  Retrieves all the brokenlinks of a zettel
   permanent    Sets the given zettel as type permanent
   fleet        Sets the given zettel as type fleet
   last         Retrieves the last opened zettel
   save         Inserts or updates the given zettel to the database, and some repairs
   sync         Sync the filesystem with the database and does some fixing on the side
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open
issues.

## License

`zet` is open-source and proudly built on the shoulders of giants. It's available
under the MIT license.
