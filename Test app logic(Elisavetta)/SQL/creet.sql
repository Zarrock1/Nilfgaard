INSERT INTO users (name, blocked)
VALUES('София Августа Фредерика',false),
('Геральт из Ривии', false),
('Йеннифер из Венгерберга', false),
('Трисс Меригольд', false),
('Цирилла Фиона Элен Рианнон', false),
('Ламберт', false),
('Эскель', false),
('Роше', false),
('Золтан Хивай', false),
('Лютик', false),
('Весемир', false),
('Филиппа Эйльхарт', false),
('Шани', false),
('Вильгефорц', false),
('Дуду', false),
('Васк', false),
('Барон', false),
('Ребекка', false),
('Артемия', false),
('Джамшут', false),
('Ясен', false),
('Ярп', false),
('Чиро', false),
('Ренд', false),
('Эредин', false),
('Кагыр', false),
('Ассирэ', false),
('Ядвига', false),
('Калантэ', false),
('Геоффрей', false),
('Николя', false),
('Король Фольтест', false),
('Кахир', false),
('Майкл', false),
('Анна', false),
('Каролина', false),
('Эгр ван Эмрэйс', false),
('Кондвирамурсу',false);


INSERT INTO roles (name)
VALUES ('Админ'),
('Студент'),
('Препод');


INSERT INTO disciplines (title, discription, prepod_id, deleted)
VALUES ('Алгоритмизация и програмирование', 'Всё хорошо, всё под контролем', 1, false),
('Английский', 'Всё плохо, всё не под контролем', 37, false),
('История Нильфгаарда','',38, false);


INSERT INTO tests (title, active, discipline_id, deleted)
VALUES 
('Тест 1', true, 1,false),
('Test 2', true, 2,false),
('Легенды и мифы древнего Нильфгаарда', false, 3, true);

INSERT INTO questions (avtor_id)
VALUES
(1),
(1),
(1),
(1);

INSERT INTO questions_versions ( question_id, title, text_q, version, corect_answer_id)
VALUES 
(1,'Что вполняет эта команда: cd', 'Что вполняет эта команда: cd', 1, 3),
(2,'Есть ли перегрузка функций в go', 'Есть ли перегрузка функций в go', 1, 6),
(3,'What is the capital of Britan?', 'What is the capital of The Greate Britan?', 1, 8),
(4,'Who is the King of Britan?', 'Who is the King of the King of The Greate Britan?', 1, 11),
(2,'Есть ли перегрузка функций в golang?', 'Есть ли перегрузка функций в golang?', 2, 14);



INSERT INTO answers(title, number, question_version_id)
VALUES 
('Такой команды не существует', 1, 1),
('Удаление директории', 2, 1),
('Изменение теущёго каталога', 3, 1),
('Создание каталога', 4, 1),
('Да', 1, 2),
('Нет', 2, 2),
('Parice', 1, 3),
('London',2,3),
('Moscow',3,3),
('Sofia Augusta Frederica',2,4),
('Karl III',3,4),
('I am', 1, 4),
('Да', 1, 5),
('Нет', 2, 5);

INSERT INTO users_roles (user_id, role_id)
VALUES (1,3),
(2, 2),
(3, 3),
(4, 2),
(5, 3),
(6, 2),
(7, 3),
(8, 2),
(9, 3),
(10, 2),
(11, 3),
(12, 2),
(13, 3),
(14, 2),
(15, 3),
(16, 2),
(17, 3),
(18, 2),
(19, 3),
(20, 2),
(21, 3),
(22, 2),
(23, 3),
(24, 2),
(25, 3),
(26, 2),
(27, 3),
(28, 2),
(29, 3),
(30, 2),
(31, 3),
(32, 2),
(33, 3),
(34, 2),
(35, 3),
(36, 2),
(37, 3),
(38, 2),
(1,1),
(27,1),
(37,1),
(20,1);


INSERT INTO users_disciplines (user_id, discipline_id)
VALUES 
(1, 1),
(1, 2);


INSERT INTO tests_questions (test_id, question_id)
VALUES 
(1, 1),
(1, 2),
(2, 3),
(2, 4);

INSERT INTO atempts (user_id, test_id, active)
VALUES
(2, 1, true),
(2, 2, false),
(3, 1, false),
(3, 2, false);

INSERT INTO atempts_questions_answers (atempt_id, question_version_id, answer_id)
VALUES 
(1, 1, 2),
(1, 2, 6),
(2, 3, 8),
(2, 4, 11),
(3, 1, 3),
(3, 5, 14),
(4, 3, 9),
(4, 4, 12);
