CREATE TABLE IF NOT EXISTS `employees` (
    `id` varchar(50) NOT NULL,
    `first_name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
    `last_name` varchar(200) COLLATE utf8mb4_unicode_ci NULL DEFAULT '',
    `birth_place` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
    `date_of_birth` DATE NOT NULL,
    `title` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
    `dept_id` varchar(50) NOT NULL,
    PRIMARY KEY (`id`),
    FULLTEXT KEY `first_name_idx` (`first_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
