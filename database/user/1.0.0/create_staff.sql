create table if not exists public.staff (
  staff_id serial primary key not null,
  email varchar(64) unique not null,
  phone_number varchar(64) unique,
  name text not null,
  password text,
  email_verification_code varchar(6),
  email_verification_completed boolean,
  two_factor_auth_code varchar(6),
  two_factor_auth_completed boolean,
  is_manager boolean,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone
)