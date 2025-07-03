-- +goose Up
-- admin/admin
INSERT INTO `user` (username, password) VALUES ('admin', '$2a$10$kSMPeUITJB9tZ2bs.Wm5gukp0Vv9qP3P.K.To.S2OrqAQv0mhzgtS');

-- +goose Down
DELETE FROM `user` WHERE username='admin';