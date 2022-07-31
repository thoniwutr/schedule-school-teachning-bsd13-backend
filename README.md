# schedule school teaching services

Schedule School Teaching Services

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
PROJECT_ID=bsd-schedule-teaching;DEBUG=true;PORT=7001;SERVER_TIMEOUT_READ=600s;SERVER_TIMEOUT_WRITE=600s;SERVER_TIMEOUT_IDLE=30s;PUBLISHER_TOPIC_ID=merchant-create;GOOGLE_APPLICATION_CREDENTIALS=D:\bsd13\schedule-school-teaching-bsd13-backend\bsd-schedule-teaching-c983423ae892.json;KYM_BUCKET_NAME=beam-development-315606_kym_documents;RECIPIENT_SERVICE_URL=httpp;ORGANISATION_SERVICE_URL=test
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