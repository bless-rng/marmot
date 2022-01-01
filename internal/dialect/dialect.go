package dialect

type SqlDialect interface {
	CreateVersionTableSql() string
	AddVersionSql(name string) string
	VerifyVersion(name string) string
	FindMigration(name string) string
	DeleteVersionSql(name string) string
	AllMigrationsSql() string
}
