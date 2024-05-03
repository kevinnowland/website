# Website

This repo is the basis of my personal website.

## Run locally

To run directly with go:

```bash
docker run website
```

To build then run with go:

```bash
go build
./website
```

To build and run with docker:

```bash
docker build -t website .
docker run --rm -p 8080:8080 website
```

and then go to `localhost:8080` in a browser.
