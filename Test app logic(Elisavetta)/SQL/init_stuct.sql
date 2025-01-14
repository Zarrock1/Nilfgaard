DROP TABLE if exists users;
DROP TABLE if exists roles;
DROP TABLE if exists disciplines;
DROP TABLE if exists tests;
DROP TABLE if exists questions;
DROP TABLE if exists questions_version;
DROP TABLE if exists questions_versions;
DROP TABLE if exists users_answers;
DROP TABLE if exists answers;
DROP TABLE if exists users_roles;
DROP TABLE if exists atemps;
DROP TABLE if exists atemps_questions_answers;
DROP TABLE if exists atempts;
DROP TABLE if exists atempts_questions_answers;
DROP TABLE if exists users_disciplines;
DROP TABLE if exists users_tests;
DROP TABLE if exists tests_questions;



CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name varchar(200),
    blocked boolean DEFAULT false,
    login varchar(200)
);

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name varchar(20)
);

CREATE TABLE disciplines (
    id SERIAL PRIMARY KEY,
    title varchar(500),
    discription varchar(10000),
    prepod_id integer,
    deleted boolean DEFAULT false
);

CREATE TABLE tests(
    id SERIAL PRIMARY KEY,
    title varchar(500),
    active boolean not null,
    discipline_id integer,
    deleted boolean DEFAULT false
);
CREATE TABLE questions(
    id SERIAL PRIMARY KEY,
    avtor_id int,
    deleted boolean DEFAULT false
);
CREATE TABLE questions_versions(
    id SERIAL PRIMARY KEY,
    question_id int,
    title varchar(500),
    text_q varchar(500),
    version integer,
    corect_answer_id integer
);

CREATE TABLE answers(
    id SERIAL PRIMARY KEY,
    title varchar(500),
    number integer,
    question_version_id integer
);

CREATE TABLE users_roles(
    user_id integer,
    role_id integer
);

CREATE TABLE atempts(
    id SERIAL PRIMARY KEY,
    user_id integer,
    test_id integer,
    active boolean DEFAULT true
);
CREATE TABLE atempts_questions_answers(
    atempt_id integer,
    question_version_id integer,
    answer_id  integer
);

CREATE TABLE users_disciplines(
    user_id integer,
    discipline_id integer 
);

CREATE TABLE users_tests(
    user_id integer,
    test_id integer
);

CREATE TABLE tests_questions(
    test_id integer,
    question_id integer,
    q_order integer
);
