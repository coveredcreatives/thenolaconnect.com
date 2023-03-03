create table if not exists qr_mapping.impression (
    impression_id serial primary key not null,
    path text not null,
    ipaddress text not null,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

select setval('qr_mapping.impression_impression_id_seq', 1);