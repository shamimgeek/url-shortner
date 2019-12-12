# An URL shortener written in Golang

I wanted to try and see if I could build an URL shortener in Go.

## Features
* Doesn’t create multiple short URLs when you try to shorten the same URL. In this case, the script will simply return the existing short URL for that long URL.
* Redirect to your main website when no shortid, or incorrect shortid, is entered, e.g. `http://localhost:8080/xdf` → `http://localhost:8080/`.
* Generates short URLs using only `[a-z0-9]` characters.

## Installation
1. Download the source code and install it using the `go install` command.
2. Use `database.sql` to create the `tiny_urls` table.
3. Create a config file in `/path/to/url-shortener/` named `config.(json|yaml|toml)`. Use `config-example.json` as a example.


## To-do
* Add tests
* Add json request and response. currently only plain txt is requested and serverd.
* Build Docker image for kuberntes. 
* Add deployment

