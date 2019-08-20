package grimoire

type Limit int

func (l Limit) Build(query *Query) {
	query.LimitClause = l
}