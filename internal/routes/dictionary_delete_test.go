package routes

import (
	"context"
	"errors"
	"testing"

	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/stretchr/testify/require"
	"github.com/swaggest/usecase/status"
)

type isErr interface {
	Is(target error) bool
}

type testDictionaryDeleter struct {
	err error
}

func (f *testDictionaryDeleter) Delete(code string) error {
	return f.err
}

func Test_DictionaryDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		deleter  *testDictionaryDeleter
		input    DictionaryDeleteRequest
		wantErr  bool
		wantCode status.Code
	}{
		{
			name:     "success",
			deleter:  &testDictionaryDeleter{err: nil},
			input:    DictionaryDeleteRequest{Code: "en"},
			wantErr:  false,
			wantCode: status.OK,
		},
		{
			name:     "not found",
			deleter:  &testDictionaryDeleter{err: spellchecker.ErrNotFound},
			input:    DictionaryDeleteRequest{Code: "xx"},
			wantErr:  true,
			wantCode: status.NotFound,
		},
		{
			name:     "internal error",
			deleter:  &testDictionaryDeleter{err: errors.New("boom")},
			input:    DictionaryDeleteRequest{Code: "en"},
			wantErr:  true,
			wantCode: status.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			interactor := dictionaryDelete(tt.deleter)

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
