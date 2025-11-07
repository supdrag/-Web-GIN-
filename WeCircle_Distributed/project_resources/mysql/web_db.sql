-- MySQL dump 10.13  Distrib 8.0.39, for Win64 (x86_64)
--
-- Host: localhost    Database: web_db_copy
-- ------------------------------------------------------
-- Server version	8.0.39

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `circle_messages`
--

DROP TABLE IF EXISTS `circle_messages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `circle_messages` (
  `MESSAGE_ID` bigint NOT NULL AUTO_INCREMENT,
  `USER_ID` int NOT NULL COMMENT '用户ID',
  `CIRCLE_ID` int NOT NULL COMMENT '圈子ID',
  `CONTENT` varchar(512) NOT NULL COMMENT '消息内容',
  `TIME` datetime NOT NULL COMMENT '发送时间',
  `USER_NAME` varchar(30) NOT NULL COMMENT '用户名',
  PRIMARY KEY (`MESSAGE_ID`)
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `circle_messages`
--

LOCK TABLES `circle_messages` WRITE;
/*!40000 ALTER TABLE `circle_messages` DISABLE KEYS */;
INSERT INTO `circle_messages` VALUES (1,1,1,'Hello!My name is Super Xiaoming.Nice to meet you!','2025-10-07 22:39:22','SuperMing'),(2,1,1,'From now on, I’ll be the host of this managers’ group chat family—thank you all for your support! I believe that with our joint efforts, a bright future will surely come!','2025-10-07 22:55:10','SuperMing'),(3,1,1,'Thank you! And I love you!','2025-10-07 22:55:29','SuperMing'),(4,2,1,'Hello, my name is Yuyi.Nice to meet you!','2025-10-08 15:03:40','YuYi'),(5,1,1,'哈哈哈，大家可以说中文的！欢迎你，宇熠！','2025-10-08 15:04:40','SuperMing'),(6,2,1,'谢谢你，小明！','2025-10-08 15:05:28','YuYi'),(7,1,1,'乌拉呀哈乌拉！','2025-10-08 15:13:09','SuperMing'),(8,2,1,'我回来啦！','2025-10-08 16:54:01','YuYi'),(9,1,1,'hi','2025-10-11 20:46:12','SuperMing'),(10,1,1,'hello','2025-10-11 20:50:24','SuperMing'),(11,2,1,'hello!','2025-10-11 20:50:30','YuYi'),(12,3,1,'hi','2025-10-11 22:40:05','KangKang'),(13,3,1,'hello','2025-10-12 16:40:27','KangKang'),(14,2,1,'Are you OK?','2025-10-12 16:54:03','YuYi'),(15,3,1,'May be fine?','2025-10-12 16:54:13','KangKang'),(16,2,1,'Hi,bro? Are you there?','2025-10-12 17:50:10','YuYi'),(17,1,1,'He is kicked by me for a while just now.Haha!','2025-10-12 18:05:48','SuperMing'),(18,2,1,'Oh,it\'s funny.Haha!','2025-10-12 18:06:16','YuYi'),(19,1,1,'OK,I\'ll invite him to come here again.','2025-10-12 18:07:02','SuperMing'),(20,2,1,'Waiting for you haha!','2025-10-12 18:07:50','YuYi'),(21,3,1,'What happened to me??? I exit abnormally.God!','2025-10-12 18:11:10','KangKang'),(22,1,1,'Maybe the ghost did. Haha!','2025-10-12 18:11:35','SuperMing'),(23,2,1,'Yeah,maybe.Haha','2025-10-12 18:11:48','YuYi'),(24,3,1,'It\'s not funny,bros.','2025-10-12 18:12:23','KangKang'),(28,1,2,'\n--==[ROLE CHANGE]==--\n<|User_ID:2|>\n[Common Member]=>[Circle Manager]\n---------------------\n','2025-10-14 22:49:12','SuperMing'),(29,3,1,'Hello,guys.It\'s a good day!','2025-10-26 12:51:12','SuperMing'),(30,2,1,'Hello!','2025-10-26 12:56:54','LiTianyu'),(31,3,1,'exit','2025-10-26 13:01:15','SuperMing');
/*!40000 ALTER TABLE `circle_messages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `circle_ships`
--

