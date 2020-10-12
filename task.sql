
BEGIN;

DROP TABLE IF EXISTS "gotask"."task";
CREATE TABLE "gotask"."task" (
  "name" TEXT NOT NULL PRIMARY KEY,
  "type" TEXT NOT NULL,
  "host" TEXT NOT NULL,
);

COMMIT;