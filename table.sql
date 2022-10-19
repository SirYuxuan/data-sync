/*
 Navicat Premium Data Transfer

 Source Server         : DataSync
 Source Server Type    : SQLite
 Source Server Version : 3030001
 Source Schema         : main

 Target Server Type    : SQLite
 Target Server Version : 3030001
 File Encoding         : 65001

 Date: 19/10/2022 17:22:30
*/

PRAGMA foreign_keys = false;

-- ----------------------------
-- Table structure for config_table
-- ----------------------------
DROP TABLE IF EXISTS "config_table";
CREATE TABLE "config_table" (
  "id" INTEGER NOT NULL,
  "name" TEXT,
  "source_id" INTEGER,
  "target_id" INTEGER,
  "source_table" TEXT,
  "target_table" TEXT,
  "interval" integer,
  "last_time" text,
  "where_sql" TEXT,
  PRIMARY KEY ("id")
);

-- ----------------------------
-- Table structure for config_table_field
-- ----------------------------
DROP TABLE IF EXISTS "config_table_field";
CREATE TABLE "config_table_field" (
  "id" INTEGER NOT NULL,
  "config_id" INTEGER,
  "source_field" TEXT,
  "target_field" TEXT,
  PRIMARY KEY ("id")
);

-- ----------------------------
-- Table structure for db_list
-- ----------------------------
DROP TABLE IF EXISTS "db_list";
CREATE TABLE "db_list" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "name" TEXT,
  "db_type" integer,
  "ip" TEXT,
  "port" integer,
  "username" TEXT,
  "password" TEXT,
  "db_name" TEXT
);

-- ----------------------------
-- Auto increment value for db_list
-- ----------------------------
UPDATE "sqlite_sequence" SET seq = 2 WHERE name = 'db_list';

PRAGMA foreign_keys = true;
