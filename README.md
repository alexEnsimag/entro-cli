### Entro-CLI

Entro-CLI is a CLI client for Entro.

### Installation

- Clone the repository.
- At the root of the repository, run `go build`.
- To get the list of operations: run `./entro-cli -h`.

### Usage

Before anything, set your AWS credentials as environment variables: 
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `AWS_SESSION_TOKEN`
- `AWS_REGION`

Then,
- To request a report, run: `./entro-cli create`, save the printed report ID.
- To get the status of a report, run: `./entro-cli status <reportID>`.
- To download the report, run :`./entro-cli download <reportID>`.
