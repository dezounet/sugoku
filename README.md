# Sugokud

[![CircleCI](https://circleci.com/gh/dezounet/sugokud.svg?style=shield)](https://app.circleci.com/pipelines/github/dezounet/sugokud)

Go daemon serving a multiplayer sudoku game.

## Build from sources

```bash
# Linux users
make linux

# OSX users
make darwin
```

Binary files are generated in `cmd/sugokud` directory.

### Docker containder

You can package a linux docker image with this:

```bash
# This create docker image sugokud:latest
make docker
```

Once instanciated, this container launch sugokud on TCP port 8080

## Contribute

### Commit changes

Before any commit, install pre-commit hooks:

```bash
pre-commit install
```

### Understand project layout

Inspired from <https://github.com/golang-standards/project-layout>
