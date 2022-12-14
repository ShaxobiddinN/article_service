package postgres

import (
	"blogpost/article_service/protogen/blogpost"
	"errors"
	"time"
)

// AddAuthor...
func (stg Postgres) AddAuthor(id string, entity *blogpost.CreateAuthorRequest) error {

	_, err := stg.db.Exec(`INSERT INTO author 
	(id,fullname) 
	VALUES ($1, $2)`,
		id,
		entity.Fullname,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetAuthorById...
func (stg Postgres) GetAuthorById(id string) (*blogpost.GetAuthorByIdResponse, error) {
	result := &blogpost.GetAuthorByIdResponse{}

	var deletedAt *time.Time
	var updatedAt *string

	err := stg.db.QueryRow(`SELECT 
	id, 
	fullname,
	created_at, 
	updated_at, 
	deleted_at
	FROM author
	WHERE id = $1`, id).Scan(
		&result.Id,
		&result.Fullname,
		&result.CreatedAt,
		&updatedAt,
		&deletedAt,
		// &result.DeleteAt,
	)
	if err != nil {
		return result, err
	}

	if updatedAt != nil {
		result.UpdatedAt = *updatedAt
	}

	if deletedAt != nil {
		return result, errors.New("article not found")
	}

	return result, nil
}

// GetAuthorList...
func (stg Postgres) GetAuthorList(offset, limit int, search string) (*blogpost.GetAuthorListResponse, error) {
	resp := &blogpost.GetAuthorListResponse{
		Authors: make([]*blogpost.Author, 0),
	}
	var deletedAt *time.Time

	rows, err := stg.db.Queryx(`SELECT 
	id, 
	fullname,
	created_at, 
	updated_at, 
	deleted_at 
	FROM author 
	WHERE deleted_at IS NULL AND 
	(fullname ILIKE '%' || $1 || '%') 
	LIMIT $2
	OFFSET $3`, search, limit, offset)
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		var a *blogpost.Author
		var updatedAt *string
		//------------------------bloomrpcda getauthlist qiganda otvorvotti----------------
		err := rows.Scan(
			&a.Id,
			&a.Fullname,
			&a.CreatedAt,
			&updatedAt,
			&deletedAt,
			// &a.DeleteAt,
		)
		if err != nil {
			return resp, err
		}

		if updatedAt != nil {
			a.UpdatedAt = *updatedAt
		}
		resp.Authors = append(resp.Authors, a)
	}
	//comment ochish extimoli bor
	// if deletedAt != nil {
	// 	return resp, errors.New("article not found")
	// }

	return resp, err
}

// UpdateAuthor...
func (stg Postgres) UpdateAuthor(entity *blogpost.UpdateAuthorRequest) error {
	res, err := stg.db.NamedExec(`
	UPDATE  author SET 
		fullname =:fn,
		updated_at=now() 
		WHERE id =:i AND deleted_at IS NULL `, map[string]interface{}{
		"fn": entity.Fullname,
		"i":  entity.Id,
	})
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n > 0 {
		return nil
	}

	return errors.New("author not found")
}

// RemoveAuthor...
func (stg Postgres) RemoveAuthor(id string) error {
	res, err := stg.db.Exec(`UPDATE author SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL`, id)
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

	return errors.New("author not found or already deleted")
}
