create table if not exists order_communication.pdf_generation_worker (
  pdf_generation_worker_id serial primary key not null,
  form_id text,
  order_id int references order_communication.order not null,
  process_id int,
  start_at timestamp with time zone,
  completed_at timestamp with time zone,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone
);

select setval('order_communication.pdf_generation_worker_pdf_generation_worker_id_seq', 1);