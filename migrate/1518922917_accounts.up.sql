CREATE TABLE IF NOT EXISTS account (
    id varchar(255),
    publickey varchar(255) UNIQUE,
    balance Integer
);
