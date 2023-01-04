-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS news_articles_id_seq;

-- Table Definition
CREATE TABLE "news_articles" (
    "id" int4 NOT NULL DEFAULT nextval('news_articles_id_seq'::regclass),
    "title" text,
    "content" text,
    "author" text,
    PRIMARY KEY ("id")
);

INSERT INTO "news_articles" ("id", "title", "content", "author") VALUES (1, 'test-title', 'test-content', 'test-author');