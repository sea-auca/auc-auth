package model

import "fmt"

type Permissions uint64

const (
	None Permissions = 1 << iota
	AUCA_student
	SEA_moderating
	SEA_development
	SEA_maintanance
)

var permissions_map = map[Permissions]string{
	AUCA_student:    "Auca student",
	SEA_maintanance: "AU cloud engineer",
	SEA_development: "SEA certified developer",
	SEA_moderating:  "SEA club moderator",
}

func NewPermissions() *Permissions {
	var perm Permissions = 0
	return &perm
}

func (p *Permissions) Assing(perm Permissions) {
	*p = *p | perm
}

func (p *Permissions) Has(perm Permissions) bool {
	return (*p)&perm != 0
}

func (p *Permissions) Revoke(perm Permissions) {
	*p = *p &^ perm
}

func (p *Permissions) Print() {
	for _, v := range []Permissions{AUCA_student, SEA_moderating, SEA_development, SEA_maintanance} {
		if p.Has(v) {
			fmt.Println(permissions_map[v])
		}
	}
}
