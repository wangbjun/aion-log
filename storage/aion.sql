CREATE TABLE `aion_chat_log`
(
    `id`      int NOT NULL AUTO_INCREMENT,
    `player`  varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `skill`   varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `target`  varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `value`   int                                                           DEFAULT NULL,
    `time`    datetime                                                      DEFAULT NULL,
    `raw_msg` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_player_skill_time` (`player`, `skill`, `time`),
    KEY `idx_skill` (`skill`),
    KEY `idx_time` (`time`),
    KEY `idx_target` (`target`),
    KEY `idx_value` (`value`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

CREATE TABLE `aion_player_info`
(
    `id`    int NOT NULL AUTO_INCREMENT,
    `name`  varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `type`  int                                                           DEFAULT NULL,
    `class` int                                                           DEFAULT NULL,
    `time`  datetime                                                      DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `index_uniq` (`name`, `type`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 6458
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

CREATE TABLE `aion_player_rank`
(
    `id`     int NOT NULL AUTO_INCREMENT,
    `player` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
    `count`  int                                                           DEFAULT NULL,
    `time`   datetime                                                      DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `index_uniq` (`player`, `count`, `time`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 6002
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

CREATE TABLE `aion_player_skill`
(
    `skill` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
    `class` int                                                           NOT NULL,
    UNIQUE KEY `aion_player_skill_pk` (`class`, `skill`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

CREATE TABLE `aion_timeline`
(
    `time`  datetime NOT NULL,
    `value` int      NOT NULL,
    `type`  int      NOT NULL DEFAULT '0'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci
