CREATE TABLE IF NOT EXISTS mangas (
  slug TEXT PRIMARY KEY NOT NULL,
  url  TEXT NOT NULL,

  name TEXT NOT NULL,
  description TEXT NOT NULL,

  image TEXT NOT NULL, -- BASE64 (expensive)
  image_type TEXT NOT NULL,

  status TEXT DEFAULT NULL,
  time TEXT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS manga_tags (
  tag TEXT NOT NULL,
  manga_slug TEXT NOT NULL,

  CONSTRAINT manga_tags_pk PRIMARY KEY (manga_slug, tag),
  CONSTRAINT manga_fk FOREIGN KEY (manga_slug) REFERENCES mangas (slug)
);

CREATE TABLE IF NOT EXISTS chapters (
  id INTEGER PRIMARY KEY,
  manga_slug TEXT NOT NULL,

  title TEXT NOT NULL,
  "link" TEXT NOT NULL,

  time TEXT NOT NULL,

  CONSTRAINT manga_fk FOREIGN KEY (manga_slug) REFERENCES mangas (slug)
);

CREATE TABLE IF NOT EXISTS pages (
  "order" INTEGER NOT NULL,
  chapter_id INTEGER NOT NULL,

  image TEXT NOT NULL,
  image_type TEXT NOT NULL,

  CONSTRAINT chapter_order_pk PRIMARY KEY (chapter_id, "order"),
  CONSTRAINT chapter_fk FOREIGN KEY (chapter_id) REFERENCES chapters (id)
);
