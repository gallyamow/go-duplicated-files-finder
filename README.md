## go-duplicated-files-finder

A simple command-line tool written in Go that scans a directory tree, identifies duplicate files by content hash, and
prints results in a machine-friendly format.

### Usage

```shell
dupfinder [flags] <path>
```

| Flag            | Default            | Description                                          |
|-----------------|--------------------|------------------------------------------------------|
| `--min-size`    | `1B`               | Minimum file size to include (e.g., `10MB`, `500KB`) |
| `--exclude-ext` | `""`               | Comma-separated file extensions to ignore            |
| `--exclude-dir` | `""`               | Comma-separated directory names to skip              |
| `--algo`        | `md5`              | Hash algorithm (`md5`, `sha1`, `sha256`)             |
| `--workers`     | 2 x number of CPUs | Number of concurrent workers                         |
| `--format`      | `paths`            | Output format: `plain`, `paths`, `json`              |
| `--delete`      | `false`            | Delete the duplicate files (use with caution)        |

### Example

```shell
./dupfinder --exclude-dir=node_modules,vendor,venv,cache,.gradle --min-size=1MB --workers=32 --format=plain /home/ramil/Projects
```