create table if not exists twilio.role (
    account_sid varchar(34),
    conversation_service_sid varchar(34),
    sid varchar(34) primary key not null,
    identity text,
    friendly_name text not null,
    type varchar(20) not null,
    permissions text,
    url text not null,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);