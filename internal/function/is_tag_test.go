package function

import (
	"context"
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/require"
	"gopkg.in/src-d/go-mysql-server.v0/sql"
	"gopkg.in/src-d/go-mysql-server.v0/sql/expression"
)

func TestIsTag(t *testing.T) {
	f := NewIsTag(expression.NewGetField(0, sql.Text, "name", true))

	testCases := []struct {
		name     string
		row      sql.Row
		expected bool
		err      bool
	}{
		{"null", sql.NewRow(nil), false, false},
		{"not a ref name", sql.NewRow("foo bar"), false, false},
		{"not a tag ref", sql.NewRow("refs/heads/v1.x"), false, false},
		{"a tag", sql.NewRow("refs/tags/v1.0.0"), true, false},
		{"mismatched type", sql.NewRow(1), false, true},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)

			session := sql.NewBaseSession()
			ctx := sql.NewContext(context.TODO(), session, opentracing.NoopTracer{})

			val, err := f.Eval(ctx, tt.row)
			if tt.err {
				require.Error(err)
				require.True(sql.ErrInvalidType.Is(err))
			} else {
				require.NoError(err)
				require.Equal(tt.expected, val)
			}
		})
	}
}
