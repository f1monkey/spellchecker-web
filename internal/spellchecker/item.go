package spellchecker

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/f1monkey/spellchecker"
)

type RegistryItem struct {
	Spellchecker *spellchecker.Spellchecker
	Options      Options
}

type Options struct {
	Alphabet  string `json:"alphabet"`
	MaxErrors uint   `json:"maxErrors"`
}

type src struct {
	Options      Options `json:"options"`
	Spellchecker []byte  `json:"spellchecker"`
}

func (r *RegistryItem) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if r.Spellchecker != nil {
		if err := r.Spellchecker.Save(&buf); err != nil {
			return nil, err
		}
	}

	return json.Marshal(src{
		Options:      r.Options,
		Spellchecker: buf.Bytes(),
	})
}

func (r *RegistryItem) UnmarshalJSON(data []byte) error {
	var value src

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	if len(value.Spellchecker) == 0 {
		return fmt.Errorf("unable to initialize spellchecker with an empty slice")
	}

	sc, err := spellchecker.Load(bytes.NewReader(value.Spellchecker))
	if err != nil {
		return err
	}

	r.Spellchecker = sc
	r.Options = value.Options

	return nil
}
