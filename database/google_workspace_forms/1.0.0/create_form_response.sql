create table if not exists google_workspace_forms.form_response (
    create_time timestamp with time zone,
    form_id text,
    last_submitted_time timestamp with time zone,
    respondent_email text,
    response_id text,
    total_score int,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);