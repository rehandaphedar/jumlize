# Introduction

A program to segment Qurʾān translation text.

# Installation

```sh
go install git.sr.ht/~rehandaphedar/jumlize/v3@latest
```

# Usage

From the [Quranic Universal Library (QUL)](https://qul.tarteel.ai/resources/translation) (or from any other source with the same schema) obtain the following:
- [Ayah Metadata](https://qul.tarteel.ai/resources/quran-metadata) (`quran-metadata-ayah.json`).
- A [translation](https://qul.tarteel.ai/resources/translation) (`-simple.json`).

The program comes with 3 subcommands:
- `segment`
- `check`
- `clear`

For each subcommand, you can get documentation for flags by running:

```
jumlize [subcommand] -h
```

## Segment

Run:

```sh
jumlize segment -api_key "[GEMINI_API_KEY]"
```

The command will add a `segments` field to all verses in the file.

## Check

Run:

```sh
jumlize check
```

The command will sanity check the segments and print a list of verse keys with invalid segments.

## Clear

Run:

```sh
jumlize clear
```

The command will remove the `segments` field from all verses in the file.

# Results

Segmented translations can be found under `results/`.
