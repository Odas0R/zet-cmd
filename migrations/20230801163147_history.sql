-- +goose Up
-- +goose StatementBegin
create table history (
    zettel_id text not null,
    updated_at text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')), -- use ISO8601/RFC3339
    created_at text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')), -- use ISO8601/RFC3339

    primary key (zettel_id),
    foreign key (zettel_id) references zettel(id) on delete cascade
) strict;

create index history_created_idx on history (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table history;
-- +goose StatementEnd
