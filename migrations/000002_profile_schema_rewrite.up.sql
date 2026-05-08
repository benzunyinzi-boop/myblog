-- 000002_profile_schema_rewrite.up.sql
-- 调整 profiles 表结构,匹配 M3 版本 model:
--   去掉 user_id/nickname/title/skills/experiences/contacts
--   改为单条记录的扁平结构:name/bio/avatar/email/github/twitter/linkedin/website

SET NAMES utf8mb4;

-- 先清空老数据(M1/M2 没有实际写入 profiles,安全)
DELETE FROM `profiles`;

-- 删除老列
ALTER TABLE `profiles`
  DROP INDEX `uniq_user_id`,
  DROP COLUMN `user_id`,
  DROP COLUMN `nickname`,
  DROP COLUMN `title`,
  DROP COLUMN `skills`,
  DROP COLUMN `experiences`,
  DROP COLUMN `contacts`;

-- 新增列
ALTER TABLE `profiles`
  ADD COLUMN `name`     VARCHAR(100) NOT NULL DEFAULT ''  AFTER `id`,
  ADD COLUMN `email`    VARCHAR(100) NOT NULL DEFAULT ''  AFTER `avatar`,
  ADD COLUMN `github`   VARCHAR(255) NOT NULL DEFAULT ''  AFTER `email`,
  ADD COLUMN `twitter`  VARCHAR(255) NOT NULL DEFAULT ''  AFTER `github`,
  ADD COLUMN `linkedin` VARCHAR(255) NOT NULL DEFAULT ''  AFTER `twitter`,
  ADD COLUMN `website`  VARCHAR(255) NOT NULL DEFAULT ''  AFTER `linkedin`;

-- bio 从 TEXT 保留,但去掉 nullable
ALTER TABLE `profiles`
  MODIFY `bio` TEXT NOT NULL;
