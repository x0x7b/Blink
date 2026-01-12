# Changelog

## v1.0

### Added
- Stable release: all core scanning, diffing, and output features finalized.
- CLI flags and behavior normalized.
- Full 1.0 documentation updated.

### Changed
- Minor internal fixes (wordlist path, safe iterations, normalized CLI flags).

### Fixed
- Panic calls removed.
- File handling improved (defer closes, safe hash computation).
- Typo fixes in function names.


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


