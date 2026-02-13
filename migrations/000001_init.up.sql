-- MySQL database

CREATE TABLE IF NOT EXISTS telegram_users (
    telegram_id BIGINT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    locale VARCHAR(10) DEFAULT 'en',
    preferences JSON NOT NULL,
    is_bot_blocked BOOLEAN DEFAULT FALSE,
    technical_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS feedback (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    telegram_id BIGINT,
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (telegram_id) REFERENCES telegram_users(telegram_id)
);

CREATE TABLE IF NOT EXISTS admins (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    telegram_id BIGINT,
    discord_id VARCHAR(255),

    FOREIGN KEY (telegram_id) REFERENCES telegram_users(telegram_id)
);

CREATE TABLE IF NOT EXISTS admin_otp (
    id INT PRIMARY KEY AUTO_INCREMENT,
    admin_id INT NOT NULL,
    code VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (admin_id) REFERENCES admins(id)
);

CREATE TABLE IF NOT EXISTS telegram_interactions (
    id INT PRIMARY KEY AUTO_INCREMENT,
    telegram_user_id BIGINT NOT NULL,
    telegram_chat_id BIGINT NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_telegram_chat (telegram_user_id, telegram_chat_id),
    INDEX idx_telegram_created (telegram_user_id, created_at),

    FOREIGN KEY (telegram_user_id) REFERENCES telegram_users(telegram_id)
);