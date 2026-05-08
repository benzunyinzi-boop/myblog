-- 000002_profile_schema_rewrite.down.sql
-- 回滚到 M2 版本的 profiles 结构

SET NAMES utf8mb4;

DELETE FROM `profiles`;

ALTER TABLE `profiles`
  DROP COLUMN `name`,
  DROP COLUMN `email`,
  DROP COLUMN `github`,
  DROP COLUMN `twitter`,
  DROP COLUMN `linkedin`,
  DROP COLUMN `website`;

ALTER TABLE `profiles`
  MODIFY `bio` TEXT NULL,
  ADD COLUMN `user_id`     BIGINT UNSIGNED NOT NULL AFTER `id`,
  ADD COLUMN `nickname`    VARCHAR(64)  NOT NULL DEFAULT ''  AFTER `avatar`,
  ADD COLUMN `title`       VARCHAR(128) NOT NULL DEFAULT ''  AFTER `nickname`,
  ADD COLUMN `skills`      JSON                              AFTER `bio`,
  ADD COLUMN `experiences` JSON                              AFTER `skills`,
  ADD COLUMN `contacts`    JSON                              AFTER `experiences`,
  ADD UNIQUE KEY `uniq_user_id` (`user_id`);
