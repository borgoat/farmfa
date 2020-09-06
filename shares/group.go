package shares

type Group []Token

func (g *Group) Add(token *Token) error {
	newTokens := append(*g, *token)

	if err := validateTokensFromSameSecret(newTokens); err != nil {
		return err
	}

	*g = removeDuplicateTokens(newTokens)

	return nil
}
