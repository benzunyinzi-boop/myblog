-- 000001_init_schema.up.sql
-- myblog 初始 schema

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ============== users ==============
CREATE TABLE IF NOT EXISTS `users` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `username`      VARCHAR(32)     NOT NULL,
  `email`         VARCHAR(128)    NOT NULL DEFAULT '',
  `password_hash` VARCHAR(128)    NOT NULL,
  `nickname`      VARCHAR(64)     NOT NULL DEFAULT '',
  `avatar`        VARCHAR(255)    NOT NULL DEFAULT '',
  `bio`           VARCHAR(255)    NOT NULL DEFAULT '',
  `role`          VARCHAR(16)     NOT NULL DEFAULT 'author',
  `status`        TINYINT         NOT NULL DEFAULT 1,
  `created_at`    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at`    DATETIME(3)     NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_username` (`username`),
  UNIQUE KEY `uniq_email`    (`email`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '用户';

-- ============== categories ==============
CREATE TABLE IF NOT EXISTS `categories` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name`        VARCHAR(64)     NOT NULL,
  `slug`        VARCHAR(64)     NOT NULL,
  `description` VARCHAR(255)    NOT NULL DEFAULT '',
  `sort_order`  INT             NOT NULL DEFAULT 0,
  `created_at`  DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`  DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at`  DATETIME(3)     NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_name` (`name`),
  UNIQUE KEY `uniq_slug` (`slug`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '分类';

-- ============== tags ==============
CREATE TABLE IF NOT EXISTS `tags` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name`       VARCHAR(64)     NOT NULL,
  `slug`       VARCHAR(64)     NOT NULL,
  `created_at` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3)     NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_name` (`name`),
  UNIQUE KEY `uniq_slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '标签';

-- ============== articles ==============
CREATE TABLE IF NOT EXISTS `articles` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `title`        VARCHAR(200)    NOT NULL,
  `slug`         VARCHAR(200)    NOT NULL,
  `summary`      VARCHAR(500)    NOT NULL DEFAULT '',
  `content`      MEDIUMTEXT      NOT NULL,
  `cover_image`  VARCHAR(255)    NOT NULL DEFAULT '',
  `category_id`  BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `author_id`    BIGINT UNSIGNED NOT NULL,
  `status`       VARCHAR(16)     NOT NULL DEFAULT 'draft',
  `view_count`   INT             NOT NULL DEFAULT 0,
  `published_at` DATETIME(3)     NULL DEFAULT NULL,
  `created_at`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`   DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at`   DATETIME(3)     NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_slug` (`slug`),
  KEY `idx_status_published_at` (`status`, `published_at`),
  KEY `idx_category_id`         (`category_id`),
  KEY `idx_author_id`           (`author_id`),
  KEY `idx_deleted_at`          (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '文章';

-- ============== article_tags ==============
CREATE TABLE IF NOT EXISTS `article_tags` (
  `article_id` BIGINT UNSIGNED NOT NULL,
  `tag_id`     BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`article_id`, `tag_id`),
  KEY `idx_tag_id` (`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '文章-标签';

-- ============== profiles ==============
CREATE TABLE IF NOT EXISTS `profiles` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id`     BIGINT UNSIGNED NOT NULL,
  `avatar`      VARCHAR(255)    NOT NULL DEFAULT '',
  `nickname`    VARCHAR(64)     NOT NULL DEFAULT '',
  `title`       VARCHAR(128)    NOT NULL DEFAULT '',
  `bio`         TEXT,
  `skills`      JSON,
  `experiences` JSON,
  `contacts`    JSON,
  `created_at`  DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`  DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at`  DATETIME(3)     NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '关于我';

SET FOREIGN_KEY_CHECKS = 1;
