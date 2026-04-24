-- ============================================================
-- テストデータ投入（開発・UAT環境用）
-- test-data.md に基づくシードデータ
-- ============================================================

-- テストユーザー（パスワードは全員 "password123"）
-- bcrypt hash of "password123"
INSERT INTO users (id, name, email, password_hash) VALUES
    ('11111111-1111-1111-1111-111111111111', '一般 太郎', 'taro.ippan@example.com', '$2a$10$3OvQvEL7KtZdS/banz6vwuu4wP8Y3cRUyeeX36ntJ/uKgNj/o3q96'),
    ('22222222-2222-2222-2222-222222222222', '一般 花子', 'hanako.ippan@example.com', '$2a$10$3OvQvEL7KtZdS/banz6vwuu4wP8Y3cRUyeeX36ntJ/uKgNj/o3q96'),
    ('33333333-3333-3333-3333-333333333333', '上長 一郎', 'ichiro.manager@example.com', '$2a$10$3OvQvEL7KtZdS/banz6vwuu4wP8Y3cRUyeeX36ntJ/uKgNj/o3q96'),
    ('44444444-4444-4444-4444-444444444444', '経理 担当', 'tanto.keiri@example.com', '$2a$10$3OvQvEL7KtZdS/banz6vwuu4wP8Y3cRUyeeX36ntJ/uKgNj/o3q96'),
    ('55555555-5555-5555-5555-555555555555', '経理 部長', 'bucho.keiri@example.com', '$2a$10$3OvQvEL7KtZdS/banz6vwuu4wP8Y3cRUyeeX36ntJ/uKgNj/o3q96')
ON CONFLICT (email) DO NOTHING;

-- ユーザーロール
INSERT INTO user_roles (user_id, role) VALUES
    ('11111111-1111-1111-1111-111111111111', 'applicant'),
    ('22222222-2222-2222-2222-222222222222', 'applicant'),
    ('33333333-3333-3333-3333-333333333333', 'manager'),
    ('44444444-4444-4444-4444-444444444444', 'expense_admin'),
    ('55555555-5555-5555-5555-555555555555', 'finance_director')
ON CONFLICT (user_id, role) DO NOTHING;

-- 勘定科目（test-data.md の拡張版、001_initial.sql と重複しないようUPSERT）
INSERT INTO categories (code, name, sort_order, is_active) VALUES
    ('travel', '旅費', 5, true),
    ('books', '書籍・研修費', 6, true),
    ('disabled_item', '廃止された科目', 100, false)
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    sort_order = EXCLUDED.sort_order,
    is_active = EXCLUDED.is_active;

-- sort_order を test-data.md の値に合わせて更新
UPDATE categories SET sort_order = 10 WHERE code = 'transportation';
UPDATE categories SET sort_order = 20 WHERE code = 'entertainment';
UPDATE categories SET sort_order = 30 WHERE code = 'supplies';
UPDATE categories SET sort_order = 40 WHERE code = 'communication';
UPDATE categories SET sort_order = 50 WHERE code = 'travel';
UPDATE categories SET sort_order = 60 WHERE code = 'books';
UPDATE categories SET sort_order = 99 WHERE code = 'other';
