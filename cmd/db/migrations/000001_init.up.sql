CREATE TABLE IF NOT EXISTS users (
                                     id INTEGER PRIMARY KEY,
                                     username VARCHAR(255) NOT NULL,
                                     password VARCHAR(255) NOT NULL,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     taskd_uuid VARCHAR(255),
                                     timew_id INT,
                                     timew_hook BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS sessions (
                                        id INTEGER PRIMARY KEY,
                                        user_id INT NOT NULL,
                                        token VARCHAR(255) NOT NULL,
                                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);