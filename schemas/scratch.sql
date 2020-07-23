GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO dbuser;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO dbuser;

ALTER TABLE "listing" ADD COLUMN "is_active" boolean DEFAULT False;
ALTER TABLE "auditlisting" ADD COLUMN "is_active" boolean DEFAULT False;

ALTER TABLE RedditQueue 
RENAME full_id TO submission_id;

Copy RedditSubmission to '/home/patrick/Documents/Projects/Caterpillar/test/sub.csv' DELIMITER ',' CSV HEADER;

Copy (select body from newsarticle where host='www.foxnews.com' limit 10) To '/home/patrick/Documents/Projects/Caterpillar/test/foxnews.csv' With CSV DELIMITER ',' HEADER;