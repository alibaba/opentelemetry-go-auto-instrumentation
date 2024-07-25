//go:build ignore

package mongo

type mongoRequest struct {
	CommandName  string
	DatabaseName string
	ConnectionID string
}
