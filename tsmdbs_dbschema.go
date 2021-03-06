package tsmdbs

const DBSQL=`
PRAGMA foreign_keys = ON;
DROP TABLE IF EXISTS TS;
DROP TABLE IF EXISTS HOST;
DROP TABLE IF EXISTS KEY;
DROP TABLE IF EXISTS RELATION;
DROP TABLE IF EXISTS DATA;
DROP TABLE IF EXISTS DREL;
DROP INDEX IF EXISTS TS1;
DROP INDEX IF EXISTS HOST1;
DROP INDEX IF EXISTS KEY1;
DROP INDEX IF EXISTS RELATION1;
CREATE TABLE TS (
  ID    INTEGER PRIMARY KEY AUTOINCREMENT,
  START INTEGER,
  END   INTEGER
);
CREATE UNIQUE INDEX TS1 ON TS (START, END);
CREATE TABLE HOST (
  ID    INTEGER PRIMARY KEY AUTOINCREMENT,
  NAME  TEXT
);
CREATE UNIQUE INDEX HOST1 ON HOST (NAME);
CREATE TABLE KEY (
  ID    INTEGER PRIMARY KEY AUTOINCREMENT,
  NAME  TEXT
);
CREATE UNIQUE INDEX KEY1 ON KEY (NAME);
CREATE TABLE RELATION (
  ID    INTEGER PRIMARY KEY AUTOINCREMENT,
  NAME  TEXT
);
CREATE UNIQUE INDEX RELATION1 ON RELATION(NAME);
CREATE TABLE DATA (
  ID      TEXT PRIMARY KEY,
  TS      INTEGER,
  HOST    INTEGER,
  KEY     INTEGER,
  VALUE   BLOB NOT NULL,
  FOREIGN KEY(TS) REFERENCES TS(ID),
  FOREIGN KEY(HOST) REFERENCES HOST(ID),
  FOREIGN KEY(KEY) REFERENCES KEY(ID)
);
CREATE TABLE DREL (
  REL     INTEGER,
  DATA    TEXT,
  FOREIGN KEY(REL) REFERENCES RELATION(ID),
  FOREIGN KEY(DATA) REFERENCES DATA(ID)
);
`
