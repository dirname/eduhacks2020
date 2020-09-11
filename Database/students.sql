-- ----------------------------
-- Table structure for ban
-- ----------------------------
-- CREATE TABLE IF NOT EXISTS "student"."ban"
-- (
--     "id"         bigserial,
--     "ban_id"     text NOT NULL UNIQUE,
--     "message"    text,
--     "created_at" timestamptz,
--     "updated_at" timestamptz,
--     "deleted_at" timestamptz,
--     PRIMARY KEY ("id")
-- );
-- CREATE INDEX IF NOT EXISTS "idx_bans_deleted_at" ON "student"."ban" ("deleted_at");
-- ----------------------------
-- Table structure for users
-- ----------------------------
CREATE TABLE IF NOT EXISTS "student"."users"
(
    "id"         bigserial,
    "user_id"    text NOT NULL UNIQUE,
    "username"   text NOT NULL UNIQUE,
    "password"   text NOT NULL,
    "nickname"   text NOT NULL,
    "gender"     boolean,
    "phone"      text UNIQUE,
    "email"      text UNIQUE,
    "avatar"     text,
    "birthday"   timestamptz,
    "banned"     boolean DEFAULT false,
    "class_id"   bigint,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_students_class" FOREIGN KEY ("class_id") REFERENCES "college"."classes" ("id")
);
CREATE INDEX IF NOT EXISTS "idx_students_deleted_at" ON "student"."users" ("deleted_at");
-- INSERT INTO "student".ban ("id", "ban_id", "message", "created_at", "updated_at", "deleted_at")
-- VALUES (1, '00000000-0000-0000-0000-000000000000', 'Normal', '2020-09-09 17:05:16.446', '2020-09-09 17:05:16.446', NULL)
-- RETURNING "id"; -- This is a default record which make column of 'banned_id' valid.