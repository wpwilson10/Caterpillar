CREATE TABLE Listing(
	listing_id serial PRIMARY KEY,
	update_time timestamptz, -- when this current contact became valid
	is_enabled boolean, -- IEX enabled
	is_active boolean DEFAULT TRUE, -- whether the stocks app is using this listing. Default to true to ensure new listings are extracted cause IPOs are important.
	symbol text,
	name text,
	iex_id text, -- unique ID applied by IEX to track securities through symbol changes.
	type text,
	region text,
	currency text,
	exchange text,
	is_sp500 boolean DEFAULT False, -- true if the listing is currently on the index, false otherwise.
	is_russell3000 boolean DEFAULT False -- true if the listing is currently on the index, false otherwise. 
);

CREATE TABLE AuditListing(
	audit_id serial PRIMARY KEY,  -- this audit contacts
	listing_id int REFERENCES Listing, -- reference to current listing record
	update_time timestamptz, -- when this current contact became valid
	is_enabled boolean,
	symbol text,
	name text,
	iex_id text, -- unique ID applied by IEX to track securities through symbol changes.
	type text,
	region text,
	currency text,
	exchange text,
	is_sp500 boolean,
	is_russell3000 boolean
);

CREATE TABLE DataSource(
	source_id smallint PRIMARY KEY,
	source_name text
);

CREATE TABLE Intraday(
	intraday_id bigserial PRIMARY KEY,
	listing_id int REFERENCES Listing,
	source_id smallint REFERENCES DataSource,
	data_time timestamptz,
	update_time timestamptz,
	open numeric,
	close numeric,
	high numeric,
	low numeric,
	volume numeric,
	notional numeric,
	num_trades numeric
);

-- Set id to start at 1000 instead of 1
ALTER SEQUENCE Listing_listing_id_seq RESTART WITH 1000;
ALTER SEQUENCE intraday_intraday_id_seq RESTART WITH 1000;
-- Useful index
CREATE INDEX index_intra_time ON Intraday (listing_id, data_time);

-- Add values to data source 
INSERT INTO DataSource (source_id, source_name) VALUES (1, 'IEX'), (2, 'Alpha Vantage');