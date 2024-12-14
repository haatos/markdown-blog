-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS roles (
    id INTEGER PRIMARY KEY,
    name TEXT
);

INSERT INTO roles (id, name) VALUES (1, 'free');
INSERT INTO roles (id, name) VALUES (1000, 'superuser');

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id INTEGER,
    first_name TEXT,
    last_name TEXT,
    email TEXT,
    avatar_url TEXT,
    FOREIGN KEY (role_id) REFERENCES roles (id),
    UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id INTEGER,
    expires TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    title TEXT,
    slug TEXT,
    description TEXT,
    content TEXT,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published_on TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL,
    UNIQUE (title),
    UNIQUE (slug)
);

CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    UNIQUE (name)
);

INSERT INTO tags (name) VALUES ('algorithms');
INSERT INTO tags (name) VALUES ('data structures');
INSERT INTO tags (name) VALUES ('tutorials');
INSERT INTO tags (name) VALUES ('golang');

CREATE TABLE IF NOT EXISTS articles_tags (
    tag_id INTEGER,
    article_id INTEGER,
    FOREIGN KEY (tag_id) REFERENCES tags (id),
    FOREIGN KEY (article_id) REFERENCES articles (id)
);

CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    article_id INTEGER,
    user_id INTEGER,
    comment_id INTEGER,
    content TEXT,
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (article_id) REFERENCES articles (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS images (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    image_key TEXT,
    UNIQUE (name)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE images;
DROP TABLE article_comments;
DROP TABLE articles_tags;
DROP TABLE tags;
DROP TABLE articles;
DROP TABLE sessions;
DROP TABLE users;
DROP TABLE roles;
-- +goose StatementEnd
