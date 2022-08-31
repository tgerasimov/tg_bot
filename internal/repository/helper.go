package repository

import "github.com/Masterminds/squirrel"

func pgQb() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}
