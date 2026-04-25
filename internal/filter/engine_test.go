package filter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompileAcceptsStandardTagEqualityPredicate(t *testing.T) {
	t.Parallel()

	engine, err := NewEngine(NewSchema())
	require.NoError(t, err)

	_, err = engine.Compile(context.Background(), `tags.exists(t, t == "1231")`)
	require.NoError(t, err)
}

func TestCompileRejectsLegacyNumericLogicalOperand(t *testing.T) {
	t.Parallel()

	engine, err := NewEngine(NewSchema())
	require.NoError(t, err)

	_, err = engine.Compile(context.Background(), `pinned && 1`)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to compile filter")
}

func TestCompileRejectsNonBooleanTopLevelConstant(t *testing.T) {
	t.Parallel()

	engine, err := NewEngine(NewSchema())
	require.NoError(t, err)

	_, err = engine.Compile(context.Background(), `1`)
	require.EqualError(t, err, "filter must evaluate to a boolean value")
}

func TestCompileToStatementSupportsUIDContains(t *testing.T) {
	t.Parallel()

	engine, err := NewEngine(NewSchema())
	require.NoError(t, err)

	stmt, err := engine.CompileToStatement(context.Background(), `uid.contains("searchable")`, RenderOptions{
		Dialect: DialectSQLite,
	})
	require.NoError(t, err)
	require.Contains(t, stmt.SQL, "memo`.`uid")
	require.Len(t, stmt.Args, 1)
	require.Equal(t, "%searchable%", stmt.Args[0])
}

func TestCompileToStatementContentContainsIncludesCommentSearch(t *testing.T) {
	t.Parallel()

	engine, err := NewEngine(NewSchema())
	require.NoError(t, err)

	stmt, err := engine.CompileToStatement(context.Background(), `content.contains("#SB")`, RenderOptions{
		Dialect: DialectSQLite,
	})
	require.NoError(t, err)
	require.Contains(t, stmt.SQL, `comment_relation`)
	require.Contains(t, stmt.SQL, `comment_memo`)
	require.Len(t, stmt.Args, 2)
	require.Equal(t, "%#SB%", stmt.Args[0])
	require.Equal(t, "%#SB%", stmt.Args[1])
}
