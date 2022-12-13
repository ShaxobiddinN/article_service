package author

import (
	"blogpost/article_service/models"
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
	err := s.stg.AddAuthor(id.String(), models.CreateAuthorModel{
		Fullname: req.Fullname,
	})
	if err != nil {

		return nil, status.Errorf(codes.Internal, "s.stg.AddAuthor: %s", err.Error())
	}

	author, err := s.stg.GetAuthorById(id.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorById: %s", err.Error())

	}

	var updatedAt string
	if author.UpdateAt != nil {
		updatedAt = author.UpdateAt.String()
	}

	return &blogpost.Author{
		Id:        author.Id,
		Fullname:  author.Fullname,
		CreatedAt: author.CreatedAt.String(),
		UpdatedAt: updatedAt,
	}, nil
}

// UpdateAuthor...
func (s *authorService) UpdateAuthor(ctx context.Context, req *blogpost.UpdateAuthorRequest) (*blogpost.Author, error) {
	err := s.stg.UpdateAuthor(models.UpdateAuthorModel{
		Id:       req.Id,
		Fullname: req.Fullname,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.UpdateAuthor: %s", err.Error())

	}

	author, err := s.stg.GetAuthorById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorById: %s", err.Error())
	}

	var updatedAt string
	if author.UpdateAt != nil {
		updatedAt = author.UpdateAt.String()
	}

	return &blogpost.Author{
		Id:        author.Id,
		Fullname:  author.Fullname,
		CreatedAt: author.CreatedAt.String(),
		UpdatedAt: updatedAt, //.String() bolishi mumkin
	}, nil
}

// DeleteAuthor...
func (s *authorService) DeleteAuthor(ctx context.Context, req *blogpost.DeleteAuthorRequest) (*blogpost.Author, error) {
	author, err := s.stg.GetAuthorById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorById: %s", err.Error())
	}

	var updatedAt string
	if author.UpdateAt != nil {
		updatedAt = author.UpdateAt.String()
	}

	err = s.stg.RemoveAuthor(author.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.DeleteAuthor: %s", err.Error())
	}

	return &blogpost.Author{
		Id:        author.Id,
		Fullname:  author.Fullname,
		CreatedAt: author.CreatedAt.String(),
		UpdatedAt: updatedAt,
	}, nil
}

// GetAuthorList...
func (s *authorService) GetAuthorList(ctx context.Context, req *blogpost.GetAuthorListRequest) (*blogpost.GetAuthorListResponse, error) {
	fmt.Println("----------GetAuthorList----------->")

	res := &blogpost.GetAuthorListResponse{

		Authors: make([]*blogpost.Author, 0),
	}

	authorList, err := s.stg.GetAuthorList(int(req.Offset), int(req.Limit), string(req.Search))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.DeleteAuthor: %s", err.Error())
	}
	for _, v := range authorList {
		var updatedAt string
		if v.UpdateAt != nil {
			updatedAt = v.UpdateAt.String()
		}
		res.Authors = append(res.Authors, &blogpost.Author{

			Id:        v.Id,
			Fullname:  v.Fullname,
			CreatedAt: v.CreatedAt.String(),
			UpdatedAt: updatedAt,
		})
	}
	return res, nil
}

// GetAuthorById...
func (s *authorService) GetAuthorById(ctx context.Context, req *blogpost.GetAuthorByIdRequest) (*blogpost.GetAuthorByIdResponse, error) {
	author, err := s.stg.GetAuthorById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorById: %s", err.Error())
	}

	if author.DeleteAt != nil {
		return nil, status.Errorf(codes.NotFound, "s.stg.GetAuthorById: author with id: %s not found", req.Id)

	}

	var updatedAt string
	if author.UpdateAt != nil {
		updatedAt = author.UpdateAt.String()
	}
	// var authorupdatedAt string
	// if author.GetAuthor.UpdateAt != nil {
	// 	updatedAt = author.GetAuthor.UpdateAt.String()
	// }
	return &blogpost.GetAuthorByIdResponse{
		Id:        author.Id,
		Fullname:  author.Fullname,
		CreatedAt: author.CreatedAt.String(),
		UpdatedAt: updatedAt,
	}, nil
}
