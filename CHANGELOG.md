# Changelog

## v0.6

### Added
- Payload scoring system
- Result sorting by score
- `--top N` flag for highest-signal results
- Structured diff output (machine-usable)
- Progress bar for query scans
- Wordlist CLI flag

### Changed
- Core returns structured diff data instead of printing
- Output layer fully separated from core logic
- Diff analysis based on typed data structures
- Link discovery via HTML parsing instead of regex
- Wordlist loaded once per scan (performance)

### Improved
- Reflection detection (raw vs encoded, no duplicates)
- URL parameter diff output formatting
- Error handling and output stability
- Progress bar correctness

### Removed
- Verbose diff mode
- Duplicate reflected diff entries
- Trash / built binaries from repo


## v0.5

### Added
- Diff-based response analysis (body, headers)
- URL query parameter testing
- HTML forms parsing and testing
- Multiple diff output modes
- Redirect chain analysis
- Server fingerprinting
- Network timings (DNS, TCP, TLS, TTFB)
- Basic error handling and reporting
- ASCII logo

### Changed
- Output logic split into separate functions
- Project structure refactored into multiple packages
- Manual redirect handling instead of default client behavior

### Fixed
- Redirect handling edge cases
- Body output formatting
- URL query diff styling
- Import cycles and architecture issues


