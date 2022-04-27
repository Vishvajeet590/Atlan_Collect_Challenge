
BEGIN;

SET client_encoding = 'LATIN1';

CREATE TABLE form_store
(
id INT GENERATED ALWAYS AS IDENTITY,
form_id INT NOT NULL,
owner_id INT NOT NULL,
total_question INT,
total_response INT,
is_active BOOLEAN NOT NULL,
created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
form_name character varying(100) NOT NULL,
PRIMARY KEY(form_id)
);



CREATE TABLE question_store(
	question_id INT GENERATED ALWAYS AS IDENTITY,
	question text,
	question_type text NOT NULL,
	form_id INT NOT NULL,
	PRIMARY KEY(question_id),
	CONSTRAINT fk_question
		FOREIGN KEY(form_id)
			REFERENCES form_store(form_id)
			ON DELETE CASCADE
);


CREATE TABLE response_store(
	id INT GENERATED ALWAYS AS IDENTITY,
	response_id INT NOT NULL,
	response text,
	response_type text NOT NULL,
	question_id INT NOT NULL,
	form_id integer NOT NULL,
	user_id integer NOT NULL,
	PRIMARY KEY(id),
	CONSTRAINT fk_response
		FOREIGN KEY(question_id)
			REFERENCES question_store(question_id)
			ON DELETE CASCADE
);



CREATE TABLE job_store(
	job_id INT GENERATED ALWAYS AS IDENTITY,
	job_status text,
	job_status_code INT NOT NULL,
	plugin_code INT NOT NULL,
	PRIMARY KEY(job_id)
);
COMMIT;