-- +goose Up
-- +goose StatementBegin
create table zettel (
    id integer primary key autoincrement,
    title text not null,
    path text not null,
    type text not null,
    content text not null,
    created_at text not null,
    updated_at text not null
) strict;

create index zettel_created_idx on zettel (created_at);

create trigger zettel_updated_timestamp after update on zettel begin
  -- use ISO8601/RFC3339
  update zettel set updated_at = strftime('%Y-%m-%dT%H:%M:%fZ') where id = old.id;
end;

create table link (
    id integer primary key autoincrement,
    zettel_id integer not null,
    link_id integer not null,
    created_at text not null,
    updated_at text not null,

    foreign key (zettel_id) references zettel(id),
    foreign key (link_id) references zettel(id)
) strict;

create index link_created_idx on link (created_at);

create trigger link_updated_timestamp after update on link begin
  -- use ISO8601/RFC3339
  update link set updated_at = strftime('%Y-%m-%dT%H:%M:%fZ') where id = old.id;
end;
-- +goose StatementEnd
