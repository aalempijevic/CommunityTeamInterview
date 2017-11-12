package repository

import (
	"github.com/aalempijevic/communityteaminterview/model"
	"database/sql"
	"log"
	"strings"
	"strconv"
)

type WordRepo struct {
	db *sql.DB
}

const selectWordsByTag string =
	`
		select w.word
		from words w
		inner join tag_word_frequencies twf
		on w.id = twf.wordId
		where twf.tagId = ?
		order by frequency desc
	`
const selectSentences string =
	`
		select tagIds, sentence from (
		select
			sentence_tags.tagIds,
			substr(c.text, sentence_tags.annotationStart, sentence_tags.annotationEnd-sentence_tags.annotationStart) sentence
		from (
			-- Select all sentences and their max score (excluding the whole comment scores by checking count)
			select
				cs.commentId,
				cs.annotationStart,
				cs.annotationEnd,
				GROUP_CONCAT(tagId SEPARATOR ',') tagIds
			from comment_scores cs
			where score > ?
			group by cs.commentId, cs.annotationStart, cs.annotationEnd
			having count(*) > 1
		) sentence_tags
		inner join comments c
			on c.id = sentence_tags.commentId
		) s
		where -- Exclude white space only sentences
		NULLIF(REPLACE(REPLACE(REPLACE(sentence, ' ', ''), '\t', ''), '\n', ''), ' ') IS NOT NULL
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
	for word, _:= range words {
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

func (repo *WordRepo) GetWordsByTag(tagId int) []string {
	rows, err := repo.db.Query(selectWordsByTag, tagId)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	words := make([]string, 0)
	for rows.Next() {
		var word string
		rows.Scan(&word)
		words = append(words, word)
	}

	return words
}

//GetSentences returns all sentences that have been scored and
// the tag ids for that sentence that exceed the threshold
func (repo *WordRepo) GetSentenceTags(threshold float32) model.Sentences {
	rows, err := repo.db.Query(selectSentences, threshold)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	return extractSentences(rows)
}


//ExtractSentences converts rows to Sentences
func extractSentences(rows *sql.Rows) model.Sentences {

	sentences := model.Sentences{}
	for rows.Next() {
		var (
			tagsString string
			sentenceText string
		)

		err := rows.Scan(&tagsString, &sentenceText)
		if err != nil {
			log.Fatal(err)
		}
		stringTagIds := strings.Split(tagsString, ",")
		tagIds := make([]int, len(stringTagIds))
		for i, idStr := range stringTagIds {
			id,err := strconv.Atoi(idStr)
			if err != nil {
				log.Fatal(err)
			}
			tagIds[i] = id
		}

		sentences = append(sentences, model.Sentence{TagIds:tagIds, Text:sentenceText})

	}
	return sentences
}
