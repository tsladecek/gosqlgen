CREATE TABLE IF NOT EXISTS users (
    _id INTEGER PRIMARY KEY AUTOINCREMENT,
    id TEXT NOT NULL UNIQUE,
    name BLOB NOT NULL,
    payload TEXT NOT NULL,
    age INTEGER,
    drives_car INTEGER,
    birthday TEXT,
    registered TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS admins (
    _id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    FOREIGN KEY (_id) REFERENCES users(_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS countries (
    _id INTEGER PRIMARY KEY AUTOINCREMENT,
    id TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    gps TEXT NOT NULL,
    continent TEXT NOT NULL CHECK (continent IN ('Asia', 'Europe', 'Africa'))
);

CREATE TABLE IF NOT EXISTS addresses (
    _id INTEGER PRIMARY KEY AUTOINCREMENT,
    id TEXT NOT NULL UNIQUE,
    address TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    country_id INTEGER NOT NULL,
    deleted_at TEXT,
    ipv4 TEXT NOT NULL,
    ipv6 TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(_id),
    FOREIGN KEY (country_id) REFERENCES countries(_id)
);

CREATE TABLE IF NOT EXISTS addresses_book (
    _id INTEGER PRIMARY KEY AUTOINCREMENT,
    id TEXT NOT NULL UNIQUE,
    address_id INTEGER NOT NULL,
    FOREIGN KEY (address_id) REFERENCES addresses(_id) ON DELETE CASCADE
);

