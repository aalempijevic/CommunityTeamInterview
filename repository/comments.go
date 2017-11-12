package repository

import
(
	"github.com/aalempijevic/communityteaminterview/model"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"log"
)

//SelectCommentByWord is a query to get the comments with text matching the search term
//I wrapped the query to add a secondary ordering to make sure that we preserve order of rows
//even if two rows have the same relevance because we are paginating results
const selectCommentsByWord =
	`
		SELECT c.id, c.text
		FROM comments c
			INNER JOIN (
				SELECT c.id, MATCH (text) AGAINST (? IN NATURAL LANGUAGE MODE)  relevance
				FROM comments c
					WHERE MATCH (text) AGAINST (? IN NATURAL LANGUAGE MODE)
			) cs
		ON c.id = cs.id
		ORDER BY cs.relevance desc, cs.id
		LIMIT ?, ?
	`

type CommentRepo struct {
	db *sql.DB
}

//NewCommentRepo creates a new repo that uses the database object passed in
func NewCommentRepo(database *sql.DB) *CommentRepo {
	return &CommentRepo{db: database}
}

//ExtractComments converts rows to comments
func extractComments(rows *sql.Rows) model.Comments {

	comments := model.Comments{}
	for rows.Next() {
		c := model.Comment{}
		err := rows.Scan(&c.Id, &c.Text)
		if err != nil {
			log.Fatal(err)
		}
		comments = append(comments, c)
	}

	return comments
}


//GetCommentsByWord returns all comments that have a match on the search word and paginates the results
func (repo *CommentRepo) GetCommentsByWord(word string, skip int, limit int) model.Comments {
	rows, err := repo.db.Query(selectCommentsByWord, word, word, skip, limit)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	return extractComments(rows)
}
