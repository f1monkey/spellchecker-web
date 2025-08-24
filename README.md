# spellchecker-web

`spellchecker-web` is a Go-based ready-to-use spellchecker web service. It leverages [f1monkey/spellchecker](https://github.com/f1monkey/spellchecker) as its core for spell checking and correction, supports multiple dictionaries and has an autosave feature.

## Configuration

Before running the service, set the following environment variables:

|Variable        | 	Description | Example | Default value | Required |
|----------------|------------- |---------|---------------|----------|
|SPELLCHECKER_DIR| 	Directory to store dictionaries |	/tmp/spellchecker | none | yes |
|SPELLCHECKER_AUTOSAVE_INTERVAL| 	Auto-save interval (Go time.Duration) | 5m | none | no |
|SPELLCHECKER_WORD_SPLIT_REGEXP| Regular expression used to split phrases by words | ['\pL]+ | ['\pL]+| no |
|SPELLCHECKER_HTTP_ADDR| 	HTTP server address and port | localhost:8011 | localhost:8011 | no |
|SPELLCHECKER_LOG_LEVEL| 	Logging level |	error | info | no |

## Swagger Docs

The OpenAPI (Swagger) documentation is available at /docs.

## Usage Example

__Notes__: The alphabet is case sensitive, as are phrases, so you'll need to convert them to lowercase or add the letters A-Z to the alphabet to make the spellchecker work with capitals.

1) Create a dictionary `my-dictionary`:

```
POST /v1/dictionaries/my-dictionary
Content-Type: application/json

{
  "alphabet": "abcdefghijklmnopqrstuvwxyz",
  "maxErrors": 2
}
```

2) Add some words to `my-dictionary`
```
POST /v1/dictionaries/my-dictionary/add
Content-Type: application/json

{
    "phrases": [
        {
            "text": "weapon",
            "weight": 1
        },
        {
            "text": "the knight raised his weapon before charging into battle",
            "weight": 1
        }
    ]
}
```

3) Fix phrase
```
POST /v1/dictionaries/my-dictionary/fix
Content-Type: application/json

{
    "text": "the knight raised his waapon befor charging into battl"
}
```

Response:

```
{
    "fixes": [
        {
            "start": 22,
            "end": 28,
            "suggestions": [
                {
                    "text": "weapon",
                    "score": 0.8239592165010822
                }
            ],
            "error": "invalid_word"
        },
        {
            "start": 29,
            "end": 34,
            "suggestions": [
                {
                    "text": "before",
                    "score": 0.7797905781299383
                }
            ],
            "error": "invalid_word"
        },
        {
            "start": 49,
            "end": 54,
            "suggestions": [
                {
                    "text": "battle",
                    "score": 0.7797905781299383
                }
            ],
            "error": "invalid_word"
        }
    ],
    "correct": [
        {
            "start": 0,
            "end": 3
        },
        {
            "start": 4,
            "end": 10
        },
        {
            "start": 11,
            "end": 17
        },
        {
            "start": 18,
            "end": 21
        },
        {
            "start": 35,
            "end": 43
        },
        {
            "start": 44,
            "end": 48
        }
    ]
}
```