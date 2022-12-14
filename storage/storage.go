package storage

import (
	"blogpost/article_service/protogen/blogpost"
)

type StorageI interface {
	AddArticle(id string, entity *blogpost.CreateArticleRequest) error
	GetArticleById(id string) (*blogpost.GetArticleByIdResponse, error)
	GetArticleList(offset, limit int, search string) (resp *blogpost.GetArticleListResponse, err error)
	UpdateArticle(entity *blogpost.UpdateArticleRequest) error
	RemoveArticle(id string) error

	AddAuthor(id string, entity *blogpost.CreateAuthorRequest) error
	GetAuthorById(id string) (*blogpost.GetAuthorByIdResponse, error)
	GetAuthorList(offset, limit int, search string) (resp *blogpost.GetAuthorListResponse, err error)
	UpdateAuthor(entity *blogpost.UpdateAuthorRequest) error
	RemoveAuthor(id string) error
}
