package utils

import "go.opentelemetry.io/otel/attribute"

const DB_CLIENT_KEY = attribute.Key("opentelemetry-traces-span-key-db-client")
const RPC_SERVER_KEY = attribute.Key("opentelemetry-traces-span-key-rpc-server")
const RPC_CLIENT_KEY = attribute.Key("opentelemetry-traces-span-key-rpc-client")
const PRODUCER_KEY = attribute.Key("opentelemetry-traces-span-key-producer")
const CONSUMER_RECEIVE_KEY = attribute.Key("opentelemetry-traces-span-key-consumer-receive")
const CONSUMER_PROCESS_KEY = attribute.Key("opentelemetry-traces-span-key-consumer-process")
