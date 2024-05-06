-- +goose Up
-- +goose StatementBegin
drop trigger link_updated_timestamp;
create trigger link_updated_timestamp after update on link begin
  -- use ISO8601/RFC3339
  update link set updated_at = strftime('%Y-%m-%dT%H:%M:%fZ') where zettel_id = old.zettel_id and link_id = old.link_id;
end;

alter table zettel add column slug text;
UPDATE zettel
SET slug = CASE
WHEN path LIKE '/home/odas0r/github.com/odas0r/zet/fleet/%' THEN
substr(substr(path, length('/home/odas0r/github.com/odas0r/zet/fleet/') + 1), 1, instr(substr(path, length('/home/odas0r/github.com/odas0r/zet/fleet/') + 1), '.') - 1)
WHEN path LIKE '/home/odas0r/github.com/odas0r/zet/permanent/%' THEN
substr(substr(path, length('/home/odas0r/github.com/odas0r/zet/permanent/') + 1), 1, instr(substr(path, length('/home/odas0r/github.com/odas0r/zet/permanent/') + 1), '.') - 1)
ELSE
substr(substr(path, length('/home/odas0r/github.com/odas0r/zet/') + 1), 1, instr(substr(path, length('/home/odas0r/github.com/odas0r/zet/') + 1), '.') - 1)
END
WHERE path like '/home/odas0r/github.com/odas0r/zet/%';

-- create an index on the slug
create unique index idx_unique_slug on zettel (slug, id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- drop index
drop index idx_unique_slug;
alter table zettel drop column slug;
-- +goose StatementEnd
