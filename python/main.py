from typing import cast
import sys
import json
import spacy


def split():
    language = "en"
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

    with open("data.json", "r") as dataFile:
        data = cast(dict[str, dict[str, str]], json.load(dataFile))
        verse_keys = list(data.keys())
        translations = [data[key]["t"] for key in verse_keys]

        sentences = {
            verse_key: [sent.text for sent in doc.sents]
            for verse_key, doc in zip(verse_keys, nlp.pipe(translations))
        }

        with open("sentences.json", "w") as sentencesFile:
            json.dump(sentences, sentencesFile)


def check():
    with open("data.json", "r") as dataFile:
        data = cast(dict[str, dict[str, str]], json.load(dataFile))
        verses = [(key, value["t"]) for key, value in data.items()]

        with open("sentences.json", "r") as sentencesFile:
            sentences = cast(dict[str, list[str]], json.load(sentencesFile))

            for verse_key, translation in verses:
                if translation != " ".join(sentences[verse_key]):
                    print(verse_key)


match sys.argv[1]:
    case "split":
        split()
    case "check":
        check()
    case default:
        exit()
