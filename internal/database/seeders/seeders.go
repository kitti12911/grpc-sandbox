package seeders

import "embed"

const Users = "users.yml"

//go:embed *.yml
var Fixtures embed.FS
