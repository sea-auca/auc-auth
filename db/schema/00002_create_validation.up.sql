create table if not exists users.validation_requests(
    id uuid primary key default uuid_generate_v4(),
    user_id uuid not null,
    was_utilised boolean not null default('false'),
    created_at timestamptz not null default('now'),
    updated_at timestamptz not null default('now'),

    constraint fk_validation_requests_users
        foreign key user_id
            references users.users(id) 
);