package clauses

import (
	"github.com/Zone16/gorm/clause"
)

type WhenNotMatched struct {
	clause.Values
	Where clause.Where
}

func (w WhenNotMatched) Name() string {
	return "WHEN NOT MATCHED"
}

func (w WhenNotMatched) Build(builder clause.Builder) {
	if len(w.Columns) > 0 {
		if len(w.Values.Values) != 1 {
			panic("cannot insert more than one rows due to DM SQL language restriction")
		}

		builder.WriteString(" THEN")
		builder.WriteString(" INSERT ")
		w.Values.Build(builder)

		if len(w.Where.Exprs) > 0 {
			builder.WriteString(w.Where.Name())
			builder.WriteByte(' ')
			w.Where.Build(builder)
		}
	}
}

func (w WhenNotMatched) MergeClause(clause *clause.Clause) {
	clause.Name = w.Name()
	clause.Expression = w
}
