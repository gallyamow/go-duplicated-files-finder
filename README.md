## go-duplicated-files-finder

A lightweight utility written in Go that allows you to find duplicated files in a directory.

### Usage

```shell
./dupfinder --exclude-dir=node_modules,vendor,venv,cache,.gradle,ramil --min-size=1MB --workers=32 --format=plain /home/ramil/Projects
```