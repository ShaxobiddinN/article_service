package author

import (
	blogpost "blogpost/article_service/protogen/blogpost"
	"context"
	"log"
)

type AuthorService struct {
	blogpost.UnimplementedAuthorServiceServer
} 


//Ping ...
func (s *AuthorService) Ping(ctx context.Context, req *blogpost.Empty) (*blogpost.Pong, error) {
	log.Println("Pingg")
	return &blogpost.Pong{
		Message: "Ok",
	}, nil
}