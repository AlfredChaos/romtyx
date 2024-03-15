-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE `users` (
    `id` BIGINT NOT NULL,
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `last_sign_in` datetime DEFAULT NULL,
    `user_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `email` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `password_hash` varbinary(255) DEFAULT NULL,
    PRIMARY KEY (`id`)
);
CREATE TABLE `files` (
    `id` BIGINT NOT NULL,
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `path` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `last_version` BIGINT NOT NULL,
    `file_hash` varbinary(255) DEFAULT NULL,
    `owner_id` BIGINT NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (owner_id) REFERENCES users(id),
    FOREIGN KEY (last_version) REFERENCES file_versions(id)
);
CREATE TABLE `file_versions` (
    `id` BIGINT NOT NULL,
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `file_id` BIGINT NOT NULL,
    `version_number` BIGINT NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (file_id) REFERENCES files(id)
);
CREATE TABLE `devices` (
    `id` BIGINT NOT NULL,
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `user_id` BIGINT NOT NULL,
    `name` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS `users`;
DROP TABLE IF EXISTS `files`;
DROP TABLE IF EXISTS `file_versions`;
DROP TABLE IF EXISTS `devices`;
-- +goose StatementEnd