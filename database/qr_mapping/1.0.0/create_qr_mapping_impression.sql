create table if not exists qr_mapping.qr_mapping_impression (
    qr_mapping_impression_id serial primary key not null,
    qr_encoded_data text references qr_mapping.qr_mapping not null,
    path text not null,
    ipaddress text not null,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,
    deleted_at timestamp with time zone
);

select setval('qr_mapping.qr_mapping_impression_qr_mapping_impression_id_seq', 1);