/*
 Navicat Premium Data Transfer

 Source Server         : localdb1
 Source Server Type    : SQLite
 Source Server Version : 3030001
 Source Schema         : main

 Target Server Type    : SQLite
 Target Server Version : 3030001
 File Encoding         : 65001

 Date: 11/04/2023 10:38:05
*/

PRAGMA foreign_keys = false;

-- ----------------------------
-- Table structure for file_list
-- ----------------------------
DROP TABLE IF EXISTS "file_list";
CREATE TABLE "file_list" (
  "file_id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "site" text NOT NULL,
  "file_dir" text NOT NULL,
  "file_name" text NOT NULL,
  "hash" text NOT NULL,
  "perm" integer NOT NULL,
  "policy" text NOT NULL,
  "create_time" integer NOT NULL,
  "update_time" integer NOT NULL
);

-- ----------------------------
-- Auto increment value for file_list
-- ----------------------------
UPDATE "sqlite_sequence" SET seq = 0 WHERE name = 'file_list';

PRAGMA foreign_keys = true;
