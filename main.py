from typing import cast
import sys
import json
import spacy


def split():
    language = sys.argv[2]
    data_path = sys.argv[3]
    sentences_path = sys.argv[4]

    nlp = spacy.load(
        f"{language}_core_web_sm",
        exclude=[
            "tok0vec",
            "tagger",
            "parser",
            "senter",
            "attribute_ruler",
            "lemmatizer",
            "ner",
        ],
    )
    _ = nlp.add_pipe("sentencizer")

    with open(data_path, "r") as data_file:
        data = cast(dict[str, dict[str, str]], json.load(data_file))
        verse_keys = list(data.keys())
        translations = [data[key]["t"] for key in verse_keys]

        sentences = {
            verse_key: [sent.text for sent in doc.sents]
            for verse_key, doc in zip(verse_keys, nlp.pipe(translations))
        }

        with open(sentences_path, "w") as sentences_file:
            json.dump(sentences, sentences_file)


def check():
    data_path = sys.argv[2]
    sentences_path = sys.argv[3]

    with open(data_path, "r") as data_file:
        data = cast(dict[str, dict[str, str]], json.load(data_file))
        verses = [(key, value["t"]) for key, value in data.items()]

        with open(sentences_path, "r") as sentences_file:
            sentences = cast(dict[str, list[str]], json.load(sentences_file))

            for verse_key, translation in verses:
                if translation != " ".join(sentences[verse_key]):
                    print(verse_key)


def main():
    match sys.argv[1]:
        case "split":
            split()
        case "check":
            check()
        case _:
            exit()
