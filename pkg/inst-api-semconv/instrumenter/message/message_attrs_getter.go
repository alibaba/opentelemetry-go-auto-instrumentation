package message

type MessageAttrsGetter[REQUEST any, RESPONSE any] interface {
	GetSystem(request REQUEST) string
	GetDestination(request REQUEST) string
	GetDestinationTemplate(request REQUEST) string
	IsTemporaryDestination(request REQUEST) bool
	isAnonymousDestination(request REQUEST) bool
	GetConversationId(request REQUEST) string
	GetMessageBodySize(request REQUEST) int64
	GetMessageEnvelopSize(request REQUEST) int64
	GetMessageId(request REQUEST, response RESPONSE) string
	GetClientId(request REQUEST) string
	GetBatchMessageCount(request REQUEST, response RESPONSE) int64
	GetMessageHeader(request REQUEST, name string) []string
}
