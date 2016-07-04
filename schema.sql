CREATE TABLE `pastebin` (
  `id` varchar(30) NOT NULL,
  `url` char(256) default NULL,
  `expiry` TIMESTAMP,
  PRIMARY KEY (`id`)
);
