package dubbo

type dubboRequest struct {
	methodName    string
	serviceKey    string
	serverAddress string
	attachments   map[string]any
}

type dubboResponse struct {
	hasError bool
	errorMsg string
}
