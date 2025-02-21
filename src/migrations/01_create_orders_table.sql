CREATE TABLE IF NOT EXISTS orders (
    id INT auto_increment PRIMARY KEY,
    user_id INT NOT NULL,
    item_id INT NOT NULL,
    price FLOAT NOT NULL,
    created BIGINT UNSIGNED,
    updated BIGINT UNSIGNED,
    completed TINYINT DEFAULT 0