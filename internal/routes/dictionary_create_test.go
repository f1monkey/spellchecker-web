package routes

import (
	"context"
	"errors"
	"testing"

	f1mspellchecker "github.com/f1monkey/spellchecker"
	"github.com/f1monkey/spellchecker-web/internal/spellchecker"
	"github.com/stretchr/testify/require"
	"github.com/swaggest/usecase/status"
)

type testRegistryAdder struct {
	err error
}

func (f *testRegistryAdder) Add(code string, options spellchecker.Options) (*f1mspellchecker.Spellchecker, error) {
	return nil, f.err
}

func Test_DictionaryCreate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		adder    *testRegistryAdder
		input    DictionaryCreateRequest
		wantErr  bool
		wantCode status.Code
	}{
		{
			name: "success",
			adder: &testRegistryAdder{
				err: nil,
			},
			input: DictionaryCreateRequest{
				Code:      "en",
				Alphabet:  "abcdefghijklmnopqrstuvwxyz",
				MaxErrors: 2,
			},
			wantErr:  false,
			wantCode: status.OK,
		},
		{
			name: "already exists",
			adder: &testRegistryAdder{
				err: spellchecker.ErrAlreadyExists,
			},
			input: DictionaryCreateRequest{
				Code:      "en",
				Alphabet:  "abcdefghijklmnopqrstuvwxyz",
				MaxErrors: 2,
			},
			wantErr:  true,
			wantCode: status.AlreadyExists,
		},
		{
			name: "internal error",
			adder: &testRegistryAdder{
				err: errors.New("boom"),
			},
			input: DictionaryCreateRequest{
				Code:      "fr",
				Alphabet:  "abcdefghijklmnopqrstuvwxyz",
				MaxErrors: 2,
			},
			wantErr:  true,
			wantCode: status.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			interactor := dictionaryCreate(tt.adder)

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
