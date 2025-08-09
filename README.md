# Introduction

A program to split Qurʾān translation text into sentences.

This program was mainly created for use by [tilavid](https://sr.ht/~rehandaphedar/tilavid).
Therefore, the `data.json` used, and subsequent configuration etc. was with the English translation [Saheeh International from Quranic Universal Library](https://qul.tarteel.ai/resources/translation/193) (`simple.json`) in mind.
However, this program can be adapted to other languages and translations as well.

Readily segmented and manually fixed translations can be found under `results/`.

# Dependencies

- `python`
- `spacy`

# Installation

```sh
pip install git+https://git.sr.ht/~rehandaphedar/jumlize
```

You can also use `pipx`.

# Usage

Run:

```sh
jumlize split [language_code] [data_file] [sentences_file]
```

Where:

- `language_code`: Two letter language code

- `data_file`: JSON file of the following format:

```json
{
	"[verse_key]": {
		"t": "[translation]"
	},
	...
}
```

- `sentences_file`: JSON file to output segmented text to in the following format:

```json
{
	"[verse_key]": [
		"[sentence_1]",
		"[sentence_2]",
		...
	],
	...
}
```

Example command:

```sh
jumlize split en data.json sentences.json
```

The generated `sentences.json` will likely have some errors. For English, these mostly tend to be due to spacy splitting verses on `[` or `"` and then failing to preserve the exact text.
Run the following command to print a list of verse keys whose segmented text doesn't yield the original text on joining with a `" "`:

```sh
jumlize check [data_file] [sentences_file]
```

Example command:

```sh
jumlize check data.json sentences.json
```

After this, fix these errors manually or using another tool. Rerun the command and repeat until no verse keys are printed.

Do note that as stated above, the `check` command only checks if the segments of a particular verse key yield the original text on joinig. It does not check for anomalous sentence boundaries being used or any other error that still causes the segments to yield the original text.
