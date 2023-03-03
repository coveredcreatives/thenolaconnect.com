create table if not exists twilio.conversation (
    account_sid varchar(34),
    chat_service_sid varchar(34),
    messaging_service_sid varchar(34),
    sid varchar(34) primary key not null,
    friendly_name text not null,
    unique_name text not null,
    attributes text,
    state varchar(10) not null,
    url text not null,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);