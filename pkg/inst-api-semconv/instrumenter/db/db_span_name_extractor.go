package db

type DBSpanNameExtractor[REQUEST any] struct {
	getter DbClientAttrsGetter[REQUEST]
}

func (d *DBSpanNameExtractor[REQUEST]) Extract(request REQUEST) string {
	dbName := d.getter.GetName(request)
	operation := d.getter.GetOperation(request)
	if operation == "" {
		if dbName == "" {
			return "DB Query"
		} else {
			return dbName
		}
	}
	return operation
}
