//ctrl+fn+f2 belgilangan 1ta sozni hammasini belgilab ozgartirish 
package article

import (
	blogpost "blogpost/article_service/protogen/blogpost"
	"context"
	"log"
)

type ArticleService struct {
	blogpost.UnimplementedArticleServiceServer
} 


//Ping ...
func (s *ArticleService) Ping(ctx context.Context, req *blogpost.Empty) (*blogpost.Pong, error) {
	log.Println("Pingg")
	return &blogpost.Pong{
		Message: "Ok",
	}, nil
}