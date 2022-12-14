package postgres

import (
	"blogpost/article_service/protogen/blogpost"
	"errors"
	"time"
)

// AddArticle...
func (stg Postgres) AddArticle(id string, entity *blogpost.CreateArticleRequest) error {

	_, err := stg.GetAuthorById(entity.AuthorId)
	if err != nil {
		return err
	}

	if entity.Content == nil {
		entity.Content = &blogpost.Content{}
	}
	_, err = stg.db.Exec(`INSERT INTO article 
	(id, title, body, author_id) 
	VALUES ($1, $2, $3, $4)`,
		id,
		entity.Content.Title,
		entity.Content.Body,
		entity.AuthorId,
	)
	if err != nil {
		return err
	}
	return nil
}

// GetArticleById...
func (stg Postgres) GetArticleById(id string) (*blogpost.GetArticleByIdResponse, error) {
	res := &blogpost.GetArticleByIdResponse{
		Content: &blogpost.Content{},
		Author:  &blogpost.GetArticleByIdResponse_Author{},
	}

	var deletedAt *time.Time
	var updatedAt *string
	var authorUpdatedAt *string

	err := stg.db.QueryRow(`SELECT 
	ar.id,
	ar.title, 
	ar.body, 
	ar.created_at, 
	ar.updated_at, 
	ar.deleted_at,
	au.id, 
	au.fullname,
	au.created_at, 
	au.updated_at
	FROM article AS ar JOIN author AS au ON ar.author_id = au.id WHERE ar.id = $1`, id).Scan(
		&res.Id,
		&res.Content.Title,
		&res.Content.Body,
		&res.CreatedAt,
		&updatedAt,
		&deletedAt,
		&res.Author.Id,
		&res.Author.Fullname,
		&res.Author.CreatedAt,
		&authorUpdatedAt,
		// &res.Author.DeleteAt,
	)
	if err != nil {
		return res, err
	}

	if updatedAt != nil {
		res.UpdatedAt = *updatedAt
	}
	if authorUpdatedAt != nil {
		res.Author.UpdatedAt = *authorUpdatedAt
	}

	if deletedAt != nil {
		return res, errors.New("article not found")
	}
	return res, err
}

// GetArticleList ...
func (stg Postgres) GetArticleList(offset, limit int, search string) (*blogpost.GetArticleListResponse, error) {
	resp := &blogpost.GetArticleListResponse{
		Articles: make([]*blogpost.Article, 0),
	}
	var deletedAt *time.Time

	rows, err := stg.db.Queryx(`SELECT 
	id, 
	title, 
	body, 
	author_id, 
	created_at, 
	updated_at, 
	deleted_at 
	FROM article 
	WHERE deleted_at IS NULL AND (title ILIKE '%' || $1 || '%') OR (body ILIKE '%' || $1 || '%')
	LIMIT $2
	OFFSET $3`, search, limit, offset)
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		a := &blogpost.Article{
			Content: &blogpost.Content{},
		}
		var updatedAt *string

		err := rows.Scan(
			&a.Id,
			&a.Content.Title,
			&a.Content.Body,
			&a.AuthorId,
			&a.CreatedAt,
			&updatedAt,
			&deletedAt,
		)
		if err != nil {
			return resp, err
		}
		if updatedAt != nil {
			a.UpdatedAt = *updatedAt
		}
		resp.Articles = append(resp.Articles, a)
	}
	// if deletedAt != nil {
	// 	return res, errors.New("article not found")
	// }

	return resp, err
}

// UpdateArticle ...
func (stg Postgres) UpdateArticle(entity *blogpost.UpdateArticleRequest) error {
	if entity.Content == nil {
		entity.Content = &blogpost.Content{}
	}
	res, err := stg.db.NamedExec("UPDATE article  SET title=:t, body=:b, updated_at=now() WHERE deleted_at IS NULL AND id=:id", map[string]interface{}{
		"id": entity.Id,
		"t":  entity.Content.Title,
		"b":  entity.Content.Body,
	})
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}

	return errors.New("article not found")
}

// DeleteArticle ...
func (stg Postgres) RemoveArticle(id string) error {
	res, err := stg.db.Exec("UPDATE article  SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}

	return errors.New("article not found or already deleted")

}
