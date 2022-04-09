
create table if not exists users.authentication_settings (
    user_id uuid not null,
    refresh_token_retention_in_hours int not null default(336),
    main_token_retention_in_minutes int not null default(5),
    enforce_two_factor_auth boolean not null default('false'),

    constraint fk_authentication_settings_users
        foreign key (user_id)
            references users.users(id)
);

create table if not exists users.credentials(
    user_id uuid not null,
    hash varchar(60),
    constraint fk_credentials_users
        foreign key (user_id)
            references users.users(id)
)