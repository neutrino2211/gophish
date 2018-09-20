# Gophish

Gophish is a tool for tracking users over the internet who don't use vpns

## Usage

```bash
$ go run gophish.go --help
Gophish v0.0.1

[ gophish ] gophish <url> [options?]
                url -> url to redirect target after getting information (url should not start with http:// )
                options
                        --admin -> Admin link provided by first launch to continue (url not needed)
                        --timeout -> optional interval for checking if the linked was clicked
```