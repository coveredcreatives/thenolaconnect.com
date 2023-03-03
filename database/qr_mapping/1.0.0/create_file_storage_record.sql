create table if not exists qr_mapping.file_storage_record (
    file_storage_record_id serial primary key not null,
    file_storage_url text,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

select setval('qr_mapping.file_storage_record_file_storage_record_id_seq', 1);