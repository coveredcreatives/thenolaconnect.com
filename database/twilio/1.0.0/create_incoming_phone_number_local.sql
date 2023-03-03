create table if not exists twilio.incoming_phone_number_local (
  account_sid varchar(34),
  address_sid varchar(34),
  address_requirements text,
  api_version varchar(10),
  beta boolean,
  friendly_name text,
  identity_sid varchar(34),
  phone_number varchar(20),
  origin text,
  sid varchar(34),
  sms_application_sid varchar(34),
  sms_fallback_method text,
  sms_fallback_url text,
  sms_method text,
  sms_url text,
  status_callback text,
  status_callback_method text,
  trunk_sid varchar(34),
  uri text,
  voice_receive_mode text,
  voice_application_sid varchar(34),
  voice_callier_id_lookup boolean,
  voice_fallback_method text,
  voice_fallback_url text,
  voice_method text,
  voice_url text,
  emergency_status text,
  emergency_address_sid varchar(34),
  emergency_address_status text,
  bundle_sid varchar(34),
  status text
)