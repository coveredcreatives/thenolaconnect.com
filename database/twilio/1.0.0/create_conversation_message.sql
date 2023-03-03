create table if not exists twilio.conversation_message (
    account_sid varchar(34),
    conversation_sid varchar(34),
    sid varchar(34) primary key not null,
    index integer not null,
    author  text not null,
    body text,
    media text,
    attributes text,
    participant_sid varchar(34),
    url text not null,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);