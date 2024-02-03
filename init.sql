CREATE DATABASE app_db;
\c app_db;

CREATE TABLE units (
    name TEXT PRIMARY KEY
);

CREATE TABLE users (
    id BIGINT primary key,
    unit TEXT REFERENCES units(name) ON DELETE CASCADE DEFAULT '',
    is_admin INT DEFAULT 0
);

CREATE TABLE classes (
    day DATE,
    num INT CHECK (num >= 1 AND num <= 8),
    unit TEXT REFERENCES units(name) ON DELETE CASCADE,
    name TEXT,
    room TEXT,
    PRIMARY KEY (day, num, unit)
);

INSERT INTO units(name) VALUES(''); -- for admins;

-- тестовые данные для демонстрации
INSERT INTO units(name) VALUES('h3223');
INSERT INTO units(name) VALUES('p3232');
INSERT INTO units(name) VALUES('q1010');

INSERT INTO classes(day, num, unit, name, room) VALUES(TO_DATE('06.02.2024', 'DD.MM.YYYY'), 1, 'q1010', 'История', '427/1');
INSERT INTO classes(day, num, unit, name, room) VALUES(TO_DATE('06.02.2024', 'DD.MM.YYYY'), 3, 'q1010', 'Философия', '403/2');
INSERT INTO classes(day, num, unit, name, room) VALUES(TO_DATE('07.02.2024', 'DD.MM.YYYY'), 5, 'q1010', 'Базы данных', '504/1');
INSERT INTO classes(day, num, unit, name, room) VALUES(TO_DATE('08.02.2024', 'DD.MM.YYYY'), 1, 'q1010', 'Математический анализ', '427/1');
INSERT INTO classes(day, num, unit, name, room) VALUES(TO_DATE('09.02.2024', 'DD.MM.YYYY'), 2, 'q1010', 'Линейная алгерба', '234/1');


INSERT INTO classes(day, num, unit, name, room) VALUES(TO_DATE('06.02.2024', 'DD.MM.YYYY'), 2, 'p3232', 'Системное программирование', '427/1');
INSERT INTO classes(day, num, unit, name, room) VALUES(TO_DATE('06.02.2024', 'DD.MM.YYYY'), 5, 'p3232', 'Базы данных', '403/2');
INSERT INTO classes(day, num, unit, name, room) VALUES(TO_DATE('07.02.2024', 'DD.MM.YYYY'), 3, 'p3232', 'Веб-программирование', '504/1');
INSERT INTO classes(day, num, unit, name, room) VALUES(TO_DATE('08.02.2024', 'DD.MM.YYYY'), 3, 'p3232', 'Дискретная математика', '427/1');
INSERT INTO classes(day, num, unit, name, room) VALUES(TO_DATE('09.02.2024', 'DD.MM.YYYY'), 4, 'p3232', 'Линейная алгерба', '234/1');


INSERT INTO classes(day, num, unit, name, room) VALUES(TO_DATE('06.02.2024', 'DD.MM.YYYY'), 4, 'h3223', 'Биохимия', '427/1');
INSERT INTO classes(day, num, unit, name, room) VALUES(TO_DATE('06.02.2024', 'DD.MM.YYYY'), 6, 'h3223', 'Математический анализ', '403/2');
INSERT INTO classes(day, num, unit, name, room) VALUES(TO_DATE('07.02.2024', 'DD.MM.YYYY'), 2, 'h3223', 'Физика', '504/1');
INSERT INTO classes(day, num, unit, name, room) VALUES(TO_DATE('08.02.2024', 'DD.MM.YYYY'), 4, 'h3223', 'Аналитическая химия', '427/1');
INSERT INTO classes(day, num, unit, name, room) VALUES(TO_DATE('09.02.2024', 'DD.MM.YYYY'), 3, 'h3223', 'Молекулярная биология', '234/1');
