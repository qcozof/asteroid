PRAGMA foreign_keys = false;

-- ----------------------------
-- Table structure for file_list
-- ----------------------------
DROP TABLE IF EXISTS "file_list";
CREATE TABLE "file_list" (
  "file_id" integer null,--NOT NULL PRIMARY KEY AUTOINCREMENT,
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
DROP TABLE IF EXISTS  "sqlite_sequence";

PRAGMA foreign_keys = true;
