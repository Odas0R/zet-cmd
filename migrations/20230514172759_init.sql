-- +goose Up
-- +goose StatementBegin
create table zettel (
    id text not null primary key,
    title text not null,
    path text not null,
    type text not null,
    content text not null,
    created_at text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),
    updated_at text not null default (strftime('%Y-%m-%dT%H:%M:%fZ'))
) strict;

create index zettel_created_idx on zettel (created_at);
create index zettel_path_idx on zettel (path);

create trigger zettel_updated_timestamp after update on zettel begin
  -- use ISO8601/RFC3339
  update zettel set updated_at = strftime('%Y-%m-%dT%H:%M:%fZ') where id = old.id;
end;

create table link (
    zettel_id text not null,
    link_id text not null,
    created_at text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),
    updated_at text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),

    primary key (zettel_id, link_id),

    foreign key (zettel_id) references zettel(id) on delete cascade,
    foreign key (link_id) references zettel(id) on delete cascade
) strict;

create index link_created_idx on link (created_at);

create trigger link_updated_timestamp after update on link begin
  -- use ISO8601/RFC3339
  update link set updated_at = strftime('%Y-%m-%dT%H:%M:%fZ') where id = old.id;
end;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table link;
drop table zettel;
-- +goose StatementEnd
