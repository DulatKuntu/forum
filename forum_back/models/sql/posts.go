package sqls

import (
	"database/sql"
	"time"

	"awesome_forum/forum_back/models"
)

type PostModel struct {
	DB *sql.DB
}

func (pm *PostModel) Insert(userId int, title, text string, category string) (int, error) {
	stmt := `INSERT INTO posts(userid, title, text, category, createdAt)
		VALUES(?, ?, ?, ?, ?)`

	result, err := pm.DB.Exec(stmt, userId, title, text, category, time.Now())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (pm *PostModel) Get(postid int) (*models.Post, error) {
	stmt := `SELECT rowid, userid, title, text, category, createdAt FROM posts WHERE rowid = ?`
	row := pm.DB.QueryRow(stmt, postid)
	p := &models.Post{}
	err := row.Scan(&p.ID, &p.UserID, &p.Title, &p.Text, &p.Category, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (pm *PostModel) GetAll() ([]*models.Post, error) {
	stmt := `SELECT rowid, userid, title, text, category, createdAt FROM posts`
	rows, err := pm.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	var posts []*models.Post

	for rows.Next() {
		p := &models.Post{}
		err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Text, &p.Category, &p.CreatedAt)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
