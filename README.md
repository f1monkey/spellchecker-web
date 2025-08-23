# spellchecker-web

`spellchecker-web` is a Go-based ready-to-use spellchecker web service. It leverages [f1monkey/spellchecker](https://github.com/f1monkey/spellchecker) as its core for spell checking and correction, supports multiple dictionaries and has an autosave feature.

## Configuration

Before running the service, set the following environment variables:

|Variable        | 	Description | Example |
|----------------|------------- |---------|
|SPELLCHECKER_DIR| 	Directory to store dictionaries |	/tmp/spellchecker |
|SPELLCHECKER_AUTOSAVE_INTERVAL| 	Auto-save interval (Go time.Duration) | 5m |
|HTTP_ADDR| 	HTTP server address and port | localhost:8011 |
|LOG_LEVEL| 	Logging level |	info, debug |

## Swagger Docs

The OpenAPI (Swagger) documentation is available at /docs.

## Usage Example

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
            "text": "The knight raised his weapon before charging into battle.",
            "weight": 1
        }
    ]
}
```
