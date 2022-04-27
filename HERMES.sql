BEGIN;


CREATE TABLE IF NOT EXISTS public.form_store
(
    id integer NOT NULL DEFAULT nextval('form_store_id_seq'::regclass),
    form_id integer NOT NULL,
    owner_id integer NOT NULL,
    total_question integer,
    total_response integer,
    is_active boolean NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    form_name character varying(100) COLLATE pg_catalog."default",
    CONSTRAINT form_store_pkey PRIMARY KEY (form_id)
);

CREATE TABLE IF NOT EXISTS public.job_store
(
    job_id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    job_status text COLLATE pg_catalog."default",
    job_status_code integer NOT NULL,
    plugin_code integer NOT NULL,
    CONSTRAINT job_store_pkey PRIMARY KEY (job_id)
);

CREATE TABLE IF NOT EXISTS public.question_store
(
    question_id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    question text COLLATE pg_catalog."default",
    question_type text COLLATE pg_catalog."default" NOT NULL,
    form_id integer NOT NULL,
    CONSTRAINT question_store_pkey PRIMARY KEY (question_id)
);

CREATE TABLE IF NOT EXISTS public.response_store
(
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    response_id integer NOT NULL,
    response text COLLATE pg_catalog."default",
    response_type text COLLATE pg_catalog."default" NOT NULL,
    question_id integer NOT NULL,
    form_id integer NOT NULL,
    user_id integer NOT NULL,
    CONSTRAINT response_store_pkey PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.question_store
    ADD CONSTRAINT fk_form FOREIGN KEY (form_id)
    REFERENCES public.form_store (form_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;


ALTER TABLE IF EXISTS public.response_store
    ADD CONSTRAINT fk_question FOREIGN KEY (question_id)
    REFERENCES public.question_store (question_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;

COMMIT;