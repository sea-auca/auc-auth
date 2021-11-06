package user

import "strings"

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

func NewPermissions(perms ...Permissions) *Permissions {
	var perm Permissions = 0
	if len(perms) != 0 {
		for _, v := range perms {
			perm = perm | v
		}
	} else {
		perm = None
	}
	return &perm
}

func (p *Permissions) Assing(perm Permissions) {
	if *p == None {
		*p = perm
		return
	}
	*p = *p | perm
}

func (p *Permissions) Has(perm Permissions) bool {
	return (*p)&perm != 0
}

func (p *Permissions) Revoke(perm Permissions) {
	*p = *p &^ perm
	if *p == 0 {
		*p = None
	}
}

func (p *Permissions) Sprint() string {
	var vals []string
	for key, v := range permissions_map {
		if p.Has(key) {
			vals = append(vals, v)
		}
	}
	return strings.Join(vals, "|")
}
