-- ============================================================
-- 経費精算システム 初期スキーマ
-- ============================================================

-- UUID生成用拡張
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================
-- 1. ユーザーマスタ
-- ============================================================
CREATE TABLE users (
    id            UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    name          VARCHAR(100) NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================
-- 2. ユーザーロール紐付け（複合主キー）
-- ============================================================
CREATE TABLE user_roles (
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role       VARCHAR(50) NOT NULL,
    created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role)
);

-- ============================================================
-- 3. 勘定科目マスタ
-- ============================================================
CREATE TABLE categories (
    id         UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    code       VARCHAR(50)  NOT NULL UNIQUE,
    name       VARCHAR(100) NOT NULL,
    sort_order INTEGER      NOT NULL DEFAULT 0,
    is_active  BOOLEAN      NOT NULL DEFAULT true
);

-- ============================================================
-- 4. 経費データ
-- ============================================================
CREATE TABLE expenses (
    id                 UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id            UUID         NOT NULL REFERENCES users(id),
    expense_date       DATE         NOT NULL,
    category_id        UUID         NOT NULL REFERENCES categories(id),
    amount             INTEGER      NOT NULL,
    description        VARCHAR(500) NOT NULL,
    receipt_image_path VARCHAR(500),
    status             VARCHAR(50)  NOT NULL DEFAULT 'draft',
    submitted_at       TIMESTAMP,
    created_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================
-- 5. 承認履歴
-- ============================================================
CREATE TABLE approval_histories (
    id         UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    expense_id UUID        NOT NULL REFERENCES expenses(id) ON DELETE CASCADE,
    action     VARCHAR(50) NOT NULL,
    actor_id   UUID        NOT NULL REFERENCES users(id),
    actor_role VARCHAR(50) NOT NULL,
    comment    TEXT,
    created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================
-- 6. インデックス
-- ============================================================
CREATE INDEX idx_expenses_user_id ON expenses(user_id);
CREATE INDEX idx_expenses_status ON expenses(status);
CREATE INDEX idx_approval_histories_expense_id ON approval_histories(expense_id);

-- ============================================================
-- 7. 初期マスタデータ（勘定科目）
-- ============================================================
INSERT INTO categories (code, name, sort_order) VALUES
    ('transportation', '交通費',   1),
    ('entertainment',  '交際費',   2),
    ('supplies',       '消耗品費', 3),
    ('communication',  '通信費',   4),
    ('other',          'その他',   5);
