BEGIN;
-- QR Mapping insert statements
\i qr_mapping/1.0.0/insert_file_storage_record.sql;
\i qr_mapping/1.0.0/insert_qr_mapping.sql;
-- User insert statements
\i user/1.0.0/insert_staff.sql;
COMMIT;