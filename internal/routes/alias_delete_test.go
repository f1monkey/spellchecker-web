package routes

import (
	"context"
	"errors"
	"testing"

	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/stretchr/testify/require"
	"github.com/swaggest/usecase/status"
)

type testAliasDeleter struct {
	err error
}

func (f *testAliasDeleter) DeleteAlias(alias string) error {
	return f.err
}

func Test_AliasDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		deleter  *testAliasDeleter
		input    AliasDeleteRequest
		wantErr  bool
		wantCode status.Code
	}{
		{
			name:     "success",
			deleter:  &testAliasDeleter{err: nil},
			input:    AliasDeleteRequest{Alias: "eng"},
			wantErr:  false,
			wantCode: status.OK,
		},
		{
			name:     "alias not found",
			deleter:  &testAliasDeleter{err: spellchecker.ErrAliasNotFound},
			input:    AliasDeleteRequest{Alias: "xx"},
			wantErr:  true,
			wantCode: status.NotFound,
		},
		{
			name:     "internal error",
			deleter:  &testAliasDeleter{err: errors.New("boom")},
			input:    AliasDeleteRequest{Alias: "eng"},
			wantErr:  true,
			wantCode: status.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			interactor := aliasDelete(tt.deleter)

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
