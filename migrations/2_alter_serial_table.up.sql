-- Создание последовательности
CREATE SEQUENCE users_id_seq;

-- Изменение колонки id для использования этой последовательности
ALTER TABLE users ALTER COLUMN id SET DEFAULT nextval('users_id_seq');

-- Установка типа данных bigint
ALTER TABLE users ALTER COLUMN id SET DATA TYPE bigint;

