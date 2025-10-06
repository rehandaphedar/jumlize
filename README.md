# Introduction

A program to segment Qurʾān translation text.

# Installation

```sh
go install git.sr.ht/~rehandaphedar/lafzize/v2@latest
```

# Usage

Obtain a translation from the [Quranic Universal Library (QUL)](https://qul.tarteel.ai/resources/translation) (`simple.json`), or from any other source with the same schema.

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

The command will print a list of verse keys whose segmented text doesn't yield the original text on joining with a `" "`.

## Clear

Run:

```sh
jumlize clear
```

The command will remove the `segments` field from all verses in the file.
