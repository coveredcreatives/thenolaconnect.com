create table if not exists google_workspace_forms.form (
  form_id text,
  description text,
  document_title text,
  title text,
  linked_sheet_id text,
  responder_uri text,
  revision_id text,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone
);