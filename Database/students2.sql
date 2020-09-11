/*
 Navicat Premium Data Transfer

 Source Server         : PostgreSQL
 Source Server Type    : PostgreSQL
 Source Server Version : 120004
 Source Host           : 139.9.240.38:5432
 Source Catalog        : education
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 120004
 File Encoding         : 65001

 Date: 09/09/2020 17:01:54
*/


-- ----------------------------
-- Table structure for students
-- ----------------------------
DROP TABLE IF EXISTS "public"."students";
CREATE TABLE "public"."students" (
  "id" int8 NOT NULL DEFAULT nextval('students_id_seq'::regclass),
  "username" text COLLATE "pg_catalog"."default" NOT NULL,
  "password" text COLLATE "pg_catalog"."default" NOT NULL,
  "nickname" text COLLATE "pg_catalog"."default" NOT NULL,
  "gender" bool,
  "phone" text COLLATE "pg_catalog"."default",
  "email" text COLLATE "pg_catalog"."default",
  "avatar" text COLLATE "pg_catalog"."default",
  "birthday" timestamptz(6),
  "banned" bool DEFAULT false,
  "banned_id" int8,
  "class_id" int8,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6)
)
;

-- ----------------------------
-- Indexes structure for table students
-- ----------------------------
CREATE INDEX "idx_students_deleted_at" ON "public"."students" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table students
-- ----------------------------
ALTER TABLE "public"."students" ADD CONSTRAINT "students_username_key" UNIQUE ("username");
ALTER TABLE "public"."students" ADD CONSTRAINT "students_phone_key" UNIQUE ("phone");
ALTER TABLE "public"."students" ADD CONSTRAINT "students_email_key" UNIQUE ("email");

-- ----------------------------
-- Primary Key structure for table students
-- ----------------------------
ALTER TABLE "public"."students" ADD CONSTRAINT "students_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table students
-- ----------------------------
ALTER TABLE "public"."students" ADD CONSTRAINT "fk_students_ban" FOREIGN KEY ("banned_id") REFERENCES "public"."bans" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."students" ADD CONSTRAINT "fk_students_class" FOREIGN KEY ("class_id") REFERENCES "public"."classes" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
