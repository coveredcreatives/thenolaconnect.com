create table if not exists google_workspace_forms.answer (
  question_id text,
  form_id text,
  form_response_id text,
  text_answers text,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone
)