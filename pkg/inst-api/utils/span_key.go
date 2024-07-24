package utils

import "go.opentelemetry.io/otel/attribute"

const ARMS_DB_CLIENT_START_TIME = "arms-db-client-start-time"

const DB_CLIENT_KEY = attribute.Key("opentelemetry-traces-span-key-db-client")
const RPC_SERVER_KEY = attribute.Key("opentelemetry-traces-span-key-rpc-server")
const RPC_CLIENT_KEY = attribute.Key("opentelemetry-traces-span-key-rpc-client")
const PRODUCER_KEY = attribute.Key("opentelemetry-traces-span-key-producer")
const CONSUMER_RECEIVE_KEY = attribute.Key("opentelemetry-traces-span-key-consumer-receive")
const CONSUMER_PROCESS_KEY = attribute.Key("opentelemetry-traces-span-key-consumer-process")

const ARMS_DB_CLIENT_METRICS_STATE = "arms-db-client-metrics-state"

const ARMS_DB_CLIENT_METRICS_ATTR_KEY = attribute.Key(ARMS_DB_CLIENT_METRICS_STATE)

const STATUS = "status"
