-- internal スキーマ（審査用）を作成する初期化スクリプト
CREATE DATABASE IF NOT EXISTS `internal` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON `internal`.* TO 'user'@'%';
FLUSH PRIVILEGES;
