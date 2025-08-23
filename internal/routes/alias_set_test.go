package routes

import (
	"context"
	"errors"
	"testing"

	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/stretchr/testify/require"
	"github.com/swaggest/usecase/status"
)

type testAliasSetter struct {
	err error
}

func (f *testAliasSetter) SetAlias(alias string, code string) error {
	return f.err
}

func Test_AliasSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setter   *testAliasSetter
		input    AliasSetRequest
		wantErr  bool
		wantCode status.Code
	}{
		{
			name:     "success",
			setter:   &testAliasSetter{err: nil},
			input:    AliasSetRequest{Alias: "eng", Dictionary: "en"},
			wantErr:  false,
			wantCode: status.OK,
		},
		{
			name:     "alias not found",
			setter:   &testAliasSetter{err: spellchecker.ErrAliasNotFound},
			input:    AliasSetRequest{Alias: "xx", Dictionary: "en"},
			wantErr:  true,
			wantCode: status.NotFound,
		},
		{
			name:     "internal error",
			setter:   &testAliasSetter{err: errors.New("boom")},
			input:    AliasSetRequest{Alias: "eng", Dictionary: "en"},
			wantErr:  true,
			wantCode: status.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			interactor := aliasSet(tt.setter)

			var out Empty
			err := interactor.Interact(context.Background(), tt.input, &out)

			if tt.wantErr {
				require.Error(t, err)
				require.True(t, err.(isErr).Is(tt.wantCode))
			} else {
				require.NoError(t, err)
			}
		})
	}
}