DROP TABLE IF EXISTS `circle_ships`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `circle_ships` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `USER_ID` int NOT NULL COMMENT '用户ID',
  `CIRCLE_ID` int NOT NULL COMMENT '圈子ID',
  `ROLE` int NOT NULL COMMENT '用户身份',
  `STATUS` int NOT NULL COMMENT '关系状态',
  `TIME` datetime NOT NULL COMMENT '关系创建时间',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `circle_ships`
--

LOCK TABLES `circle_ships` WRITE;
/*!40000 ALTER TABLE `circle_ships` DISABLE KEYS */;
INSERT INTO `circle_ships` VALUES (1,1,1,1,1,'2025-10-07 15:07:00'),(2,3,1,2,1,'2025-10-07 15:07:00'),(4,4,1,2,1,'2025-10-07 15:07:00'),(5,1,2,1,1,'2025-10-08 19:57:00'),(6,2,1,2,1,'2025-10-07 15:07:00'),(19,2,2,2,1,'2025-10-12 19:40:45'),(20,3,2,3,1,'2025-10-14 16:55:35');
/*!40000 ALTER TABLE `circle_ships` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `circles`
--

DROP TABLE IF EXISTS `circles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `circles` (
  `CIRCLE_ID` int NOT NULL AUTO_INCREMENT,
  `PROFILE` varchar(512) NOT NULL COMMENT '群聊介绍',
  `NUM` int NOT NULL COMMENT '群聊当前人数',
  `LMT_NUM` int NOT NULL COMMENT '群聊人数上限',
  `STATUS` int NOT NULL COMMENT '群聊状态',
  `TIME` datetime NOT NULL COMMENT '群聊创建时间',
  PRIMARY KEY (`CIRCLE_ID`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `circles`
--

LOCK TABLES `circles` WRITE;
/*!40000 ALTER TABLE `circles` DISABLE KEYS */;
INSERT INTO `circles` VALUES (1,'本交流平台最权威的圈子，核心工作人员都在里面。',4,50,1,'2025-10-07 15:03:00'),(2,'官方建立的第一个闲聊群，没事可进来唠嗑两句。',3,200,1,'2025-10-08 19:54:00');
/*!40000 ALTER TABLE `circles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `collects`
--

DROP TABLE IF EXISTS `collects`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `collects` (
  `ID` bigint NOT NULL AUTO_INCREMENT COMMENT '收藏ID',
  `USER_ID` int NOT NULL COMMENT '用户ID',
  `DATA` varchar(255) NOT NULL COMMENT '收藏内容',
  `TIME` datetime NOT NULL COMMENT '收藏时间',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `collects`
--

LOCK TABLES `collects` WRITE;
/*!40000 ALTER TABLE `collects` DISABLE KEYS */;
INSERT INTO `collects` VALUES (1,3,'这是一个很重要的链接:https://www.bilibili.com/','2025-10-21 20:53:45');
/*!40000 ALTER TABLE `collects` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `friend_ships`
--

DROP TABLE IF EXISTS `friend_ships`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `friend_ships` (
  `ID` bigint NOT NULL AUTO_INCREMENT COMMENT '好友关系',
  `SMALL_ID` int NOT NULL COMMENT '较小用户ID',
  `BIG_ID` int NOT NULL COMMENT '较大用户ID',
  `STATUS_SMALL` int NOT NULL COMMENT '关系状态',
  `TIME` datetime NOT NULL COMMENT '创建时间',
  `STATUS_BIG` int NOT NULL COMMENT '关系状态',
  PRIMARY KEY (`ID`),
  UNIQUE KEY `SMALL_BIG_ID` (`SMALL_ID`,`BIG_ID`),
  CONSTRAINT `sma_bg_id` CHECK ((`SMALL_ID` < `BIG_ID`))
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `friend_ships`
--

LOCK TABLES `friend_ships` WRITE;
/*!40000 ALTER TABLE `friend_ships` DISABLE KEYS */;
INSERT INTO `friend_ships` VALUES (1,1,2,1,'2025-10-05 21:41:00',1),(2,1,3,1,'2025-10-05 21:41:00',1),(3,1,4,1,'2025-10-05 21:41:00',1),(4,2,3,1,'2025-10-05 21:41:00',1),(5,2,4,1,'2025-10-05 21:41:00',1),(6,3,4,1,'2025-10-05 21:41:00',1);
/*!40000 ALTER TABLE `friend_ships` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `game_recommends`
--

DROP TABLE IF EXISTS `game_recommends`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `game_recommends` (
  `ID` bigint NOT NULL AUTO_INCREMENT COMMENT '游戏推荐ID',
  `USER_ID` int NOT NULL COMMENT '用户ID',
  `GAME_URL` varchar(255) DEFAULT NULL COMMENT '游戏链接',
  `REASON` varchar(255) NOT NULL COMMENT '推荐理由',
  `SCORE` int DEFAULT NULL COMMENT '推荐打分',
  `TIME` datetime NOT NULL COMMENT '推荐时间',
  PRIMARY KEY (`ID`),
  CONSTRAINT `scr_limi` CHECK (((`SCORE` >= 0) and (`SCORE` <= 100)))
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `game_recommends`
--

LOCK TABLES `game_recommends` WRITE;
/*!40000 ALTER TABLE `game_recommends` DISABLE KEYS */;
INSERT INTO `game_recommends` VALUES (1,3,'','',0,'2025-10-21 17:14:43'),(2,3,'https://www.lol.com/chat/22836917398472962','good game!',90,'2025-10-21 17:17:30');
/*!40000 ALTER TABLE `game_recommends` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `goods_comments`
--

DROP TABLE IF EXISTS `goods_comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `goods_comments` (
  `COMMENT_ID` bigint NOT NULL AUTO_INCREMENT COMMENT '商品评论ID',
  `GOODS_CODE` varchar(50) NOT NULL COMMENT '商品编号',
  `USER_ID` int NOT NULL COMMENT '用户ID',
  `TIME` datetime NOT NULL COMMENT '发表时间',
  `CONTENT` varchar(255) NOT NULL COMMENT '商品评论',
  `SCORE` int NOT NULL COMMENT '打分',
  PRIMARY KEY (`COMMENT_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `goods_comments`
--

LOCK TABLES `goods_comments` WRITE;
/*!40000 ALTER TABLE `goods_comments` DISABLE KEYS */;
/*!40000 ALTER TABLE `goods_comments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `goods_datas`
--

DROP TABLE IF EXISTS `goods_datas`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `goods_datas` (
  `GOODS_CODE` varchar(50) NOT NULL COMMENT '商品编号',
  `NAME` varchar(50) NOT NULL COMMENT '商品名称',
  `PRICE` float NOT NULL COMMENT '商品价格',
  `TYPE` varchar(30) NOT NULL COMMENT '商品类型',
  `LIKES` int NOT NULL COMMENT '商品推荐人数',
  `SCORE` int NOT NULL COMMENT '商品评分',
  `PROFILE` varchar(150) NOT NULL COMMENT '商品介绍',
  `LOGO` varchar(20) NOT NULL COMMENT '商品品牌',
  `URL` varchar(150) DEFAULT NULL COMMENT '跳转链接',
  PRIMARY KEY (`GOODS_CODE`),
  UNIQUE KEY `Name` (`NAME`),
  UNIQUE KEY `URL` (`URL`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `goods_datas`
--

LOCK TABLES `goods_datas` WRITE;
/*!40000 ALTER TABLE `goods_datas` DISABLE KEYS */;
/*!40000 ALTER TABLE `goods_datas` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `goods_recommends`
--

DROP TABLE IF EXISTS `goods_recommends`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `goods_recommends` (
  `ID` bigint NOT NULL AUTO_INCREMENT COMMENT '商品推荐ID',
  `USER_ID` int NOT NULL COMMENT '用户ID',
  `GOODS_URL` varchar(255) NOT NULL COMMENT '商品链接',
  `REASON` varchar(255) NOT NULL COMMENT '推荐原因',
  `SCORE` int DEFAULT NULL COMMENT '推荐打分',
  `TIME` datetime NOT NULL COMMENT '推荐时间',
  PRIMARY KEY (`ID`),
  CONSTRAINT `score_limit` CHECK (((`SCORE` >= 0) and (`SCORE` <= 100)))
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `goods_recommends`
--

LOCK TABLES `goods_recommends` WRITE;
/*!40000 ALTER TABLE `goods_recommends` DISABLE KEYS */;
INSERT INTO `goods_recommends` VALUES (1,3,'https://www.doubao.com/chat/22836917398472962','good!',0,'2025-10-21 15:54:00');
/*!40000 ALTER TABLE `goods_recommends` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `moment_comments`
--

DROP TABLE IF EXISTS `moment_comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `moment_comments` (
  `COMMENT_ID` bigint NOT NULL AUTO_INCREMENT COMMENT '评论ID',
  `MOMENT_ID` bigint NOT NULL COMMENT '所属动态ID',
  `USER_ID` int NOT NULL COMMENT '用户ID',
  `TIME` datetime NOT NULL COMMENT '发表时间',
  `CONTENT` varchar(255) NOT NULL COMMENT '评论内容',
  PRIMARY KEY (`COMMENT_ID`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `moment_comments`
--

LOCK TABLES `moment_comments` WRITE;
/*!40000 ALTER TABLE `moment_comments` DISABLE KEYS */;
INSERT INTO `moment_comments` VALUES (1,2,3,'2025-10-21 13:00:42','哈哈哈，期待作者大大后续更新！');
/*!40000 ALTER TABLE `moment_comments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `moments`
--

DROP TABLE IF EXISTS `moments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `moments` (
  `MOMENT_ID` bigint NOT NULL AUTO_INCREMENT COMMENT '动态ID',
  `USER_ID` int NOT NULL COMMENT '发表用户ID',
  `CONTENT` varchar(400) NOT NULL COMMENT '动态内容',
  `TIME` datetime NOT NULL COMMENT '发表时间',
  `COMMENT_NUM` int DEFAULT NULL COMMENT '评论数',
  `LIKES` int DEFAULT NULL COMMENT '点赞数',
  `STATUS` int DEFAULT NULL COMMENT '查看权限',
  PRIMARY KEY (`MOMENT_ID`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `moments`
--

LOCK TABLES `moments` WRITE;
/*!40000 ALTER TABLE `moments` DISABLE KEYS */;
INSERT INTO `moments` VALUES (2,3,'这是我的第一条动态，祝大家幸运、开心每一天！','2025-10-21 09:39:15',1,1,1);
/*!40000 ALTER TABLE `moments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `playings`
--

DROP TABLE IF EXISTS `playings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `playings` (
  `USER_ID` int NOT NULL COMMENT '用户ID',
  `GAME_NAME` varchar(30) NOT NULL COMMENT '游戏名'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `playings`
--

LOCK TABLES `playings` WRITE;
/*!40000 ALTER TABLE `playings` DISABLE KEYS */;
/*!40000 ALTER TABLE `playings` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `solo_messages`
--

DROP TABLE IF EXISTS `solo_messages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `solo_messages` (
  `MESSAGE_ID` bigint NOT NULL AUTO_INCREMENT COMMENT '私聊信息ID',
  `USER_ID` int NOT NULL COMMENT '发送用户ID',
  `USER_NAME` varchar(30) NOT NULL COMMENT '用户名称',
  `FRIEND_ID` int NOT NULL COMMENT '接收方ID',
  `CONTENT` varchar(512) NOT NULL COMMENT '消息内容',
  `TIME` datetime NOT NULL COMMENT '发送时间',
  PRIMARY KEY (`MESSAGE_ID`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `solo_messages`
--

LOCK TABLES `solo_messages` WRITE;
/*!40000 ALTER TABLE `solo_messages` DISABLE KEYS */;
INSERT INTO `solo_messages` VALUES (1,3,'SuperMing',2,'Hi','2025-10-17 22:32:13'),(2,2,'LiTianyu',3,'Hi,I come just now','2025-10-17 22:44:10'),(3,3,'SuperMing',2,'What do you do these days?','2025-10-17 22:44:40'),(4,2,'LiTianyu',3,'Playing games;','2025-10-17 22:45:04'),(5,3,'SuperMing',2,'Me too,bro.So boring.','2025-10-17 22:47:05'),(6,2,'LiTianyu',3,'Do you have any intersting ideas to pass the time.','2025-10-17 22:47:56'),(7,3,'SuperMing',2,'Nothing now.','2025-10-17 22:48:07'),(8,3,'SuperMing',2,'It\'s a little late now.See you tommorow.','2025-10-17 22:51:15'),(9,2,'LiTianyu',3,'Bye!','2025-10-17 22:51:19'),(10,2,'LiTianyu',3,'emm','2025-10-19 14:05:27'),(11,2,'LiTianyu',3,'em','2025-10-19 14:05:57'),(12,3,'SuperMing',2,'hello','2025-10-26 13:05:47'),(13,2,'LiTianyu',3,'Long time no see!','2025-10-26 13:06:45'),(14,3,'SuperMing',2,'How are you today?','2025-10-26 13:08:26'),(15,2,'LiTianyu',3,'Not very good,and you?','2025-10-26 13:08:40'),(16,3,'SuperMing',2,'Yeah.Me,too','2025-10-26 13:08:56'),(17,3,'SuperMing',2,'Bye!','2025-10-26 13:10:55'),(18,2,'LiTianyu',3,'Bye','2025-10-26 13:11:00'),(19,2,'LiTianyu',3,'[exit]','2025-10-26 14:54:07');
/*!40000 ALTER TABLE `solo_messages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_profiles`
--

DROP TABLE IF EXISTS `user_profiles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_profiles` (
  `USER_ID` int NOT NULL COMMENT '用户ID',
  `SIGNATURE` varchar(255) NOT NULL COMMENT '用户签名',
  `POPULARITY` int NOT NULL COMMENT '用户人气',
  `AGE` int NOT NULL COMMENT '用户年龄',
  `GENDER` tinyint(1) NOT NULL COMMENT '用户性别',
  `LOCATION` varchar(50) NOT NULL COMMENT '用户住址',
  `JOB` varchar(30) NOT NULL COMMENT '工作',
  `TIME` datetime NOT NULL COMMENT '最近更新时间',
  PRIMARY KEY (`USER_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_profiles`
--

LOCK TABLES `user_profiles` WRITE;
/*!40000 ALTER TABLE `user_profiles` DISABLE KEYS */;
INSERT INTO `user_profiles` VALUES (1,'',0,21,1,'Chongqing','Student','2025-09-21 21:39:00'),(2,'',0,21,1,'Chongqing','Student','2025-09-21 21:39:00'),(3,'你好，世界！',0,21,1,'Chongqing','Student','2025-10-21 15:43:06'),(4,'',0,21,1,'Chongqing','Student','2025-09-21 21:39:00');
/*!40000 ALTER TABLE `user_profiles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `ID` int NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `NAME` varchar(30) NOT NULL COMMENT '用户名字',
  `ACCOUNT` varchar(20) NOT NULL COMMENT '用户账号',
  `CTF` varchar(10) DEFAULT NULL COMMENT '用户身份',
  `CreateAt` datetime DEFAULT NULL COMMENT '用户创建时间',
  `UpdateAt` datetime DEFAULT NULL COMMENT '用户更新时间',
  `PASSWD` varchar(30) NOT NULL COMMENT '用户密码',
  `STATUS` int NOT NULL COMMENT '用户状态',
  `PHONE` varchar(20) DEFAULT NULL COMMENT '用户电话',
  PRIMARY KEY (`ID`),
  UNIQUE KEY `NAME` (`NAME`),
  UNIQUE KEY `QQ` (`ACCOUNT`),
  UNIQUE KEY `Account` (`ACCOUNT`),
  UNIQUE KEY `Acocount` (`ACCOUNT`),
  UNIQUE KEY `ACCOUNT_2` (`ACCOUNT`),
  UNIQUE KEY `PHONE` (`PHONE`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'KangKang','8888888888','Common','2025-09-21 21:39:00','2025-09-21 21:39:00','123456789',1,'123456789aa'),(2,'LiTianyu','2524212569','common','2025-10-03 20:13:20','2025-10-03 20:13:20','123456789',1,'123456789BB'),(3,'SuperMing','3542358746','Master','2025-09-21 21:33:00','2025-09-21 21:33:00','123456789',1,'123456789CC'),(4,'YuYi','8888866666','Manager','2025-09-21 21:37:00','2025-09-21 21:37:00','123456789',1,'123456789DD');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `visit_logs`
--

DROP TABLE IF EXISTS `visit_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `visit_logs` (
  `CODE` bigint NOT NULL AUTO_INCREMENT COMMENT '访问记录编号',
  `USER_ID` int NOT NULL COMMENT '用户id',
  `PATH` varchar(150) NOT NULL COMMENT '访问路径',
  `TYPE` varchar(10) NOT NULL COMMENT '访问类型',
  `TIME` datetime NOT NULL COMMENT '访问时间',
  PRIMARY KEY (`CODE`)
) ENGINE=InnoDB AUTO_INCREMENT=145 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `visit_logs`
--

LOCK TABLES `visit_logs` WRITE;
/*!40000 ALTER TABLE `visit_logs` DISABLE KEYS */;
INSERT INTO `visit_logs` VALUES (1,1,'/url/user/friends','GET','2025-10-16 21:16:55'),(2,1,'/url/user/friends','GET','2025-10-16 21:18:21'),(3,1,'/url/user/friends','GET','2025-10-16 21:19:05'),(4,1,'/url/user/friends','GET','2025-10-16 21:48:12'),(5,1,'/url/user/friends/contact','GET','2025-10-17 18:42:11'),(6,3,'/url/user/friends/contact','GET','2025-10-17 21:47:09'),(7,3,'/url/user/friends/contact','GET','2025-10-17 21:53:08'),(8,3,'/url/user/friends/contact','GET','2025-10-17 21:53:35'),(9,3,'/url/user/friends/contact','GET','2025-10-17 22:01:18'),(10,3,'/url/user/friends/contact','GET','2025-10-17 22:03:57'),(11,3,'/url/user/friends/contact','GET','2025-10-17 22:08:25'),(12,3,'/url/user/friends/contact','GET','2025-10-17 22:10:25'),(13,3,'/url/user/friends/contact','GET','2025-10-17 22:18:07'),(14,3,'/url/user/friends/contact','GET','2025-10-17 22:24:08'),(15,3,'/url/user/friends/contact','GET','2025-10-17 22:32:10'),(16,3,'/url/user/friends/contact','GET','2025-10-17 22:33:45'),(17,2,'/url/user/friends/contact','GET','2025-10-17 22:38:55'),(18,2,'/url/user/friends/contact','GET','2025-10-17 22:41:43'),(19,3,'/url/user/friends/contact','GET','2025-10-17 22:41:45'),(20,3,'/url/user/friends/contact','GET','2025-10-17 22:43:35'),(21,2,'/url/user/friends/contact','GET','2025-10-17 22:43:39'),(22,2,'/url/user/friends/contact','GET','2025-10-17 22:46:50'),(23,3,'/url/user/friends/contact','GET','2025-10-17 22:46:52'),(24,3,'/url/user/friends/contact','GET','2025-10-17 22:50:36'),(25,2,'/url/user/friends/contact','GET','2025-10-17 22:50:43'),(26,3,'/url/user/friends','GET','2025-10-18 12:48:14'),(27,3,'/url/user/friends','GET','2025-10-18 16:23:30'),(28,2,'/url/user/friends','GET','2025-10-18 16:24:34'),(29,2,'/url/user/myprofile','GET','2025-10-18 17:51:52'),(30,2,'/url/user/myprofile','GET','2025-10-18 17:55:10'),(31,2,'/url/user/myprofile','GET','2025-10-18 17:56:05'),(32,2,'/url/user/myprofile','GET','2025-10-18 18:03:28'),(33,2,'/url/user/myprofile','GET','2025-10-18 18:03:56'),(34,2,'/url/user/myprofile','GET','2025-10-19 10:37:58'),(35,2,'/url/user/myprofile','GET','2025-10-19 10:38:10'),(36,2,'/url/user/myprofile','GET','2025-10-19 10:38:23'),(37,2,'/url/user/myprofile','GET','2025-10-19 10:42:26'),(38,2,'/url/user/myprofile','GET','2025-10-19 10:42:44'),(39,2,'/url/user/myprofile','GET','2025-10-19 12:11:15'),(40,3,'/url/user/friends','GET','2025-10-19 12:16:24'),(41,2,'/url/user/myprofile','GET','2025-10-19 12:19:53'),(42,5,'/url/user/cancel','GET','2025-10-19 12:29:50'),(43,5,'/url/user/cancel','DELETE','2025-10-19 12:30:03'),(44,5,'/url/user/cancel','DELETE','2025-10-19 12:31:21'),(45,5,'/url/user/cancel','DELETE','2025-10-19 12:32:27'),(46,2,'/url/user/history','GET','2025-10-19 13:05:47'),(47,2,'/url/user/history','GET','2025-10-19 13:07:24'),(48,2,'/url/user/history','GET','2025-10-19 13:11:24'),(49,2,'/url/user/history','GET','2025-10-19 13:12:41'),(50,2,'/url/user/history','GET','2025-10-19 13:13:54'),(51,3,'/url/user/friends/manage','GET','2025-10-19 13:51:05'),(52,2,'/url/user/friends/manage','GET','2025-10-19 13:55:59'),(53,3,'/url/user/friends/contact','GET','2025-10-19 13:56:29'),(54,3,'/url/user/friends/contact','GET','2025-10-19 13:56:59'),(55,2,'/url/user/friends/contact','GET','2025-10-19 13:57:41'),(56,2,'/url/user/friends/contact','GET','2025-10-19 14:01:51'),(57,3,'/url/user/friends/contact','GET','2025-10-19 14:01:56'),(58,3,'/url/user/friends/contact','GET','2025-10-19 14:05:18'),(59,2,'/url/user/friends/contact','GET','2025-10-19 14:05:22'),(60,2,'/url/user/friends/manage','GET','2025-10-19 14:05:52'),(61,2,'/url/user/friends/contact','GET','2025-10-19 14:05:56'),(62,2,'/url/user/friends/manage','GET','2025-10-19 14:10:08'),(63,2,'/url/user/friends/manage','GET','2025-10-19 14:12:22'),(64,2,'/url/user/friends/manage','GET','2025-10-19 14:15:43'),(65,2,'/url/user/friends/manage','GET','2025-10-19 14:18:08'),(66,2,'/url/user/friends/manage','GET','2025-10-19 14:20:01'),(67,2,'/url/user/friends/contact','GET','2025-10-19 14:20:45'),(68,3,'/url/user/friends/contact','GET','2025-10-19 14:20:50'),(69,3,'/url/user/profile','GET','2025-10-21 09:20:39'),(70,3,'/url/user/profile','GET','2025-10-21 09:23:39'),(71,3,'/url/user/profile','GET','2025-10-21 09:24:25'),(72,3,'/url/user/profile','GET','2025-10-21 09:25:13'),(73,3,'/url/user/profile','GET','2025-10-21 09:26:50'),(74,3,'/url/user/profile','GET','2025-10-21 09:27:23'),(75,3,'/url/user/profile','GET','2025-10-21 09:27:31'),(76,3,'/url/user/profile','GET','2025-10-21 09:39:15'),(77,3,'/url/user/profile','GET','2025-10-21 09:39:21'),(78,3,'/url/user/profile','GET','2025-10-21 09:40:11'),(79,3,'/url/user/profile','GET','2025-10-21 09:41:04'),(80,3,'/url/user/interaction','GET','2025-10-21 12:43:20'),(81,3,'/url/user/friends/interaction','GET','2025-10-21 12:44:15'),(82,3,'/url/user/friends/interaction','GET','2025-10-21 12:49:45'),(83,3,'/url/user/profile','GET','2025-10-21 12:49:53'),(84,3,'/url/user/friends/interaction','GET','2025-10-21 12:58:54'),(85,3,'/url/user/friends/interaction','GET','2025-10-21 13:00:42'),(86,3,'/url/user/profile','GET','2025-10-21 13:00:48'),(87,3,'/url/user/profile','GET','2025-10-21 13:01:02'),(88,3,'/url/user/profile','GET','2025-10-21 13:04:01'),(89,3,'/url/user/profile','GET','2025-10-21 14:03:45'),(90,3,'/url/user/profile','GET','2025-10-21 14:03:50'),(91,3,'/url/user/profile','GET','2025-10-21 14:36:33'),(92,3,'/url/user/profile','GET','2025-10-21 14:37:02'),(93,3,'/url/user/profile','GET','2025-10-21 14:37:25'),(94,3,'/url/user/profile','GET','2025-10-21 14:37:28'),(95,3,'/url/user/profile','GET','2025-10-21 14:39:27'),(96,3,'/url/user/profile','GET','2025-10-21 14:39:50'),(97,3,'/url/user/profile','GET','2025-10-21 14:39:57'),(98,3,'/url/user/profile','GET','2025-10-21 14:43:20'),(99,3,'/url/user/profile','GET','2025-10-21 14:43:22'),(100,3,'/url/user/profile','GET','2025-10-21 14:43:35'),(101,3,'/url/user/profile','GET','2025-10-21 14:43:37'),(102,3,'/url/user/profile','GET','2025-10-21 14:44:07'),(103,3,'/url/user/profile','GET','2025-10-21 14:44:08'),(104,3,'/url/user/profile','GET','2025-10-21 15:00:03'),(105,3,'/url/user/profile','GET','2025-10-21 15:00:06'),(106,3,'/url/user/profile','GET','2025-10-21 15:30:19'),(107,3,'/url/user/profile','GET','2025-10-21 15:31:14'),(108,3,'/url/user/profile','GET','2025-10-21 15:31:28'),(109,3,'/url/user/profile','GET','2025-10-21 15:41:52'),(110,3,'/url/user/profile','GET','2025-10-21 15:42:01'),(111,3,'/url/user/profile','GET','2025-10-21 15:43:00'),(112,3,'/url/user/profile','GET','2025-10-21 15:43:06'),(113,3,'/url/user/profile','GET','2025-10-21 15:43:11'),(114,3,'/url/user/profile','GET','2025-10-21 15:53:39'),(115,3,'/url/user/profile','GET','2025-10-21 15:54:00'),(116,3,'/url/user/profile','GET','2025-10-21 15:54:30'),(117,3,'/url/user/profile','GET','2025-10-21 15:58:05'),(118,3,'/url/user/profile','GET','2025-10-21 15:59:47'),(119,3,'/url/user/profile','GET','2025-10-21 16:28:36'),(120,3,'/url/user/profile','GET','2025-10-21 16:29:08'),(121,3,'/url/user/profile','GET','2025-10-21 16:29:27'),(122,3,'/url/user/profile','GET','2025-10-21 17:13:47'),(123,3,'/url/user/profile','GET','2025-10-21 17:14:20'),(124,3,'/url/user/profile','GET','2025-10-21 17:14:43'),(125,3,'/url/user/profile','GET','2025-10-21 17:17:30'),(126,3,'/url/user/profile','GET','2025-10-21 20:53:45'),(127,3,'/url/user/profile','GET','2025-10-21 20:54:03'),(128,3,'/url/user/profile','GET','2025-10-26 12:44:21'),(129,3,'/url/user/profile','GET','2025-10-26 12:44:33'),(130,3,'/url/user/profile','GET','2025-10-26 12:44:51'),(131,3,'/url/user/profile','GET','2025-10-26 12:45:16'),(132,3,'/url/user/profile','GET','2025-10-26 12:46:16'),(133,3,'/url/user/profile','GET','2025-10-26 12:48:45'),(134,3,'/url/chat/chatroom/contact','GET','2025-10-26 12:50:15'),(135,3,'/url/chat/chatroom/contact','GET','2025-10-26 12:54:15'),(136,2,'/url/chat/chatroom/contact','GET','2025-10-26 12:56:47'),(137,3,'/url/chat/chatroom/contact','GET','2025-10-26 13:01:12'),(138,3,'/url/user/friends/contact','GET','2025-10-26 13:03:00'),(139,3,'/url/user/friends/contact','GET','2025-10-26 13:05:43'),(140,2,'/url/user/friends/contact','GET','2025-10-26 13:06:19'),(141,2,'/url/user/friends/contact','GET','2025-10-26 13:06:35'),(142,2,'/url/user/friends/contact','GET','2025-10-26 13:16:44'),(143,3,'/url/user/friends/contact','GET','2025-10-26 13:16:55'),(144,2,'/url/user/friends/contact','GET','2025-10-26 14:54:03');
/*!40000 ALTER TABLE `visit_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `voter`
--

DROP TABLE IF EXISTS `voter`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `voter` (
  `USER_ID` int NOT NULL COMMENT '投票用户ID',
  `GOODS_CODE` varchar(50) NOT NULL COMMENT '商品编号',
  `TIME` datetime NOT NULL COMMENT '投票时间',
  `VOTE_ID` int NOT NULL COMMENT '投票编号',
  PRIMARY KEY (`VOTE_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `voter`
--

LOCK TABLES `voter` WRITE;
/*!40000 ALTER TABLE `voter` DISABLE KEYS */;
/*!40000 ALTER TABLE `voter` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping events for database 'web_db_copy'
--

--
-- Dumping routines for database 'web_db_copy'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-10-26 17:40:17
