package repository

import (
	"database/sql"
	"github.com/aalempijevic/communityteaminterview/model"
)

type WordRepo struct {
	db *sql.DB
}

const selectWordsByTag string = `
		select w.word
		from words w
		inner join tag_word_frequencies twf
		on w.id = twf.wordId
		where twf.tagId = ?
		order by frequency desc
	`
//NewWordRepoRepo creates a new repo that uses the database object passed in
func NewWordRepo(database *sql.DB) *WordRepo {
	return &WordRepo{db: database}
}

func (repo *WordRepo) StoreWords(words map[string]bool) error {
	stmt, err := repo.db.Prepare("INSERT IGNORE INTO `words` (`word`) VALUES (?)")
	defer stmt.Close()
	if err != nil {
		return err
	}
	for word := range words {
		_, err = stmt.Exec(word)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo *WordRepo) StoreWordFrequencies(tagWordFrequencies map[int]model.WordFrequencies) interface{} {
	stmt, err := repo.db.Prepare(
		"INSERT IGNORE INTO `tag_word_frequencies` (`wordId`,`tagId`,`frequency`) " +
			"SELECT `id`, ?, ? from `words` WHERE word = ?")
	defer stmt.Close()
	if err != nil {
		return err
	}
	for tagId, wordFrequencies := range tagWordFrequencies {
		for word, frequency := range wordFrequencies {
			_, err = stmt.Exec(tagId, frequency, word)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (repo *WordRepo) GetWordsByTag(tagId int) ([]string, error) {
	rows, err := repo.db.Query(selectWordsByTag, tagId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	words := make([]string, 0)
	for rows.Next() {
		var word string
		if err = rows.Scan(&word); err != nil {
			return words, err
		}
		words = append(words, word)
	}

	return words, nil
}
//Truncate will truncate both the words and tagwordfrequencies tables so we can repopulate
func (repo *WordRepo) Truncate() error {

	statements := []string{
		"ALTER TABLE tag_word_frequencies DROP FOREIGN KEY twf_wordid_fk",
		"TRUNCATE TABLE tag_word_frequencies",
		"TRUNCATE TABLE words",
		"ALTER TABLE tag_word_frequencies ADD CONSTRAINT `twf_wordid_fk` FOREIGN KEY (`wordId`) REFERENCES `words` (`id`) ON DELETE CASCADE ON UPDATE CASCADE",
	}

	for _, statement := range statements {
		if _, err := repo.db.Exec(statement); err != nil {
			return err
		}
	}

	return nil
}