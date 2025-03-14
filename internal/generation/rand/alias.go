package rand

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
)

type AliasProvider struct {
	charset string
	len     int
}

func NewAliasProvider(charset string, length int) *AliasProvider {
	return &AliasProvider{
		charset: charset,
		len:     length,
	}
}

func (p *AliasProvider) Generate(ctx context.Context, original string) (string, error) {
	charsetLen := len(p.charset)
	res := make([]byte, p.len)

	maxNum := big.NewInt(int64(charsetLen))

	for i := range p.len {
		num, err := rand.Int(rand.Reader, maxNum)
		if err != nil {
			return "", fmt.Errorf("failed to generate random string: %w", err)
		}

		res[i] = p.charset[num.Int64()]
	}

	return string(res), nil
}
