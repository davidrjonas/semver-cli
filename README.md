semver-cli
==========

[![Build Status](https://travis-ci.org/davidrjonas/semver-cli.svg?branch=master)](https://travis-ci.org/davidrjonas/semver-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/davidrjonas/semver-cli)](https://goreportcard.com/report/github.com/davidrjonas/semver-cli)

semver-cli is a simple command line tool to compare and manipulate version strings.

It is basically a cli wrapper around the excellent [Masterminds semver library](https://github.com/Masterminds/semver), with lots of help from the also excellent [Kingpin](https://github.com/alecthomas/kingpin)

```
usage: semver [<flags>] <command> [<args> ...]

Command-line semver tools. On error, print to stderr and exit -1.

Flags:
      --help     Show context-sensitive help (also try --help-long and --help-man).
  -v, --verbose  Verbose mode.

Commands:
  help [<command>...]
    Show help.

  satisfies <VERSION> <CONSTRAINTS>
    Test if a version satisfies a constraint. Exit 0 if satisfies, 1 if not. If verbose, print an explanation to stdout.

  greater <A> <B>
    Compare two versions. Exit 0 if the first is greater, 1 if not. If verbose, print greater to stdout.

  lesser <A> <B>
    Compare two versions. Exit 0 if the first is lesser, 1 if not. If verbose, print lesser to stdout.

  equal <A> <B>
    Compare two versions. Exit 0 if they are equal, 1 if not.

  inc <COMPONENT> <VERSION>
    Increment major, minor, or patch component.

  set <COMPONENT> <VERSION> <VALUE>
    Set prerelease or metadata component.
```

Example
-------

Deploy only when the CI tag is within constraints and is greater than what is currently released. For example, with constraints 1.* and a released version of 1.4, a version of 1.5 would be released but a version of 2.0 or 1.2 would not.

```bash
#!/bin/bash

CONSTRAINTS=$(jq .labels.constraints manifest.json)
RELEASED=$(jq .image manifest.json | cut -d: -f 2)

semver satisfies "$CI_TAG" "$CONSTRAINTS" || exit 1
semver greater   "$CI_TAG" "$RELEASED"    || exit 1

deploy
```
