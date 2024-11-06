# Use extensions to print http request and response header

The document demonstrates how to use extensions to print http request and response header.

## Step1: Replace path in config.json
Replace the `Path` in `config.json` with the actual absolute path of `rules` directory, for example:
``` json
[{
"ImportPath":"net/http",
"Function":"RoundTrip",
"OnEnter":"httpClientEnterHook",
"ReceiverType": "*Transport",
"OnExit": "httpClientExitHook",
"Path": "/Users/zhanghaibin/Desktop/opentelemetry-go-auto-instrumentation/example/extension/nethttp/rules"
}]
```

## Step2: Compile the target binary with otelbuild
Use `otelbuild` to build the binary with `config.json`:
```
cd example/extension/netHttp
../../../otelbuild -rule=config.json -- demo/net_http.go
```
Users can get the `otelbuild` according to [documentation](../../../README.md)

## Step3: Run the binary compiled by otelbuild
```shell
./net_http
```
And the result will be:
```shell
request header is  {"Otelbuild":["true"]}
response header is  {"Content-Type":["application/x-gzip"],"Date":["Wed, 06 Nov 2024 11:35:37 GMT"],"Server":["bfe"]}
```
It means that the `nethttp` extension can print http headers correctly.
