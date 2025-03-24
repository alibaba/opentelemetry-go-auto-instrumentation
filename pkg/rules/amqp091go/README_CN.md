## **consume module**

监听github.com/rabbitmq/amqp091-go下consumers类的send方法，因为此处message未放入chan中。此模块只是监测消息接收，对于收到后消息的具体处理，可以根据需求在自己项目中加入自定义监测。

## **pulish module**

监听github.com/rabbitmq/amqp091-go下Channel类的send方法。CorrelationId作为传递traceparent用于关联pulish和consume。当使用者自己填写了CorrelationId时，CorrelationId不会被覆盖，此时pulish和consume不会关联