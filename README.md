# Blink
![status](https://img.shields.io/badge/status-pre--1.0-lightgrey)

Blink is a diff-based HTTP analysis tool focused on detecting behavioral changes
between baseline and injected requests.

It is designed to highlight *what changes*, not to guess vulnerabilities.

## Features

- Baseline vs payload response diffing
- HTML form testing
- URL query parameter testing
- Reflected input detection (raw)
- Timing-based behavior changes (RTT)
- Header and redirect diffing
- Clean, low-noise CLI output

## Example

```bash
$ blink http://example.com/search?q=test

/search.php q=abcd
  body_hash : changed
  reflect   : raw input
```

## Usage

```bash
blink <target>

blink -forms http://example.com/login.php
blink -url-params http://example.com/search?q=test
```

## Common flags

- `-forms`        test HTML forms
- `-url-params`   test URL query parameters
- `-diff-verbose` show detailed diffs
- `-help`         show help and usage

## Installation

```bash
git clone https://github.com/x0x7b/Blink
cd Blink
go run main.go
```

## Project status

Pre-1.0. Active development.
The CLI interface is not yet frozen and may change.

## License

MIT

