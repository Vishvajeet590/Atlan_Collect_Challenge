package repository

import (
	"Atlan_Collect_Challenge/entity"
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type ResponseDbSql struct {
	pool       *pgxpool.Pool
	keyDB_pool *pgxpool.Pool
}

var TAG_Response = "Response Repository"

func NewResponseDbSql(pool, keyDB_pool *pgxpool.Pool) *ResponseDbSql {
	return &ResponseDbSql{
		pool:       pool,
		keyDB_pool: keyDB_pool,
	}
}

func (r *ResponseDbSql) Extract(formId int8) (*[]entity.Response, error) {
	log.Printf("%s : ", TAG_Response)
	return nil, nil
}

func (r *ResponseDbSql) Add(response *entity.Response, formId, userId int8) (bool, error) {
	log.Printf("%s : Adding response to DB", TAG_Response)

	//Fetching Response Id from Key DB
	var respId int
	key_tx, err := r.keyDB_pool.Begin(context.Background())
	if err != nil {
		log.Printf("%s Extract : Transaction couldn't be started %v", TAG_Form, err)
		return false, err
	}
	key_rows, err := key_tx.Query(context.Background(), "SELECT keyid from key_store where active IS NULL limit 1;")
	if err != nil {
		return false, err
	}
	for key_rows.Next() {
		err = key_rows.Scan(&respId)
		if err != nil {
			return false, err
		}
	}
	_, err = key_tx.Exec(context.Background(), "UPDATE key_store SET active = true where keyid = $1;", respId)
	if err != nil {
		return false, err
	}
	key_tx.Commit(context.Background())

	var rows = [][]interface{}{}
	for _, r := range response.Responses {
		temp := make([]interface{}, 0)
		temp = append(temp, respId, r.Response, r.ResponseType, r.QuestionId, formId, userId)
		rows = append(rows, temp)
	}

	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		log.Printf("%s Extract : Transaction couldn't be started %v", TAG_Form, err)
		return false, err
	}

	copyCount, err := r.pool.CopyFrom(context.Background(), pgx.Identifier{"response_store"}, []string{"response_id", "response", "response_type", "question_id", "form_id", "user_id"}, pgx.CopyFromRows(rows))
	tx.Commit(context.Background())

	if err != nil {
		log.Printf("%s : Error - %v ", TAG_Form, err)
		return false, err
	}
	if copyCount < 1 {
		log.Printf("%s : Error - No row effected ", TAG_Form)
		return false, errors.New("no row effected")
	}

	return true, nil
}
