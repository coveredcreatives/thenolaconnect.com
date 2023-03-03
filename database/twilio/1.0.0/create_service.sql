create table if not exists twilio.service (
  sid text primary key,
  friendly_name text,
  account_sid text,
  date_created timestamp with time zone,
  date_updated timestamp with time zone
);