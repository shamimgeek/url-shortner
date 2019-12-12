CREATE DATABASE tinyurl;
USE tinyurl;

DROP TABLE IF EXISTS `tiny_urls`;

CREATE TABLE `tiny_urls` (
	`shortid` varchar(14) collate utf8mb4_unicode_ci NOT NULL,
	`url` varchar(620) collate utf8mb4_unicode_ci NOT NULL,
	`date` datetime NOT NULL,
	`hits` bigint(20) NOT NULL default '0',
	PRIMARY KEY (`shortid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Used for the URL shortener';

INSERT INTO `redirect` VALUES ('xd5k6d', 'https://gitlabe2.ext.net.nokia.com/dyminski/shortener', NOW(), 1);
