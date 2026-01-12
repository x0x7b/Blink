# Blink

![status](https://img.shields.io/badge/status-release--1.0-green)

**Blink** is a diff-based HTTP behavior analysis tool.

It compares baseline responses against injected payloads and **scores behavioral deviations** instead of guessing vulnerabilities.

Blink answers one question:
**“Which inputs change server behavior in a non-trivial way?”**

---

## Core idea

Traditional scanners try to label vulnerabilities.
Blink does not.

Blink:

* sends controlled payloads
* observes *behavioral differences*
* builds a **behavior profile**
* assigns a **score** to each payload based on how rare and strong its effects are

You decide what the behavior means.

---

## Features

* Baseline vs payload response comparison
* Behavioral diff types:

  * body hash changes
  * reflected input (raw / encoded)
  * timing anomalies (RTT)
  * headers and redirects
* Automatic **behavior profile** per endpoint
* **Payload scoring** based on diff rarity and weight
* Sorting by score and `--top N` filtering
* URL query parameter testing
* HTML form testing
* Clean, low-noise CLI output

---

## Example

```bash
$ blink -url-params --top 5 http://testphp.vulnweb.com/artists.php?artist=1
```

Output:

```
[ Blink 1.0 ]  

200 [ GET http://testphp.vulnweb.com/artists.php?artist=1 ] (477.83979ms)

[WARN] Showing results ONLY with diffs
 /artists.php  4.61  artist=%250a
   BODY_HASH    : b190247d32 -> b9cb738867 
   RTT          :  214ms →  800ms (x2.74)
```

Higher score → rarer + stronger behavioral change.

---

## Scoring model (high level)

Each diff type has:

* a base weight
* a frequency in the current scan

Score increases when:

* the diff is **rare**
* the diff is **high-impact** (e.g. RTT)

This suppresses noise and surfaces anomalies automatically.

---

## Usage

```bash
blink <target>

blink -url-params http://example.com/search?q=test
blink -forms http://example.com/login.php
blink -url-params -forms --top 10 http://example.com/
```

---

## Common flags

* `-forms`         test HTML forms
* `-url-params`    test URL query parameters
* `--top N`        show only top N scored results
* `-help`          show help and usage

---

## Installation

```bash
git clone https://github.com/x0x7b/Blink
cd Blink
go run main.go
```

---

## Project status


v1.0 Released

---

## License

MIT

