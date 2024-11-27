--
-- PostgreSQL database dump
--

-- Dumped from database version 16.4 (Debian 16.4-1.pgdg120+1)
-- Dumped by pg_dump version 16.4 (Debian 16.4-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'SQL_ASCII';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: users; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username character varying(50) NOT NULL,
    email character varying(100) NOT NULL,
    password_hash text NOT NULL,
    first_name character varying(50),
    last_name character varying(50),
    date_of_birth date,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    is_active boolean DEFAULT true
);


ALTER TABLE public.users OWNER TO admin;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO admin;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: admin
--

INSERT INTO public.users VALUES (1, 'johndoe', 'johndoe@example.com', 'hash_password_1', 'John', 'Doe', '1990-01-15', '2024-08-25 18:09:24.185387', '2024-08-25 18:09:24.185387', true);
INSERT INTO public.users VALUES (2, 'janedoe', 'janedoe@example.com', 'hash_password_2', 'Jane', 'Doe', '1992-02-25', '2024-08-25 18:09:24.185387', '2024-08-25 18:09:24.185387', true);
INSERT INTO public.users VALUES (3, 'samsmith', 'samsmith@example.com', 'hash_password_3', 'Sam', 'Smith', '1985-03-05', '2024-08-25 18:09:24.185387', '2024-08-25 18:09:24.185387', true);
INSERT INTO public.users VALUES (4, 'emilyjones', 'emilyjones@example.com', 'hash_password_4', 'Emily', 'Jones', '1998-04-20', '2024-08-25 18:09:24.185387', '2024-08-25 18:09:24.185387', true);
INSERT INTO public.users VALUES (5, 'michaelbrown', 'michaelbrown@example.com', 'hash_password_5', 'Michael', 'Brown', '1979-05-30', '2024-08-25 18:09:24.185387', '2024-08-25 18:09:24.185387', true);
INSERT INTO public.users VALUES (6, 'lucyliu', 'lucyliu@example.com', 'hash_password_6', 'Lucy', 'Liu', '1995-06-10', '2024-08-25 18:09:24.185387', '2024-08-25 18:09:24.185387', true);
INSERT INTO public.users VALUES (7, 'chrisevans', 'chrisevans@example.com', 'hash_password_7', 'Chris', 'Evans', '1987-07-22', '2024-08-25 18:09:24.185387', '2024-08-25 18:09:24.185387', true);
INSERT INTO public.users VALUES (8, 'oliverqueen', 'oliverqueen@example.com', 'hash_password_8', 'Oliver', 'Queen', '1991-08-08', '2024-08-25 18:09:24.185387', '2024-08-25 18:09:24.185387', true);
INSERT INTO public.users VALUES (9, 'miawong', 'miawong@example.com', 'hash_password_9', 'Mia', 'Wong', '1993-09-15', '2024-08-25 18:09:24.185387', '2024-08-25 18:09:24.185387', true);
INSERT INTO public.users VALUES (10, 'jackwilson', 'jackwilson@example.com', 'hash_password_10', 'Jack', 'Wilson', '1983-10-05', '2024-08-25 18:09:24.185387', '2024-08-25 18:09:24.185387', true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.users_id_seq', 10, true);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

