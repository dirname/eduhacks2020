-- ----------------------------
-- Table structure for colleges
-- ----------------------------
CREATE TABLE IF NOT EXISTS "college"."colleges"
(
    "id"           bigserial,
    "college_id"   text NOT NULL UNIQUE,
    "college_name" text NOT NULL UNIQUE,
    "created_at"   timestamptz,
    "updated_at"   timestamptz,
    "deleted_at"   timestamptz,
    PRIMARY KEY ("id")
);

CREATE INDEX IF NOT EXISTS "idx_colleges_deleted_at" ON "college"."colleges" ("deleted_at");
-- ----------------------------
-- Table structure for majors
-- ----------------------------
CREATE TABLE IF NOT EXISTS "college"."majors"
(
    "id"         bigserial,
    "major_id"   text NOT NULL UNIQUE,
    "major_name" text NOT NULL UNIQUE,
    "college_id" bigint,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_majors_college" FOREIGN KEY ("college_id") REFERENCES "college"."colleges" ("id")
);
CREATE INDEX IF NOT EXISTS "idx_majors_deleted_at" ON "college"."majors" ("deleted_at");
-- ----------------------------
-- Table structure for classes
-- ----------------------------
CREATE TABLE IF NOT EXISTS "college"."classes"
(
    "id"         bigserial,
    "class_id"   text NOT NULL UNIQUE,
    "class_name" text NOT NULL UNIQUE,
    "major_id"   bigint,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_classes_major" FOREIGN KEY ("major_id") REFERENCES "college"."majors" ("id")
);
CREATE INDEX IF NOT EXISTS "idx_classes_deleted_at" ON "college"."classes" ("deleted_at");