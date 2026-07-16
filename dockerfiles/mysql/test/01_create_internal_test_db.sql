-- internal_test スキーマ（テスト用審査スキーマ）を作成する初期化スクリプト
CREATE DATABASE IF NOT EXISTS `internal_test` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON `internal_test`.* TO 'tester'@'%';
FLUSH PRIVILEGES;
