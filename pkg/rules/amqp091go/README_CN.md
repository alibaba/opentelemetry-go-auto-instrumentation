## **consume module**

监听github.com/rabbitmq/amqp091-go下consumers类的send方法，因为此处message未放入chan中

## **pulish module**

监听github.com/rabbitmq/amqp091-go下Channel类的send方法。CorrelationId作为传递traceparent用于关联pulish和consume。当使用者自己填写了CorrelationId时，CorrelationId不会被覆盖，此时pulish和consume不会关联