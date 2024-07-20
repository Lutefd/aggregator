-- +goose Up
CREATE UNIQUE INDEX idx_users_api_key ON users(api_key);
CREATE INDEX idx_feeds_user_id ON feeds(user_id);
CREATE UNIQUE INDEX idx_feeds_url ON feeds(url);
CREATE INDEX idx_feed_follows_user_id ON feed_follows(user_id);
CREATE INDEX idx_feed_follows_feed_id ON feed_follows(feed_id);
CREATE UNIQUE INDEX idx_posts_url ON posts(url);
CREATE INDEX idx_posts_feed_id ON posts(feed_id);

-- +goose Down
DROP INDEX idx_users_api_key;
DROP INDEX idx_feeds_user_id;
DROP INDEX idx_feeds_url;
DROP INDEX idx_feed_follows_user_id;
DROP INDEX idx_feed_follows_feed_id;
DROP INDEX idx_posts_url;
DROP INDEX idx_posts_feed_id;
