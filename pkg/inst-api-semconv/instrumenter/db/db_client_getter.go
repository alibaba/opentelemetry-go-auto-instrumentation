package db

type DbClientCommonAttrsGetter[REQUEST any] interface {
	GetSystem(REQUEST) string
	GetUser(REQUEST) string
	GetName(REQUEST) string
	GetConnectionString(REQUEST) string
}

type DbClientAttrsGetter[REQUEST any] interface {
	DbClientCommonAttrsGetter[REQUEST]
	GetStatement(REQUEST) string
	GetOperation(REQUEST) string
}

type SqlClientAttributesGetter[REQUEST any] interface {
	DbClientCommonAttrsGetter[REQUEST]
	GetRawStatement(REQUEST) string
}
