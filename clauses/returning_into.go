package clauses

import (
	"github.com/Zone16/gorm/clause"
)

type ReturningInto struct {
	Variables []clause.Column
	Into      []*clause.Values
}
