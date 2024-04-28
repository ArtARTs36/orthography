package incline

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/artarts36/orthography/incline/word"
)

type Morpher struct {
	host  string
	token string
}

func NewMorpher(host string, token string) *Morpher {
	return &Morpher{
		host:  host,
		token: token,
	}
}

func NewMorpherDefault() *Morpher {
	return NewMorpherWithDefaultHost("")
}

func NewMorpherWithDefaultHost(token string) *Morpher {
	return &Morpher{
		host:  "ws3.morpher.ru",
		token: token,
	}
}

type morpherInclineNounResponse struct {
	Nominative    string `json:"И,omitempty"`
	Genitive      string `json:"Р,omitempty"`
	Dative        string `json:"Д,omitempty"`
	Accusative    string `json:"В,omitempty"`
	Instrumental  string `json:"Т,omitempty"`
	Prepositional string `json:"П,omitempty"`
}

func (m *Morpher) InclineNouns(ctx context.Context, nouns []string) (map[string]*word.Word, error) {
	res := map[string]*word.Word{}
	for _, noun := range nouns {
		w, err := m.inclineNoun(ctx, noun)
		if err != nil {
			return nil, err
		}
		res[noun] = w
	}

	return res, nil
}

func (m *Morpher) inclineNoun(ctx context.Context, noun string) (*word.Word, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, m.buildInclineNounURL(noun), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req) //nolint:bodyclose // already closed
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func(body io.ReadCloser) {
		closeErr := body.Close()
		if closeErr != nil {
			slog.
				With(slog.String("closeErr", closeErr.Error())).
				ErrorContext(ctx, "[orthography][morpher] failed to close request body")
		}
	}(resp.Body)

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var res morpherInclineNounResponse
	err = json.Unmarshal(respBytes, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if res.Nominative == "" {
		res.Nominative = noun
	}

	return &word.Word{
		Nominative:    res.Nominative,
		Genitive:      m.caseString(res.Genitive),
		Dative:        m.caseString(res.Dative),
		Accusative:    m.caseString(res.Accusative),
		Instrumental:  m.caseString(res.Instrumental),
		Prepositional: m.caseString(res.Prepositional),
	}, nil
}

func (m *Morpher) buildInclineNounURL(noun string) string {
	if m.token == "" {
		return fmt.Sprintf(
			"http://%s/russian/declension?s=%s",
			m.host,
			noun,
		)
	}

	return fmt.Sprintf(
		"http://%s/russian/declension?s=%s&token=%s",
		m.host,
		noun,
		m.token,
	)
}

func (m *Morpher) caseString(c string) sql.NullString {
	if c == "" {
		return sql.NullString{
			Valid: false,
		}
	}

	return sql.NullString{
		Valid:  true,
		String: c,
	}
}
