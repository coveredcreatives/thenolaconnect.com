create table if not exists twilio.user (
    account_sid varchar(34),
    conversation_service_sid varchar(34),
    sid varchar(34) primary key not null,
    role_sid varchar(34),
    identity text,
    friendly_name text not null,
    attributes text,
    is_online boolean,
    url text not null,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);