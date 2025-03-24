## **consume module**

Listen to the "send" method of the consumers class in github.com/rabbitmq/amqp091-go, as the message is not put into the chan .

## **pulish module**

Listen to the "send" method of the Channel class in github.com/rabbitmq/amqp091-go. The CorrelationId is used to pass the traceparent for associating publish and consume. If the user manually sets the CorrelationId, it will not be overwritten, and in this case, publish and consume will not be associated.