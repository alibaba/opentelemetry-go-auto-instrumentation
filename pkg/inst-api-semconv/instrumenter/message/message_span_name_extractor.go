package message

const temp_destination_name = "(temporary)"

type MessageSpanNameExtractor[REQUEST any, RESPONSE any] struct {
	getter        MessageAttrsGetter[REQUEST, RESPONSE]
	operationName MessageOperation
}

func (m *MessageSpanNameExtractor[REQUEST, RESPONSE]) Extract(request REQUEST) string {
	destinationName := ""
	if m.getter.IsTemporaryDestination(request) {
		destinationName = temp_destination_name
	} else {
		destinationName = m.getter.GetDestination(request)
	}
	if destinationName == "" {
		destinationName = "unknown"
	}
	if m.operationName != "" {
		destinationName = destinationName + " " + string(m.operationName)
	}
	return destinationName
}
