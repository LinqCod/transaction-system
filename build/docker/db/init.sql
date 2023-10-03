CREATE TABLE IF NOT EXISTS accounts (
                                        id SERIAL PRIMARY KEY,
                                        card_number VARCHAR,
                                        balance DECIMAL
);

INSERT INTO accounts (card_number, balance) VALUES ('1234-4123-4123-4123', 0);
INSERT INTO accounts (card_number, balance) VALUES ('3456-2367-2341-7842', 1244);
INSERT INTO accounts (card_number, balance) VALUES ('4556-2267-4444-2353', 12222);
INSERT INTO accounts (card_number, balance) VALUES ('2131-2356-2356-2322', 111111);
