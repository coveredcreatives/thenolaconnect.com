create table if not exists twilio.account (
    sid varchar(34) primary key not null,
    owner_account_sid varchar(34) not null,
    auth_token text,
    friendly_name varchar(64),
    status varchar(10) not null,
    type varchar(10) not null,
    uri text not null,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);