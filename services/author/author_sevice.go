package author

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

// AuthorService...
type authorService struct {
	stg storage.StorageI
	blogpost.UnimplementedAuthorServiceServer
}

func NewAuthorService(stg storage.StorageI) *authorService {
	return &authorService{
		stg: stg,
	}
}

// Ping ...
func (s *authorService) Ping(ctx context.Context, req *blogpost.Empty) (*blogpost.Pong, error) {
	log.Println("Pingg")
	return &blogpost.Pong{
		Message: "Ok",
	}, nil
}

// CreateAuthor...
func (s *authorService) CreateAuthor(ctx context.Context, req *blogpost.CreateAuthorRequest) (*blogpost.Author, error) {
	id := uuid.New()
	err := s.stg.AddAuthor(id.String(), req)
	if err != nil {

		return nil, status.Errorf(codes.Internal, "s.stg.AddAuthor: %s", err.Error())
	}

	author, err := s.stg.GetAuthorById(id.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorById: %s", err.Error())

	}

	return &blogpost.Author{
		Id:        author.Id,
		Fullname:  author.Fullname,
		CreatedAt: author.CreatedAt,
		UpdatedAt: author.UpdatedAt,
	}, nil
}

// UpdateAuthor...
func (s *authorService) UpdateAuthor(ctx context.Context, req *blogpost.UpdateAuthorRequest) (*blogpost.Author, error) {
	err := s.stg.UpdateAuthor(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.UpdateAuthor: %s", err.Error())

	}

	author, err := s.stg.GetAuthorById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorById: %s", err.Error())
	}

	return &blogpost.Author{
		Id:        author.Id,
		Fullname:  author.Fullname,
		CreatedAt: author.CreatedAt,
		UpdatedAt: author.UpdatedAt,
	}, nil
}

// DeleteAuthor...
func (s *authorService) DeleteAuthor(ctx context.Context, req *blogpost.DeleteAuthorRequest) (*blogpost.Author, error) {
	author, err := s.stg.GetAuthorById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorById: %s", err.Error())
	}

	err = s.stg.RemoveAuthor(author.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.DeleteAuthor: %s", err.Error())
	}

	return &blogpost.Author{
		Id:        author.Id,
		Fullname:  author.Fullname,
		CreatedAt: author.CreatedAt,
		UpdatedAt: author.UpdatedAt,
	}, nil
}

// GetAuthorList...
func (s *authorService) GetAuthorList(ctx context.Context, req *blogpost.GetAuthorListRequest) (*blogpost.GetAuthorListResponse, error) {
	fmt.Println("----------GetAuthorList----------->")

	res, err := s.stg.GetAuthorList(int(req.Offset), int(req.Limit), string(req.Search))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorList: %s", err.Error())
	}

	return res, nil
}

// GetAuthorById...
func (s *authorService) GetAuthorById(ctx context.Context, req *blogpost.GetAuthorByIdRequest) (*blogpost.GetAuthorByIdResponse, error) {
	author, err := s.stg.GetAuthorById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorById: %s", err.Error())
	}
	return author, nil
}
