-- +goose Up
-- +goose StatementBegin
create virtual table zettel_fts using fts5(
  title,
  content,
  tokenize = 'porter',
  content = 'zettel',
  content_rowid = 'id'
);

create trigger zettel_fts_after_insert after insert on zettel begin
  insert into zettel_fts(rowid, title, content) values (new.id, new.title, new.content);
end;

create trigger zettel_fts_after_update after update on zettel begin
  update zettel_fts set title = new.title, content = new.content where rowid = new.id;
end;

create trigger zettel_fts_after_delete after delete on zettel begin
  insert into zettel_fts (zettel_fts, rowid, title, content) values('delete', old.id, old.title, old.content);
end;

-- +goose StatementEnd
