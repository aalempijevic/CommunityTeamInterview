package repository

import (
	"database/sql"
	"github.com/aalempijevic/communityteaminterview/model"
	"strings"
	"strconv"
)

type SentenceRepo struct {
	db *sql.DB
}

const selectSentences string = `
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


//GetSentences returns all sentences that have been scored and
// the tag ids for that sentence that exceed the threshold
func (repo *WordRepo) GetSentenceTags(threshold float64) (model.Sentences, error) {
	rows, err := repo.db.Query(selectSentences, threshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return extractSentences(rows)
}


//ExtractSentences converts rows to Sentences
func extractSentences(rows *sql.Rows) (model.Sentences, error) {

	sentences := model.Sentences{}
	for rows.Next() {
		var (
			tagsString   string
			sentenceText string
		)

		err := rows.Scan(&tagsString, &sentenceText)
		if err != nil {
			return sentences, err
		}
		stringTagIds := strings.Split(tagsString, ",")
		tagIds := make([]int, len(stringTagIds))
		for i, idStr := range stringTagIds {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return sentences, err
			}
			tagIds[i] = id
		}

		sentences = append(sentences, model.Sentence{TagIds: tagIds, Text: sentenceText})

	}
	return sentences, nil
}