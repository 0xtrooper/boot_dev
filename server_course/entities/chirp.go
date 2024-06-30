package entities

import (
	"context"
	"regexp"
)

type Chirp struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Body     string `json:"body"`
}

var profaneWords = []string{"kerfuffle", "sharbert", "fornax"}

func (c *Chirp) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)
	if len(c.Body) > 140 {
		problems["too long"] = "message can only be up to and including 140 chars"
	}

	for _, profaneWord := range profaneWords {
		r := regexp.MustCompile("(?i)" + profaneWord)
		c.Body = r.ReplaceAllString(c.Body, "****")
	}

	return problems
}
