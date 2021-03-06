-- +goose Up
CREATE TABLE IF NOT EXISTS devices (
    id          serial      PRIMARY KEY,
    platform    varchar(32) NOT NULL,
    user_id     bigint      NOT NULL,
    entered_at  timestamp   NOT NULL,
    removed     bool        DEFAULT false NOT NULL,
    created_at  timestamp   DEFAULT now() NOT NULL,
    updated_at  timestamp   DEFAULT now() NOT NULL
);

INSERT INTO devices (platform, user_id, entered_at, removed, created_at, updated_at)
VALUES  ('Ios', 24154, '2021-11-09 13:51:49.870041', false, '2021-11-09 13:51:49.870416', '2021-11-09 13:51:49.870416'),
        ('Android', 28123, '2021-11-09 13:52:16.785666', false, '2021-11-09 13:52:16.785891', '2021-11-09 13:52:16.785891'),
        ('Linux', 41244, '2021-11-09 13:52:26.945600', false, '2021-11-09 13:52:26.945798', '2021-11-09 13:52:26.945798'),
        ('Android', 412414, '2021-11-09 13:52:20.094174', true, '2021-11-09 13:52:20.094411', '2021-11-09 13:52:20.094411'),
        ('Windows', 52352, '2021-11-09 13:53:07.152502', false, '2021-11-09 13:53:07.152723', '2021-11-09 13:53:07.152723'),
        ('Windows', 76262, '2021-11-09 13:52:45.590365', true, '2021-11-09 13:52:45.590560', '2021-11-09 13:52:45.590560'),
        ('Ios', 241442, '2021-11-09 13:53:26.350223', false, '2021-11-09 13:53:26.350434', '2021-11-09 13:53:26.350434'),
        ('Ios', 15515, '2021-11-09 13:53:29.226157', false, '2021-11-09 13:53:29.226243', '2021-11-09 13:53:29.226243'),
        ('Android', 74349, '2021-11-09 13:53:39.733289', false, '2021-11-09 13:53:39.733535', '2021-11-09 13:53:39.733535'),
        ('Linux', 124141, '2021-11-09 13:54:04.246343', false, '2021-11-09 13:54:04.246528', '2021-11-09 13:54:04.246528');

-- +goose Down
DROP TABLE IF EXISTS devices;
