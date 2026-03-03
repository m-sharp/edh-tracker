package game

import "github.com/m-sharp/edh-tracker/lib/repositories/base"

type Model struct {
	base.ModelBase
	Description string `db:"description"`
	PodID       int    `db:"pod_id"`
	FormatID    int    `db:"format_id"`
}
