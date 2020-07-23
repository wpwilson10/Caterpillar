-- not reallly a script, but here are the commands we care about

DROP DATABASE tst_stocks;
CREATE DATABASE tst_stocks;

pg_dump postgresql://postgres:DuncanCollege@localhost/prd_stocks --schema-only | psql postgresql://postgres:DuncanCollege@localhost/tst_stocks
pg_dump postgresql://postgres:DuncanCollege@localhost/prd_stocks --table=DataSource | psql postgresql://postgres:DuncanCollege@localhost/tst_stocks
pg_dump postgresql://postgres:DuncanCollege@localhost/prd_stocks --table=Listing | psql postgresql://postgres:DuncanCollege@localhost/tst_stocks

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO dbuser;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO dbuser;