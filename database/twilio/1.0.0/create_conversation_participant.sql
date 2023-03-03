create table if not exists twilio.conversation_participant (
    account_sid varchar(34),
    conversation_sid varchar(34),
    sid varchar(34) primary key not null,
    identity text,
    attributes text,
    messaging_binding text,
    messaging_binding_address text,
    messaging_binding_proxy_address text,
    role_sid varchar(34),
    date_created timestamp with time zone,
    date_updated timestamp with time zone,
    url text not null,
    last_read_message_index int,
    last_read_timestamp text,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);