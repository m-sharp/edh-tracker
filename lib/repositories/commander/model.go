package commander

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

type Model struct {
	base.ModelBase
	Name string `db:"name"`
}
