# Use extensions to print http request and response header

The document demonstrates how to use extensions to print http request and response header.

## Step1: Compile the target binary with otel
Use `otel` to build the binary with `config.json`:
```
cd example/extension/nethttp
../../../otel set -rule=config.json
../../../otel go build .
```
Users can get the `otel` according to [documentation](../../../README.md)

## Step2: Run the binary compiled by otel
```shell
./demo
```
And the result will be:
```shell
request header is  {"otel":["true"]}
response header is  {"Content-Type":["application/x-gzip"],"Date":["Wed, 06 Nov 2024 11:35:37 GMT"],"Server":["bfe"]}
```
It means that the `nethttp` extension can print http headers correctly.
