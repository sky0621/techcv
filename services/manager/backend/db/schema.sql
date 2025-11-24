CREATE TABLE users (
  id BINARY(16) NOT NULL COMMENT 'UUID v7（アプリケーション内部ID）',
  firebase_uid VARCHAR(128) NOT NULL COMMENT 'Firebase UID',
  email VARCHAR(255) NOT NULL COMMENT 'メールアドレス',
  email_verified TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'メール確認済みフラグ（Firebaseから取得）',
  display_name VARCHAR(255) DEFAULT NULL COMMENT '表示名（Firebaseから取得）',
  photo_url VARCHAR(500) DEFAULT NULL COMMENT 'プロフィール画像URL（Firebaseから取得）',
  phone_number VARCHAR(50) DEFAULT NULL COMMENT '電話番号（Firebaseから取得、オプション）',
  provider_id VARCHAR(50) NOT NULL COMMENT '認証プロバイダーID（例: google.com）',
  firebase_created_at DATETIME(6) DEFAULT NULL COMMENT 'Firebase上のアカウント作成日時',
  firebase_last_sign_in_at DATETIME(6) DEFAULT NULL COMMENT 'Firebase上の最終サインイン日時',
  bio TEXT DEFAULT NULL COMMENT '自己紹介（アプリケーション独自項目）',
  is_active TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'アクティブ状態',
  email_verified_at DATETIME(6) DEFAULT NULL COMMENT 'メール確認日時（アプリケーション側で管理）',
  last_login_at DATETIME(6) DEFAULT NULL COMMENT '最終ログイン日時（アプリケーション側で管理）',
  created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '作成日時',
  updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新日時',
  deleted_at DATETIME(6) DEFAULT NULL COMMENT '削除日時（論理削除用）',
  PRIMARY KEY (id),
  UNIQUE KEY uq_users_email (email),
  UNIQUE KEY uq_users_firebase_uid (firebase_uid),
  KEY idx_users_firebase_uid (firebase_uid),
  KEY idx_users_provider_id (provider_id),
  KEY idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE public_urls (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  url_key VARCHAR(64) NOT NULL,
  is_active TINYINT(1) NOT NULL DEFAULT 1,
  created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  UNIQUE KEY idx_public_urls_url_key (url_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
