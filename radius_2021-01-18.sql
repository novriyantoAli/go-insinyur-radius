# ************************************************************
# Sequel Pro SQL dump
# Version 4541
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 103.55.36.93 (MySQL 5.5.5-10.3.25-MariaDB-0+deb10u1)
# Database: radius
# Generation Time: 2021-01-18 06:49:37 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table nas
# ------------------------------------------------------------

DROP TABLE IF EXISTS `nas`;

CREATE TABLE `nas` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `nasname` varchar(128) NOT NULL,
  `shortname` varchar(32) DEFAULT NULL,
  `type` varchar(30) DEFAULT 'other',
  `ports` int(5) DEFAULT NULL,
  `secret` varchar(60) NOT NULL DEFAULT 'secret',
  `server` varchar(64) DEFAULT NULL,
  `community` varchar(50) DEFAULT NULL,
  `description` varchar(200) DEFAULT 'RADIUS Client',
  PRIMARY KEY (`id`),
  KEY `nasname` (`nasname`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `nas` WRITE;
/*!40000 ALTER TABLE `nas` DISABLE KEYS */;

INSERT INTO `nas` (`id`, `nasname`, `shortname`, `type`, `ports`, `secret`, `server`, `community`, `description`)
VALUES
	(2,'10.8.0.2','MikroTik','other',NULL,'enKl01',NULL,NULL,'RADIUS Client'),
	(3,'0.0.0.0','Server','other',NULL,'s',NULL,NULL,'RADIUS Client');

/*!40000 ALTER TABLE `nas` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table radacct
# ------------------------------------------------------------

DROP TABLE IF EXISTS `radacct`;

CREATE TABLE `radacct` (
  `radacctid` bigint(21) NOT NULL AUTO_INCREMENT,
  `acctsessionid` varchar(64) NOT NULL DEFAULT '',
  `acctuniqueid` varchar(32) NOT NULL DEFAULT '',
  `username` varchar(64) NOT NULL DEFAULT '',
  `realm` varchar(64) DEFAULT '',
  `nasipaddress` varchar(15) NOT NULL DEFAULT '',
  `nasportid` varchar(15) DEFAULT NULL,
  `nasporttype` varchar(32) DEFAULT NULL,
  `acctstarttime` datetime DEFAULT NULL,
  `acctupdatetime` datetime DEFAULT NULL,
  `acctstoptime` datetime DEFAULT NULL,
  `acctinterval` int(12) DEFAULT NULL,
  `acctsessiontime` int(12) unsigned DEFAULT NULL,
  `acctauthentic` varchar(32) DEFAULT NULL,
  `connectinfo_start` varchar(50) DEFAULT NULL,
  `connectinfo_stop` varchar(50) DEFAULT NULL,
  `acctinputoctets` bigint(20) DEFAULT NULL,
  `acctoutputoctets` bigint(20) DEFAULT NULL,
  `calledstationid` varchar(50) NOT NULL DEFAULT '',
  `callingstationid` varchar(50) NOT NULL DEFAULT '',
  `acctterminatecause` varchar(32) NOT NULL DEFAULT '',
  `servicetype` varchar(32) DEFAULT NULL,
  `framedprotocol` varchar(32) DEFAULT NULL,
  `framedipaddress` varchar(15) NOT NULL DEFAULT '',
  PRIMARY KEY (`radacctid`),
  UNIQUE KEY `acctuniqueid` (`acctuniqueid`),
  KEY `username` (`username`),
  KEY `framedipaddress` (`framedipaddress`),
  KEY `acctsessionid` (`acctsessionid`),
  KEY `acctsessiontime` (`acctsessiontime`),
  KEY `acctstarttime` (`acctstarttime`),
  KEY `acctinterval` (`acctinterval`),
  KEY `acctstoptime` (`acctstoptime`),
  KEY `nasipaddress` (`nasipaddress`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `radacct` WRITE;
/*!40000 ALTER TABLE `radacct` DISABLE KEYS */;

INSERT INTO `radacct` (`radacctid`, `acctsessionid`, `acctuniqueid`, `username`, `realm`, `nasipaddress`, `nasportid`, `nasporttype`, `acctstarttime`, `acctupdatetime`, `acctstoptime`, `acctinterval`, `acctsessiontime`, `acctauthentic`, `connectinfo_start`, `connectinfo_stop`, `acctinputoctets`, `acctoutputoctets`, `calledstationid`, `callingstationid`, `acctterminatecause`, `servicetype`, `framedprotocol`, `framedipaddress`)
VALUES
	(1,'80100005','76e6caf9d099770236c79057d0d048da','123','','10.8.0.2','bridge','Wireless-802.11','2020-12-01 11:33:43','2020-12-01 11:33:43','2020-12-03 12:24:37',NULL,175854,'','','',0,0,'hotspot1','38:C9:86:26:EE:65','NAS-Reboot','','','192.168.88.254');

/*!40000 ALTER TABLE `radacct` ENABLE KEYS */;
UNLOCK TABLES;

DELIMITER ;;
/*!50003 SET SESSION SQL_MODE="STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION" */;;
/*!50003 CREATE */ /*!50017 DEFINER=`radius`@`%` */ /*!50003 TRIGGER `check_if_add_new_item_have_expiration` AFTER INSERT ON `radacct` FOR EACH ROW BEGIN

END */;;
DELIMITER ;
/*!50003 SET SESSION SQL_MODE=@OLD_SQL_MODE */;


# Dump of table radcheck
# ------------------------------------------------------------

DROP TABLE IF EXISTS `radcheck`;

CREATE TABLE `radcheck` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(64) NOT NULL DEFAULT '',
  `attribute` varchar(64) NOT NULL DEFAULT '',
  `op` char(2) NOT NULL DEFAULT '==',
  `value` varchar(253) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `username` (`username`(32))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `radcheck` WRITE;
/*!40000 ALTER TABLE `radcheck` DISABLE KEYS */;

INSERT INTO `radcheck` (`id`, `username`, `attribute`, `op`, `value`)
VALUES
	(59,'r','Cleartext-Password',':=','r');

/*!40000 ALTER TABLE `radcheck` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table radgroupcheck
# ------------------------------------------------------------

DROP TABLE IF EXISTS `radgroupcheck`;

CREATE TABLE `radgroupcheck` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `groupname` varchar(64) NOT NULL DEFAULT '',
  `attribute` varchar(64) NOT NULL DEFAULT '',
  `op` char(2) NOT NULL DEFAULT '==',
  `value` varchar(253) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `groupname` (`groupname`(32))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `radgroupcheck` WRITE;
/*!40000 ALTER TABLE `radgroupcheck` DISABLE KEYS */;

INSERT INTO `radgroupcheck` (`id`, `groupname`, `attribute`, `op`, `value`)
VALUES
	(1,'432','Simultaneous-Use',':=','1');

/*!40000 ALTER TABLE `radgroupcheck` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table radgroupreply
# ------------------------------------------------------------

DROP TABLE IF EXISTS `radgroupreply`;

CREATE TABLE `radgroupreply` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `groupname` varchar(64) NOT NULL DEFAULT '',
  `attribute` varchar(64) NOT NULL DEFAULT '',
  `op` char(2) NOT NULL DEFAULT '=',
  `value` varchar(253) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `groupname` (`groupname`(32))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `radgroupreply` WRITE;
/*!40000 ALTER TABLE `radgroupreply` DISABLE KEYS */;

INSERT INTO `radgroupreply` (`id`, `groupname`, `attribute`, `op`, `value`)
VALUES
	(1,'432','Mikrotik-Rate-Limit',':=','1M/2M');

/*!40000 ALTER TABLE `radgroupreply` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table radpostauth
# ------------------------------------------------------------

DROP TABLE IF EXISTS `radpostauth`;

CREATE TABLE `radpostauth` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(64) NOT NULL DEFAULT '',
  `pass` varchar(64) NOT NULL DEFAULT '',
  `reply` varchar(32) NOT NULL DEFAULT '',
  `authdate` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `radpostauth` WRITE;
/*!40000 ALTER TABLE `radpostauth` DISABLE KEYS */;

INSERT INTO `radpostauth` (`id`, `username`, `pass`, `reply`, `authdate`)
VALUES
	(1,'123','123','Access-Accept','2020-11-30 23:47:34'),
	(2,'123','123','Access-Accept','2020-12-01 09:19:48'),
	(3,'123','123','Access-Accept','2020-12-01 09:57:19'),
	(4,'123','123','Access-Accept','2020-12-01 10:04:16'),
	(5,'123','0x2e7d305d31d61845728081df013917c142','Access-Accept','2020-12-01 11:20:45'),
	(6,'123','0xb39c95db513c9029b685f836b642f75b0b','Access-Accept','2020-12-01 11:25:41'),
	(7,'123','0x52a3b02546300c537b0d0218c75690d94c','Access-Accept','2020-12-01 11:28:02'),
	(8,'123','0xc944599bf56978a80346d55d72c581b04b','Access-Accept','2020-12-01 11:33:43'),
	(9,'r','','Access-Accept','2021-01-13 14:23:19'),
	(10,'r','','Access-Accept','2021-01-13 14:23:19'),
	(11,'r','','Access-Accept','2021-01-13 17:08:13'),
	(12,'r','','Access-Accept','2021-01-13 17:08:13');

/*!40000 ALTER TABLE `radpostauth` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table radreply
# ------------------------------------------------------------

DROP TABLE IF EXISTS `radreply`;

CREATE TABLE `radreply` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(64) NOT NULL DEFAULT '',
  `attribute` varchar(64) NOT NULL DEFAULT '',
  `op` char(2) NOT NULL DEFAULT '=',
  `value` varchar(253) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `username` (`username`(32))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `radreply` WRITE;
/*!40000 ALTER TABLE `radreply` DISABLE KEYS */;

INSERT INTO `radreply` (`id`, `username`, `attribute`, `op`, `value`)
VALUES
	(23,'pengguna_root','Mikrotik-Rate-Limit','=','341k/341k'),
	(24,'pengguna_child','Mikrotik-Rate-Limit','=','341k/341k'),
	(25,'toor','Mikrotik-Rate-Limit','=','341k/341k');

/*!40000 ALTER TABLE `radreply` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table radusergroup
# ------------------------------------------------------------

DROP TABLE IF EXISTS `radusergroup`;

CREATE TABLE `radusergroup` (
  `username` varchar(64) NOT NULL DEFAULT '',
  `groupname` varchar(64) NOT NULL DEFAULT '',
  `priority` int(11) NOT NULL DEFAULT 1,
  KEY `username` (`username`(32))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `radusergroup` WRITE;
/*!40000 ALTER TABLE `radusergroup` DISABLE KEYS */;

INSERT INTO `radusergroup` (`username`, `groupname`, `priority`)
VALUES
	('432','432',1);

/*!40000 ALTER TABLE `radusergroup` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table transaction
# ------------------------------------------------------------

DROP TABLE IF EXISTS `transaction`;

CREATE TABLE `transaction` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `id_sender` int(11) unsigned DEFAULT NULL,
  `id_receive` int(11) unsigned NOT NULL,
  `id_company_packet` int(11) unsigned DEFAULT NULL,
  `value` int(11) NOT NULL,
  `information` text DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
