package db

type DBSpanNameExtractor[REQUEST any] struct {
	Getter DbClientAttrsGetter[REQUEST]
}

func (d *DBSpanNameExtractor[REQUEST]) Extract(request REQUEST) string {
	dbName := d.Getter.GetName(request)
	operation := d.Getter.GetOperation(request)
	if operation == "" {
		if dbName == "" {
			return "DB Query"
		} else {
			return dbName
		}
	}
	return operation
}
