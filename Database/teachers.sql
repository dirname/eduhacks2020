-- ----------------------------
-- Table structure for Teacher
-- ----------------------------
CREATE TABLE "teacher"."users"
(
    "id"         bigserial,
    "user_id"    uuid NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
    "username"   text NOT NULL UNIQUE,
    "password"   text NOT NULL,
    "nickname"   text NOT NULL,
    "gender"     boolean,
    "phone"      text UNIQUE,
    "email"      text UNIQUE,
    "avatar"     text,
    "birthday"   timestamptz,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    PRIMARY KEY ("id")
);
CREATE INDEX "idx_teachers_deleted_at" ON "teacher"."users" ("deleted_at");
-- ----------------------------
-- Table structure for Course
-- ----------------------------
CREATE TABLE "resources"."courses"
(
    "id"          bigserial,
    "course_id"   uuid NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
    "course_name" text NOT NULL,
    "major_id"    bigint,
    "teacher_id"  bigint,
    "status"      boolean,
    "open"        boolean,
    "start_at"    timestamptz,
    "end_at"      timestamptz,
    "created_at"  timestamptz,
    "updated_at"  timestamptz,
    "deleted_at"  timestamptz,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_courses_major" FOREIGN KEY ("major_id") REFERENCES "college"."majors" ("id"),
    CONSTRAINT "fk_courses_teacher" FOREIGN KEY ("teacher_id") REFERENCES "teacher"."users" ("id")
);
CREATE INDEX "idx_courses_deleted_at" ON "resources"."courses" ("deleted_at")



