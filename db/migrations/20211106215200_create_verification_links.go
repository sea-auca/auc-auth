package migrations

import r "github.com/go-rel/rel"

func MigrateCreateVerificationLinks(schema *r.Schema) {
	schema.CreateTable("user_space.verify_links", func(t *r.Table) {
		t.BigID("id", r.Primary(true))
		t.Column("user_id", UUID)
		t.String("link", r.Required(true), r.Unique(true))
		t.DateTime("expires_at", r.Required(true))
		t.DateTime("created_at", r.Required(true))
		t.DateTime("updated_at", r.Required(true))
		t.Bool("is_password_reset", r.Required(true), r.Default(false))
		t.ForeignKey("user_id", "user_space.users", "uuid")
		t.Fragment("CHECK (expires_at > created_at)")
	})
}

func RollbackCreateVerificationLinks(schema *r.Schema) {
	schema.DropTable("user_space.verify_links")
}
