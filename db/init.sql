CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
    title TEXT,
    amount FLOAT,
    note TEXT,
    tags TEXT []
);

INSERT INTO TABLE expenses (title, amount, note, tags) VALUES ('Rent', 1000, 'Rent for the month of May', ARRAY['rent', 'housing']);