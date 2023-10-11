-- +goose Up
-- +goose StatementBegin
CREATE TABLE messages (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    sender UUID REFERENCES users(id),
    recipient UUID REFERENCES users(id),
    content TEXT,
    datetime TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE messages;
-- +goose StatementEnd
