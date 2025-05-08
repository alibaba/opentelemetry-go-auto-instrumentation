module github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg

go 1.23.0

replace github.com/alibaba/opentelemetry-go-auto-instrumentation => ../../opentelemetry-go-auto-instrumentation

replace github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg => ../../opentelemetry-go-auto-instrumentation/pkg

require (
	github.com/bytedance/gopkg v0.1.2
	github.com/cespare/xxhash/v2 v2.3.0
	github.com/cloudwego/hertz v0.9.6
	github.com/cloudwego/kitex v0.5.1
	github.com/elastic/elastic-transport-go/v8 v8.6.1
	github.com/elastic/go-elasticsearch/v8 v8.4.0
	github.com/emicklei/go-restful/v3 v3.7.0
	github.com/gin-gonic/gin v1.10.0
	github.com/go-kratos/kratos/v2 v2.6.3
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-sql-driver/mysql v1.9.1
	github.com/gofiber/fiber/v2 v2.43.0
	github.com/gomodule/redigo v1.9.0
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/golang-lru/v2 v2.0.7
	github.com/kataras/iris/v12 v12.2.0
	github.com/labstack/echo/v4 v4.0.0
	github.com/mark3labs/mcp-go v0.20.0
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/nacos-group/nacos-sdk-go/v2 v2.2.9 // not minimal
	github.com/prometheus/client_golang v1.20.5
	github.com/rabbitmq/amqp091-go v1.10.0
	github.com/redis/go-redis/v9 v9.0.5
	github.com/rs/zerolog v1.31.0
	github.com/sirupsen/logrus v1.9.3 // not minimal, but consistent with older alibaba-otel
	github.com/stretchr/testify v1.10.0
	github.com/tmc/langchaingo v0.1.13
	github.com/valyala/fasthttp v1.45.0
	github.com/xwb1989/sqlparser v0.0.0-20180606152119-120387863bf2
	go-micro.dev/v5 v5.0.0
	go.mongodb.org/mongo-driver v1.14.0 // FIXME: not minimal
	go.opentelemetry.io/contrib/instrumentation/runtime v0.60.0
	go.opentelemetry.io/otel v1.35.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.35.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.35.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.35.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.35.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.35.0
	go.opentelemetry.io/otel/exporters/prometheus v0.57.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.35.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.35.0
	go.opentelemetry.io/otel/exporters/zipkin v1.35.0
	go.opentelemetry.io/otel/metric v1.35.0
	go.opentelemetry.io/otel/sdk v1.35.0
	go.opentelemetry.io/otel/sdk/metric v1.35.0
	go.opentelemetry.io/otel/trace v1.35.0
	go.uber.org/zap v1.24.0 // FIXME: not minimal
	google.golang.org/grpc v1.71.0 // FIXME: not minimal
	gorm.io/driver/mysql v1.1.3
	gorm.io/gorm v1.22.0
	trpc.group/trpc-go/trpc-go v1.0.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.2.0 // indirect
	github.com/Masterminds/sprig/v3 v3.2.3 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20210428141323-04723f9f07d7 // indirect
	github.com/Shopify/goreferrer v0.0.0-20220729165902-8cddb4f5de06 // indirect
	github.com/acomagu/bufpipe v1.0.3 // indirect
	github.com/alibabacloud-go/alibabacloud-gateway-pop v0.0.6 // indirect
	github.com/alibabacloud-go/alibabacloud-gateway-spi v0.0.5 // indirect
	github.com/alibabacloud-go/darabonba-array v0.1.0 // indirect
	github.com/alibabacloud-go/darabonba-encode-util v0.0.2 // indirect
	github.com/alibabacloud-go/darabonba-map v0.0.2 // indirect
	github.com/alibabacloud-go/darabonba-openapi/v2 v2.0.10 // indirect
	github.com/alibabacloud-go/darabonba-signature-util v0.0.7 // indirect
	github.com/alibabacloud-go/darabonba-string v1.0.2 // indirect
	github.com/alibabacloud-go/debug v1.0.1 // indirect
	github.com/alibabacloud-go/endpoint-util v1.1.0 // indirect
	github.com/alibabacloud-go/kms-20160120/v3 v3.2.3 // indirect
	github.com/alibabacloud-go/openapi-util v0.1.0 // indirect
	github.com/alibabacloud-go/tea v1.2.2 // indirect
	github.com/alibabacloud-go/tea-utils v1.4.4 // indirect
	github.com/alibabacloud-go/tea-utils/v2 v2.0.7 // indirect
	github.com/alibabacloud-go/tea-xml v1.1.3 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.1800 // indirect
	github.com/aliyun/alibabacloud-dkms-gcs-go-sdk v0.5.1 // indirect
	github.com/aliyun/alibabacloud-dkms-transfer-go-sdk v0.1.8 // indirect
	github.com/aliyun/aliyun-secretsmanager-client-go v1.1.5 // indirect
	github.com/aliyun/credentials-go v1.4.3 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/apache/thrift v0.13.0 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/bytedance/sonic v1.13.2 // indirect
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/choleraehyq/pid v0.0.16 // indirect
	github.com/clbanning/mxj/v2 v2.5.5 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/cloudwego/fastpb v0.0.5 // indirect
	github.com/cloudwego/frugal v0.2.5 // indirect
	github.com/cloudwego/gopkg v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/cloudwego/netpoll v0.7.0 // indirect
	github.com/cloudwego/thriftgo v0.4.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dlclark/regexp2 v1.11.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-git/go-billy/v5 v5.3.1 // indirect
	github.com/go-git/go-git/v5 v5.4.2 // indirect
	github.com/go-kratos/aegis v0.2.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/form/v4 v4.2.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.20.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/flatbuffers v23.5.26+incompatible // indirect
	github.com/google/pprof v0.0.0-20240727154555-813a5fbdbec8 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/goph/emperror v0.17.2 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/huandu/xstrings v1.3.3 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/iris-contrib/schema v0.0.6 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jhump/protoreflect v1.8.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.0.0-20180206201540-c2b33e8439af // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kataras/golog v0.1.11 // indirect
	github.com/kataras/pio v0.0.13 // indirect
	github.com/kevinburke/ssh_config v0.0.0-20201106050909-4977a11b4351 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.9 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/microcosm-cc/bluemonday v1.0.26 // indirect
	github.com/miekg/dns v1.1.43 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nikolalohinski/gonja v1.5.3 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/nyaruka/phonenumbers v1.0.55 // indirect
	github.com/openzipkin/zipkin-go v0.4.3 // indirect
	github.com/orcaman/concurrent-map v0.0.0-20210501183033-44dafcb38ecc // indirect
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c // indirect
	github.com/panjf2000/ants/v2 v2.4.6 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pkoukk/tiktoken-go v0.1.6 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/savsgio/dictpool v0.0.0-20221023140959-7bf2e61cea94 // indirect
	github.com/savsgio/gotils v0.0.0-20230208104028-c358bd845dee // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/tdewolff/minify/v2 v2.20.19 // indirect
	github.com/tdewolff/parse/v2 v2.7.12 // indirect
	github.com/tidwall/gjson v1.17.3 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tinylib/msgp v1.1.8 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/urfave/cli/v2 v2.25.7 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/xanzy/ssh-agent v0.3.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	github.com/yargevad/filepathx v1.0.0 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/proto/otlp v1.5.0 // indirect
	go.starlark.net v0.0.0-20230302034142-4b1e35fe2254 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/arch v0.14.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/exp v0.0.0-20240404231335-c0f41cb1a7a0 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/time v0.8.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
	google.golang.org/genproto v0.0.0-20240528184218-531527333157 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	trpc.group/trpc-go/tnet v1.0.1 // indirect
	trpc.group/trpc/trpc-protocol/pb/go/trpc v1.0.0 // indirect
)
