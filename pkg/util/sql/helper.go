package sql

import (
	"recipe-app/pkg/domain"

	sq "github.com/Masterminds/squirrel"
)

const (
	tillSfx = "_till"
	fromSfx = "_from"
	sfxLen  = 5
)

func SB() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

type FuncOp func(left, right sq.Sqlizer) sq.Sqlizer

func andOp(left, right sq.Sqlizer) sq.Sqlizer {
	return sq.And{left, right}
}

func orOp(left, right sq.Sqlizer) sq.Sqlizer {
	return sq.Or{left, right}
}

var (
	And FuncOp = andOp
	Or  FuncOp = orOp
)

func newLeft(prevLeft, action sq.Sqlizer, op FuncOp) sq.Sqlizer {
	if prevLeft == nil {
		return action
	}

	return op(prevLeft, action)
}

func AppendReturningID(qs string) string {
	return qs + " RETURNING ID"
}

func opize(dbvMap domain.DBVMap, op FuncOp) sq.Sqlizer {
	var left sq.Sqlizer

	for key, value := range dbvMap {
		kl := len(key)
		if kl > sfxLen {
			switch key[(kl - sfxLen):] {
			case fromSfx:
				left = newLeft(left, sq.GtOrEq{key[:(kl - sfxLen)]: value}, op)
			case tillSfx:
				left = newLeft(left, sq.LtOrEq{key[:(kl - sfxLen)]: value}, op)
			default:
				left = newLeft(left, sq.Eq{key: value}, op)
			}
		} else {
			left = newLeft(left, sq.Eq{key: value}, op)
		}
	}

	return left
}

// FilterSuffixed appends corresponding where clauses into squirrel.SelectBuilder.
// It expects GT/LT fields to have a _from or _till suffix, then it maps it into
// squirrel.GtOrEq/squirrel.LtOrEq clauses.
func FilterSuffixed(builder sq.SelectBuilder, dbvMap domain.DBVMap, op FuncOp) sq.SelectBuilder {
	if clause := opize(dbvMap, op); clause != nil {
		return builder.Where(clause)
	}

	return builder
}
