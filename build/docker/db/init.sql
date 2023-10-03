DROP TABLE accounts;

CREATE TABLE IF NOT EXISTS accounts (
                                        id SERIAL PRIMARY KEY,
                                        card_number VARCHAR,
                                        balance DECIMAL
);
