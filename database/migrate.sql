BEGIN;
-- Order Communication Schema
\i order_communication/schema.sql;

-- Order Communication v1.0.0
\i order_communication/1.0.0/create_order.sql;
\i order_communication/1.0.0/create_pdf_generation_worker.sql;

-- QR Mapping Schema
\i qr_mapping/schema.sql

-- QR Mapping v1.0.0
\i qr_mapping/1.0.0/create_file_storage_record.sql;
\i qr_mapping/1.0.0/create_qr_mapping.sql;
\i qr_mapping/1.0.0/create_qr_mapping_impression.sql;
\i qr_mapping/1.0.0/create_impression.sql;

-- Twilio Schema
\i twilio/schema.sql;

-- Twilio v1.0.0
\i twilio/1.0.0/create_account.sql;
\i twilio/1.0.0/create_conversation_service.sql;
\i twilio/1.0.0/create_incoming_phone_number_local.sql;
\i twilio/1.0.0/create_role.sql;
\i twilio/1.0.0/create_user.sql;
\i twilio/1.0.0/create_conversation.sql;
\i twilio/1.0.0/create_conversation_participant.sql;
\i twilio/1.0.0/create_conversation_message.sql;
\i twilio/1.0.0/create_media.sql;
\i twilio/1.0.0/create_service.sql;

-- User schema
\i user/1.0.0/create_staff.sql;

-- Google Workspace Forms Schema
\i google_workspace_forms/schema.sql

-- Google Workspace Forms v1.0.0
\i google_workspace_forms/1.0.0/create_answer.sql;
\i google_workspace_forms/1.0.0/create_form_response.sql;
\i google_workspace_forms/1.0.0/create_form.sql;
\i google_workspace_forms/1.0.0/create_item.sql;

COMMIT;