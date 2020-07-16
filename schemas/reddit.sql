-- fields described at https://github.com/reddit-archive/reddit/wiki/JSON
CREATE TABLE RedditSubmission(
	submission_id bigserial PRIMARY KEY,
	reddit_id text, -- reddit's unique ID for this submission
	title text,
	url text, -- link to what this post is about
	permalink text, -- reddits relative url for this submission
	data_entry_time timestamptz,
	created_time timestamptz,
	user_name text,
	subreddit_name text,
	subreddit_id text,
	selftext text,
	selftext_html text,
	num_comments int,
	score int,
	up_votes int,
	down_votes int,
	is_nsfw boolean,
	is_self boolean
);

-- Set id to start at 1000 instead of 1
ALTER SEQUENCE redditsubmission_submission_id_seq RESTART WITH 1000;

CREATE TABLE RedditComment(
	comment_id bigserial PRIMARY KEY,
	submission_id bigint REFERENCES RedditSubmission(submission_id),
	reddit_id text, -- reddit's unique ID for this comment
	parent_id text, -- ID of the thing this comment is a reply to. Uses reddit's unique ID
	data_entry_time timestamptz,
	created_time timestamptz,
	user_name text,
	body text,
	body_html text,
	up_votes int,
	down_votes int,
	is_deleted boolean
);


-- Set id to start at 1000 instead of 1
ALTER SEQUENCE redditcomment_comment_id_seq RESTART WITH 1000;

CREATE TABLE RedditNews(
	CONSTRAINT link_id PRIMARY KEY(article_id, submission_id),
	article_id bigint REFERENCES NewsArticle(article_id),
	submission_id bigint REFERENCES RedditSubmission(submission_id),
	data_entry_time timestamptz
);
