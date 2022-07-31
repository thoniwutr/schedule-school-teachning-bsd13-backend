# merchant-config-svc

Merchant Config Service

## Usage

### Running Locally

Set up your environment by creating `.env` file in the project:

``` bash
PROJECT_ID=123
DEBUG=true
PORT=8080
SERVER_TIMEOUT_READ=15s
SERVER_TIMEOUT_WRITE=15s
SERVER_TIMEOUT_IDLE=30s

# Generate below by running $(gcloud beta emulators datastore env-init)
DATASTORE_DATASET=beam-payment-development
DATASTORE_EMULATOR_HOST=localhost:8081
DATASTORE_EMULATOR_HOST_PATH=localhost:8081/datastore
DATASTORE_HOST=http://localhost:8081
DATASTORE_PROJECT_ID=beam-payment-development
```

Start up the emulator like so in one terminal:

`gcloud beta emulators datastore start --no-store-on-disk`

Then, in terminal 2, run:

`$(gcloud beta emulators datastore env-init)`

Also in the terminal 2, run:

`go run cmd/merchant-config-svc/main.go`

## Generating Model Code from OpenAPI

`oapi-codegen -generate types -o api/test.gen.go api/merchantconfig.yaml`

Currently only using the models generated from the openapi specifications and nothing else.

DATASTORE_EMULATOR_HOST localhost:8081