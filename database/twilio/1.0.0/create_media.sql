create table if not exists twilio.media (
  sid text primary key,
  service_sid text,
  date_created timestamp with time zone,
  date_upload_updated timestamp with time zone,
  date_updated timestamp with time zone,
  links_content text,
  size integer,
  content_type text,
  filename text,
  author text,
  category text,
  message_sid text,
  channel_sid text,
  url text
);