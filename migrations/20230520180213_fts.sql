-- +goose Up
-- +goose StatementBegin
create virtual table zettel_fts
  using fts5(id, title, content, path, tokenize = porter, content = 'zettel', content_rowid = 'id');

create trigger zettel_after_insert after insert on zettel begin
  insert into zettel_fts(rowid, id, title, content, path)
    values (new.id, new.id, new.title, new.content, new.path);
end;

create trigger zettel_after_update after update on zettel begin
  insert into zettel_fts(zettel_fts, rowid)
    values('delete', old.id);
  insert into zettel_fts(rowid, id, title, content, path)
    values (new.id, new.id, new.title, new.content, new.path);
end;

create trigger zettel_after_delete after delete on zettel begin
  insert into zettel_fts(zettel_fts, rowid)
    values('delete', old.id);
end;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table zettel_fts;
drop trigger zettel_after_insert;
drop trigger zettel_after_update;
drop trigger zettel_after_delete;
-- +goose StatementEnd
