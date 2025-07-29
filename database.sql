-- Create a new database, e.g., 'finance_tracker'
-- CREATE DATABASE finance_tracker;
-- Connect to your new database and run the following:
-- Create the 'categories' table to store expense/income categories
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    parent_id INTEGER REFERENCES categories(id) ON DELETE
    SET NULL,
        user_id INTEGER NO NULL REFERENCES users(id) ON DELETE CASCADE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- Creates the 'transactions' table to store financial records
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    description TEXT,
    amount NUMERIC(10, 2) NOT NULL,
    type VARCHAR(7) NOT NULL CHECK (type IN ('income', 'expense')),
    date TIMESTAMPTZ NOT NULL,
    category_id INTEGER REFERENCES categories(id) ON DELETE
    SET NULL,
        user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- Creates the 'users' table to store user authentication information
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- Seed some initial categories
INSERT INTO categories (name, user_id)
VALUES ('Groceries', 1),
    ('Salary', 1),
    ('Rent', 1),
    ('Utilities', 1),
    ('Entertainment', 1);