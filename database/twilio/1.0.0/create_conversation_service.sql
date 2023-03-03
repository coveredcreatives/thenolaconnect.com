create table if not exists twilio.conversation_service (
    account_sid varchar(34),
    sid varchar(34) primary key not null,
    friendly_name text not null,
    date_created timestamp with time zone,
    date_updated timestamp with time zone,
    links text,
    url text not null,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);