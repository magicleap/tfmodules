# Changelog

## [0.3.0] - 2021-12-22
### ADDED
- Add envar `OVERWRITE` to prevent overwriting an existing version

## [0.2.0] - 2021-12-21
### ADDED
- Add source of module to upload
- Fetch source in getLatest for renovate release notes management

## [0.1.0] - 2021-11-19
### ADDED
- Initiate project
- Add API for TF modules regisrty compatibility ie
  - discovery
  - fetch versions per module
  - download module
  - upload module
- Add API compatible with renovateBot scanning ie
  - fetch latest version of a module
- Add GCS backend and fake backend for test
- Build with `ko`

