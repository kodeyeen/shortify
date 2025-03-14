package dto

type CreateURLRequest struct {
	Original string `json:"original" validate:"required,url"`
}

type CreateURLResponse struct {
	ID       int64  `json:"-"`
	Original string `json:"original"`
	Alias    string `json:"alias"`
}

type GetURLByAliasRequest struct {
	Alias string `json:"alias"`
}

type GetURLByAliasResponse struct {
	Original string `json:"original"`
	Alias    string `json:"alias"`
}
