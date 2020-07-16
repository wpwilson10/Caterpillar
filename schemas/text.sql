CREATE TABLE ProcessedText(
	text_id bigserial PRIMARY KEY,
	data_entry_time timestamptz, -- when this data was collected and inserted
	title text, -- title if it exists
	body text, -- main text after preprocessing
	
);

-- Set id to start at 1000 instead of 1
ALTER SEQUENCE ProcessedText_text_id_seq RESTART WITH 1000;