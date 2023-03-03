create table if not exists qr_mapping.qr_mapping (
    file_storage_record_id integer references qr_mapping.file_storage_record not null,
    qr_file_storage_url text not null,
    name text not null,
    qr_encoded_data text primary key not null,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,
    deleted_at timestamp with time zone
);