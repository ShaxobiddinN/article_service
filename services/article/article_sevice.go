// ctrl+fn+f2 belgilangan 1ta sozni hammasini belgilab ozgartirish
package article

import (
	blogpost "blogpost/article_service/protogen/blogpost"
	"blogpost/article_service/storage"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ArticleService...
type articleService struct {
	stg storage.StorageI
	blogpost.UnimplementedArticleServiceServer
}

func NewArticleService(stg storage.StorageI) *articleService {
	return &articleService{
		stg: stg,
	}
}

// Ping ...
func (s *articleService) Ping(ctx context.Context, req *blogpost.Empty) (*blogpost.Pong, error) {
	log.Println("Pingg")
	return &blogpost.Pong{
		Message: "Ok",
	}, nil
}

// CreateArticle...
func (s *articleService) CreateArticle(ctx context.Context, req *blogpost.CreateArticleRequest) (*blogpost.Article, error) {
	id := uuid.New()
	err := s.stg.AddArticle(id.String(), req)
	if err != nil {

		return nil, status.Errorf(codes.Internal, "s.stg.AddArticle: %s", err.Error())
	}

	article, err := s.stg.GetArticleById(id.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleById: %s", err.Error())

	}

	return &blogpost.Article{
		Id:        article.Id,
		Content:   article.Content,
		AuthorId:  article.Author.Id,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}, nil
}

// UpdateArticle...
func (s *articleService) UpdateArticle(ctx context.Context, req *blogpost.UpdateArticleRequest) (*blogpost.Article, error) {
	err := s.stg.UpdateArticle(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.UpdateArticle: %s", err.Error())

	}

	article, err := s.stg.GetArticleById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleById: %s", err.Error())
	}

	return &blogpost.Article{
		Id:        article.Id,
		Content:   article.Content,
		AuthorId:  article.Author.Id,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}, nil
}

// DeleteArticle...
func (s *articleService) DeleteArticle(ctx context.Context, req *blogpost.DeleteArticleRequest) (*blogpost.Article, error) {
	article, err := s.stg.GetArticleById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleById: %s", err.Error())
	}

	err = s.stg.RemoveArticle(article.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.DeleteArticle: %s", err.Error())
	}
	return &blogpost.Article{
		Id:        article.Id,
		Content:   article.Content,
		AuthorId:  article.Author.Id,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}, nil

}

// GetArticleList...
func (s *articleService) GetArticleList(ctx context.Context, req *blogpost.GetArticleListRequest) (*blogpost.GetArticleListResponse, error) {
	fmt.Println("----------GetArticleList----------->")

	res, err := s.stg.GetArticleList(int(req.Offset), int(req.Limit), string(req.Search))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleList: %s", err.Error())
	}

	return res, nil
}

// GetArticleById...
func (s *articleService) GetArticleById(ctx context.Context, req *blogpost.GetArticleByIdRequest) (*blogpost.GetArticleByIdResponse, error) {
	article, err := s.stg.GetArticleById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleById: %s", err.Error())
	}

	return article, nil
}
