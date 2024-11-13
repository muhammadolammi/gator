-- +goose Up
CREATE TABLE posts(
id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
 created_at TIMESTAMP NOT NULL,
 updated_at TIMESTAMP NOT NULL,
 title TEXT NOT NULL,
 url TEXT UNIQUE NOT NULL,
 description TEXT  NULL,
 published_at TIMESTAMP  NOT  NULL, 
 feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE

);

-- +goose Down
DROP TABLE posts;