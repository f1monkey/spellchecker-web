package routes

import (
	"context"
	"errors"
	"testing"

	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/stretchr/testify/require"
	"github.com/swaggest/usecase/status"
)

type testDictionarySaver struct {
	err error
}

func (f *testDictionarySaver) Save(code string) error {
	return f.err
}

func Test_DictionarySave(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		saver    *testDictionarySaver
		input    DictionarySaveRequest
		wantErr  bool
		wantCode status.Code
	}{
		{
			name:     "success",
			saver:    &testDictionarySaver{err: nil},
			input:    DictionarySaveRequest{Code: "en"},
			wantErr:  false,
			wantCode: status.OK,
		},
		{
			name:     "not found",
			saver:    &testDictionarySaver{err: spellchecker.ErrNotFound},
			input:    DictionarySaveRequest{Code: "xx"},
			wantErr:  true,
			wantCode: status.NotFound,
		},
		{
			name:     "internal error",
			saver:    &testDictionarySaver{err: errors.New("boom")},
			input:    DictionarySaveRequest{Code: "en"},
			wantErr:  true,
			wantCode: status.Internal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			interactor := dictionarySave(tt.saver)

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
