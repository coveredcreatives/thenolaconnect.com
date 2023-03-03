create table if not exists google_workspace_forms.item (
  description text,
  item_id text,
  title text,
  is_question_item boolean,
  question_id text,
  is_choice_question boolean,
  is_date_question boolean,
  is_file_upload_question boolean,
  is_required boolean,
  is_row_question boolean,
  is_scale_question boolean,
  is_text_question boolean,
  is_time_question boolean,
  form_id text,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone
);