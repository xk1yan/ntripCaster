CREATE TABLE `mountpoint` (
`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
`name` varchar(30) NOT NULL,
`password` varchar(30) NOT NULL DEFAULT '123456',
`desc` varchar(3000) DEFAULT NULL,
`status` enum('y','n') DEFAULT 'n',
PRIMARY KEY (`id`),
UNIQUE KEY `name` (`name`)
) ENGINE=MyISAM AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `rover` (
`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
`loginname` varchar(30) NOT NULL,
`password` varchar(32) NOT NULL,
`status` enum('y','n') NOT NULL DEFAULT 'n',
PRIMARY KEY (`id`),
UNIQUE KEY `loginname` (`loginname`)
) ENGINE=MyISAM AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
