package migrations

import (
	r "github.com/go-rel/rel"
)

const UUID r.ColumnType = "UUID"

func MigrateCreateUsers(schema *r.Schema) {
	schema.Exec("CREATE SCHEMA IF NOT EXISTS user_space")
	schema.CreateTableIfNotExists("user_space.users", func(t *r.Table) {
		t.Column("uuid", UUID, r.Primary(true))
		t.String("email", r.Required(true), r.Unique(true))
		t.String("hash")
		t.String("fullname", r.Required(true))
		t.BigInt("permissions", r.Default(0), r.Required(true))
		t.Bool("active", r.Default(true), r.Required(true))
		t.Bool("verified", r.Default(false), r.Required(true))
		t.DateTime("created_at", r.Required(true))
		t.DateTime("updated_at", r.Required(true))
		t.PrimaryKey("uuid")
	})

	schema.CreateUniqueIndex("user_space.users", "UI_users_id", []string{"uuid"})
	schema.CreateUniqueIndex("user_space.users", "UI_users_email", []string{"email"})
}

func RollbackCreateUsers(schema *r.Schema) {
	schema.DropIndex("user_space.users", "user_space.UI_users_id")
	schema.DropIndex("user_space.users", "user_space.UI_users_email")
	schema.DropTableIfExists("user_space.users")
	schema.Exec("DROP SCHEMA IF EXISTS user_space")
}
