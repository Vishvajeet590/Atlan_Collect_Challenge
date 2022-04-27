package repository

import (
	"Atlan_Collect_Challenge/entity"
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type FormDbSql struct {
	pool       *pgxpool.Pool
	keyDB_pool *pgxpool.Pool
}

var TAG_Form = "Form Repository"

func NewFormDbSql(pool, keyDB_pool *pgxpool.Pool) *FormDbSql {
	return &FormDbSql{
		pool:       pool,
		keyDB_pool: keyDB_pool,
	}
}

func (r *FormDbSql) Extract(formId int8) (*entity.Form, error) {
	log.Printf("%s : Fetching Form from DB", TAG_Form)

	//SELECT question_id,question,question_type FROM question_store where form_id = 98
	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		log.Printf("%s Extract : Transaction couldn't be started %v", TAG_Form, err)
		return nil, err
	}

	rows, err := tx.Query(context.Background(), "SELECT question_id,question,question_type FROM question_store where form_id = $1", formId)
	if err != nil {
		log.Printf("%s Extract : %v", TAG_Form, err)
		return nil, err
	}

	questions := make([]entity.Question, 0)
	for rows.Next() {
		var question_id int
		var question string
		var question_type string
		err = rows.Scan(&question_id, &question, &question_type)
		if err != nil {
			return nil, err
		}
		temp := entity.Question{
			Question:     question,
			QuestionType: question_type,
			QuestionId:   question_id,
		}
		questions = append(questions, temp)

	}

	rows.Close()
	tx.Commit(context.Background())

	//Form Name Fetch
	tx, err = r.pool.Begin(context.Background())
	row := tx.QueryRow(context.Background(), "SELECT form_name,owner_id FROM form_store where form_id = $1", formId)
	var form_name string
	var owner_id int
	err = row.Scan(&form_name, &owner_id)
	if err != nil {
		return nil, err
	}
	tx.Commit(context.Background())

	form := &entity.Form{
		FormName: form_name,
		FormId:   formId,
		OwnerId:  owner_id,
		Question: questions,
	}

	return form, nil
}

func (r *FormDbSql) Add(form *entity.Form) (bool, int, error) {

	//Fetching Response Id from Key DB
	var formId int
	key_tx, err := r.keyDB_pool.Begin(context.Background())
	if err != nil {
		log.Printf("%s Extract : Transaction couldn't be started %v", TAG_Form, err)
		return false, -999, err
	}
	key_rows, err := key_tx.Query(context.Background(), "SELECT keyid from key_store where active IS NULL limit 1;")
	if err != nil {
		return false, -999, err
	}
	for key_rows.Next() {
		err = key_rows.Scan(&formId)
		if err != nil {
			return false, -999, err
		}
	}
	_, err = key_tx.Exec(context.Background(), "UPDATE key_store SET active = true where keyid = $1;", formId)
	if err != nil {
		return false, -999, err
	}
	key_tx.Commit(context.Background())

	log.Printf("%s : Adding form to DB", TAG_Form)
	var rows = [][]interface{}{}
	for _, q := range form.Question {
		temp := make([]interface{}, 0)
		temp = append(temp, q.Question, q.QuestionType, formId)
		rows = append(rows, temp)
	}

	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		log.Printf("%s Add : Transaction couldn't be started", TAG_Form)
		return false, -999, err
	}

	//INSERT INTO form_store(form_id,owner_id,form_name,total_question,total_response,is_active) values ($1,$2,$3,$4,$5,$6)
	_, err = tx.Exec(context.Background(),
		"INSERT INTO form_store(form_id,owner_id,form_name,total_question,total_response,is_active) values ($1,$2,$3,$4,$5,$6)",
		formId, form.OwnerId, form.FormName, len(form.Question), 0, true,
	)

	tx.Commit(context.Background())

	tx, err = r.pool.Begin(context.Background())
	if err != nil {
		log.Printf("%s Add : Error - %v ", TAG_Form, err)
		return false, -999, err
	}
	copyCount, err := r.pool.CopyFrom(context.Background(), pgx.Identifier{"question_store"}, []string{"question", "question_type", "form_id"}, pgx.CopyFromRows(rows))
	tx.Commit(context.Background())

	if err != nil {
		log.Printf("%s : Error - %v ", TAG_Form, err)
		return false, -999, err
	}
	if copyCount < 1 {
		log.Printf("%s : Error - No row effected ", TAG_Form)
		return false, -999, errors.New("no row effected")
	}

	return true, formId, nil
}

func (r *FormDbSql) Delete(formId int8) (bool, error) {

	//TODO Add functionality
	return true, nil

}
