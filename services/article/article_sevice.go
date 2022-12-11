// ctrl+fn+f2 belgilangan 1ta sozni hammasini belgilab ozgartirish
package article

import (
	"blogpost/article_service/models"
	blogpost "blogpost/article_service/protogen/blogpost"
	"blogpost/article_service/storage"
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//ArticleService...
type articleService struct {
	stg storage.StorageI
	blogpost.UnimplementedArticleServiceServer
} 

func NewArticleService(stg storage.StorageI) *articleService{
	return &articleService{
		stg: stg,
	}
}


//Ping ...
func (s *articleService) Ping(ctx context.Context, req *blogpost.Empty) (*blogpost.Pong, error) {
	log.Println("Pingg")
	return &blogpost.Pong{
		Message: "Ok",
	}, nil
}

//CreateArticle...
func (s *articleService) CreateArticle(ctx context.Context, req *blogpost.CreateArticleRequest) (*blogpost.Article, error) {
	id := uuid.New()
	err := s.stg.AddArticle(id.String(), models.CreateArticleModel{
		Content: models.Content{
			Title: req.Content.Title,
			Body: req.Content.Body,
		},
		AuthorID: req.AuthorId,
	})
	if err != nil {
		
		return nil, status.Errorf(codes.Internal, "s.stg.AddArticle: %s",err.Error())
	}

	article, err := s.stg.GetArticleById(id.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleById: %s",err.Error())

	}

	var updatedAt string
	if article.UpdateAt !=nil{
		updatedAt = article.UpdateAt.String()
	} 

	return &blogpost.Article{
		Id: article.ID,
		Content: &blogpost.Content{
			Title: article.Title,
			Body: article.Body,
		},
		AuthorId: article.GetAuthor.Id,
		CreatedAt: article.CreatedAt.String(),
		UpdatedAt: updatedAt,
	},nil
}
//UpdateArticle...
func (s *articleService) UpdateArticle(ctx context.Context, req *blogpost.UpdateArticleRequest) (*blogpost.Article, error) {
	err := s.stg.UpdateArticle(models.UpdateArticleModel{
		ID: req.Id,
		Content: models.Content{
			Title: req.Content.Title,
			Body: req.Content.Body,
		},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.UpdateArticle: %s",err.Error())

	}

	article, err := s.stg.GetArticleById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleById: %s",err.Error())
	}

	var updatedAt string
	if article.UpdateAt !=nil{
		updatedAt = article.UpdateAt.String()
	}

	return &blogpost.Article{
		Id: article.ID,
		Content: &blogpost.Content{
			Title: article.Title,
			Body: article.Body,
		},
		AuthorId: article.GetAuthor.Id,
		CreatedAt: article.CreatedAt.String(),
		UpdatedAt: updatedAt,
	},nil
}

//DeleteArticle...
func (s *articleService) DeleteArticle(ctx context.Context, req *blogpost.DeleteArticleRequest) (*blogpost.Article, error) {
	article, err := s.stg.GetArticleById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleById: %s",err.Error())
	}

	var updatedAt string
	if article.UpdateAt !=nil{
		updatedAt = article.UpdateAt.String()
	}


	err = s.stg.RemoveArticle(article.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.DeleteArticle: %s",err.Error())
	}
	
	return &blogpost.Article{
		Id: article.ID,
		Content: &blogpost.Content{
			Title: article.Title,
			Body: article.Body,
		},
		AuthorId: article.GetAuthor.Id,
		CreatedAt: article.CreatedAt.String(),
		UpdatedAt: updatedAt,
	},nil
}
//GetArticleList...
func (s *articleService) GetArticleList(ctx context.Context, req *blogpost.GetArticleListRequest) (*blogpost.GetArticleListResponse, error) {
	 res := &blogpost.GetArticleListResponse{

		 Articles:  make([]*blogpost.Article,0),
	 }


	articleList, err := s.stg.GetArticleList(int(req.Offset), int(req.Limit), string(req.Search))
	if err != nil {
		return nil,status.Errorf(codes.Internal, "s.stg.DeleteArticle: %s",err.Error())
	}
	for _, v := range articleList {
		var updatedAt string
	if v.UpdateAt !=nil{
		updatedAt = v.UpdateAt.String()
	}
		res.Articles = append(res.Articles, &blogpost.Article{
			
				Id: v.ID,
				Content: &blogpost.Content{
					Title: v.Title,
					Body: v.Body,
				},
				AuthorId: v.AuthorID,
				CreatedAt: v.CreatedAt.String(),
				UpdatedAt: updatedAt,
		})
	}
	return res, nil
}

//GetArticleById...
func (s *articleService) GetArticleById(ctx context.Context, req *blogpost.GetArticleByIdRequest) (*blogpost.GetArticleByIdResponse, error) {
	article, err := s.stg.GetArticleById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleById: %s",err.Error())
	}

	if article.DeleteAt != nil{
		return nil, status.Errorf(codes.NotFound, "s.stg.GetArticleById: %s",err.Error())

	}

	var updatedAt string
	if article.UpdateAt !=nil{
		updatedAt = article.UpdateAt.String()
	}
	var authorupdatedAt string
	if article.GetAuthor.UpdateAt !=nil{
		updatedAt = article.GetAuthor.UpdateAt.String()
	}
	return &blogpost.GetArticleByIdResponse{
		Id: article.ID,
		Content: &blogpost.Content{
			Title: article.Title,
			Body: article.Body,
		},
		Author: &blogpost.GetArticleByIdResponse_Author{
			Id: article.GetAuthor.Id,
			Fullname: article.GetAuthor.Fullname,
			CreatedAt: article.CreatedAt.String(),
			UpdatedAt: authorupdatedAt,
		},
		CreatedAt: article.CreatedAt.String(),
		UpdatedAt: updatedAt,
	},nil
}