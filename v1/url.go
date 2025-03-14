package shortify

type CreateURLRequest struct {
	Original string `json:"original" validate:"required,url"`
}

type CreateURLResponse struct {
	Original string `json:"original"`
	Alias    string `json:"alias"`
}

type GetURLByAliasRequest struct {
}

type GetURLByAliasResponse struct {
	Original string `json:"original"`
	Alias    string `json:"alias"`
}
