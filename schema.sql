CREATE TABLE `url` (
  `id` varchar(30) NOT NULL,
  `url` longtext NOT NULL,
  `expiry` TIMESTAMP,
  PRIMARY KEY (`id`)
);
