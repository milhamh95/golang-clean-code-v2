CREATE TABLE IF NOT EXISTS `departments` (
    `id` varchar (50) NOT NULL,
    `name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
    `description` varchar(250) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
    `created_time` timestamp NULL,
    `updated_time` timestamp NULL,
    PRIMARY KEY (`id`),
    FULLTEXT KEY `name_idx` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
