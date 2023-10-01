### Entro-CLI

Entro-CLI is a CLI client for Entro.

### Installation

- Clone the repository.
- At the root of the repository, run `go build`.
- To get the list of operations: `./entro-cli -h`.

### Usage

- To request a report, run: `./entro-cli create`, save the printed report ID.
- To get the status of a report, run: `./entro-cli status <reportID>`.
- To download the report, run :`./entro-cli download <reportID>`.
