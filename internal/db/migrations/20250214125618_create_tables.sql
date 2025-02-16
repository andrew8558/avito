-- +goose Up
-- +goose StatementBegin
CREATE TABLE users(
    login varchar(256) primary key not null,
    password varchar(256) not null,
    coins integer not null
);
CREATE TABLE purchases(
    login varchar(256) not null,
    type varchar(256) not null,
    quantity integer not null,
    primary key (login, type)
);
CREATE TABLE send_coin_events(
    id serial primary key not null,
    to_user varchar(256) not null,
    from_user varchar(256) not null,
    amount integer not null
);
CREATE INDEX idx_send_coin_events_to_user ON send_coin_events(to_user);
CREATE INDEX idx_send_coin_events_from_user ON send_coin_events(from_user);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE purchases;
DROP TABLE send_coin_events;
-- +goose StatementEnd
