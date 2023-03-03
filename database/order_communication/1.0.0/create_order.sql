create table if not exists order_communication.order (
    order_id serial primary key not null,
    parent_order_id int,
    pdf_generation_worker_id int,
    order_file_storage_url text,
    is_pdf_generated boolean not null default false,
    is_viewed_by_manager boolean not null,
    is_accepted_by_manager boolean not null,
    is_delivered_to_kitchen boolean not null,
    conversation_sid varchar(34),
    first_conversation_message_sid varchar(34),
    incoming_phone_number varchar(20),
    form_id text,
    form_response_id text,
    media_sid text,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

select setval('order_communication.order_order_id_seq', 1);