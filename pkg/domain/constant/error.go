package constant

// Postgres errorCodes.
const (
	PgUniqueConstraintViolation    = "23505"
	PgForeignKeyViolation          = "23503"
	PgCheckConstraint              = "23514"
	PgNotNullViolation             = "23502"
	PgIntegrityConstraintViolation = "23000"
	PgRestrictViolation            = "23001"
	PgRelationDoesNotExist         = "41P01"
)
