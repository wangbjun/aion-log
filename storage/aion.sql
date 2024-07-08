CREATE TABLE `aion_player_battle_log` (
                                          `id` int NOT NULL AUTO_INCREMENT,
                                          `player` varchar(100) DEFAULT NULL,
                                          `skill` varchar(100) DEFAULT NULL,
                                          `target` varchar(100) DEFAULT NULL,
                                          `value` int DEFAULT NULL,
                                          `time` datetime DEFAULT NULL,
                                          `raw_msg` varchar(255) DEFAULT NULL,
                                          PRIMARY KEY (`id`),
                                          KEY `index_time` (`time`),
                                          KEY `index_player` (`player`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `aion_player_info` (
                                    `id` int NOT NULL AUTO_INCREMENT,
                                    `name` varchar(100) DEFAULT NULL,
                                    `type` int DEFAULT NULL,
                                    `class` int DEFAULT NULL,
                                    `time` datetime DEFAULT NULL,
                                    PRIMARY KEY (`id`),
                                    UNIQUE KEY `index_uniq` (`name`,`type`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `aion_player_rank` (
                                    `id` int NOT NULL AUTO_INCREMENT,
                                    `player` varchar(100) DEFAULT NULL,
                                    `count` int DEFAULT NULL,
                                    `time` datetime DEFAULT NULL,
                                    PRIMARY KEY (`id`),
                                    UNIQUE KEY `index_uniq` (`player`,`count`,`time`)
) ENGINE=InnoDB AUTO_INCREMENT=4093 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `aion_player_skill` (
                                     `skill` varchar(100) NOT NULL,
                                     `class` int NOT NULL,
                                     UNIQUE KEY `aion_player_skill_pk` (`class`,`skill`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `users` (
                         `id` int NOT NULL AUTO_INCREMENT,
                         `name` varchar(100) DEFAULT NULL,
                         `email` varchar(100) DEFAULT NULL,
                         `password` varchar(100) DEFAULT NULL,
                         `salt` varchar(100) DEFAULT NULL,
                         `created_at` datetime DEFAULT NULL,
                         `updated_at` datetime DEFAULT NULL,
                         `deleted_at` datetime DEFAULT NULL,
                         PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
