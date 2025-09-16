CREATE TABLE IF NOT EXISTS `users` (
    `_id` INT(11) NOT NULL AUTO_INCREMENT,
    `id` VARCHAR(255) NOT NULL UNIQUE,
    `name` BLOB NOT NULL,
    `payload` JSON NOT NULL,
    `age` INT(11),
    `drives_car` TINYINT(1),
    `birthday` DATETIME,
    `registered` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS `admins` (
    `_id` INT(11) NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(31) NOT NULL,
    PRIMARY KEY (`_id`),
    FOREIGN KEY (`_id`) REFERENCES `users`(`_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS `countries` (
    `_id` INT(11) NOT NULL AUTO_INCREMENT,
    `id` VARCHAR(255) NOT NULL UNIQUE,
    `name` VARCHAR(255) NOT NULL,
    `gps` VARCHAR(255) NOT NULL,
    `continent` ENUM('Asia', 'Europe', 'Africa') NOT NULL,
    PRIMARY KEY (`_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS `addresses` (
    `_id` INT(11) NOT NULL AUTO_INCREMENT,
    `id` VARCHAR(255) NOT NULL UNIQUE,
    `address` VARCHAR(255) NOT NULL,
    `user_id` INT(11) NOT NULL,
    `country_id` INT(11) NOT NULL,
    `deleted_at` DATETIME,
    PRIMARY KEY (`_id`),
    FOREIGN KEY (`user_id`) REFERENCES `users`(`_id`),
    FOREIGN KEY (`country_id`) REFERENCES `countries`(`_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS `addresses_book` (
    `_id` INT(11) NOT NULL AUTO_INCREMENT,
    `id` VARCHAR(255) NOT NULL UNIQUE,
    `address_id` INT(11) NOT NULL,
    PRIMARY KEY (`_id`),
    FOREIGN KEY (`address_id`) REFERENCES `addresses`(`_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

