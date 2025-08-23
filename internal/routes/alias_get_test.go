package routes

import (
	"context"
	"errors"
	"testing"

	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/stretchr/testify/require"
	"github.com/swaggest/usecase/status"
)

type testAliasGetter struct {
	code string
	err  error
}

func (f *testAliasGetter) GetCodeByAlias(alias string) (string, error) {
	return f.code, f.err
}

func Test_AliasGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		getter   *testAliasGetter
		input    AliasGetRequest
		wantErr  bool
		wantCode status.Code
		wantDict string
	}{
		{
			name:     "success",
			getter:   &testAliasGetter{code: "en", err: nil},
			input:    AliasGetRequest{Alias: "eng"},
			wantErr:  false,
			wantCode: status.OK,
			wantDict: "en",
		},
		{
			name:     "not found",
			getter:   &testAliasGetter{code: "", err: spellchecker.ErrNotFound},
			input:    AliasGetRequest{Alias: "xx"},
			wantErr:  true,
			wantCode: status.NotFound,
		},
		{
			name:     "internal error",
			getter:   &testAliasGetter{code: "", err: errors.New("boom")},
			input:    AliasGetRequest{Alias: "eng"},
			wantErr:  true,
			wantCode: status.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			interactor := aliasGet(tt.getter)

			var out AliasGetResponse
			err := interactor.Interact(context.Background(), tt.input, &out)

			if tt.wantErr {
				require.Error(t, err)
				require.True(t, err.(isErr).Is(tt.wantCode))
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantDict, out.Dictionary)
			}
		})
	}
}
