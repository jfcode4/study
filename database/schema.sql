CREATE TABLE IF NOT EXISTS "Decks" (
	"id"         INTEGER NOT NULL DEFAULT 1,
	"name"       TEXT NOT NULL,
	"date"  TEXT NOT NULL,
	"day"        INTEGER NOT NULL DEFAULT 1,
	"owner"      INTEGER NOT NULL,
	PRIMARY KEY("id"),
	FOREIGN KEY("owner") REFERENCES "Users"("id")
);

CREATE TABLE IF NOT EXISTS "Users" (
	"id"         INTEGER NOT NULL,
	"username"   INTEGER NOT NULL UNIQUE,
	"hash"       TEXT NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
);

CREATE TABLE IF NOT EXISTS "Cards" (
	"id"         INTEGER NOT NULL,
	"deck_id"    INTEGER NOT NULL,
	"question"   TEXT NOT NULL,
	"answer"     TEXT NOT NULL,
	"due"        TEXT NOT NULL DEFAULT '',
	"interval"   INTEGER NOT NULL DEFAULT 1,
	"tags"       TEXT NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
);
