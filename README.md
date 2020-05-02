# seqren - sequence rename tool

## Installation

Download from [Release](https://github.com/kazuya0202/seqren/releases).

## Usage

```sh
$ seqren -h
Rename filename in sequence.

Usage:
  seqren [name] [flags]

Flags:
  -a, --all-show      display all lines of renaming target
  -f, --force         execute command without confirmation
  -h, --help          help for seqren
  -n, --num int       display lines of renaming target (default 10)
  -p, --path string   path of target directory
  -s, --seq int       N-digit 0 filling (default 3)
```

### Default Usage

```sh
$ seqren [name] [flags]
```