CREATE TABLE NewsArticle(
	article_id bigserial PRIMARY KEY,
	data_entry_time timestamptz, -- when this data was collected and inserted
	source text, -- source of the article/where we got the link
	host text, -- hostname parsed from url
	link text, -- originally queried link url
	source_published_time timestamptz, -- time published from reference source
	published_time timestamptz, -- time published from newspaper3k
	source_title text, -- title from reference source
	title text, -- title from newspaper3k
	canonical_link text, -- canonical link if one exists
	body text, -- main article text
	authors text -- json array of authors
);

-- Set id to start at 1000 instead of 1
ALTER SEQUENCE NewsArticle_article_id_seq RESTART WITH 1000;
-- Indecies
CREATE INDEX article_source_published_time_index ON newsarticle(source_published_time NULLS LAST);
CREATE INDEX article_source_published_time_desc_index ON newsarticle(source_published_time DESC NULLS LAST);
