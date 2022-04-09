create table if not exists users.auca_departments (
    id smallserial primary key,
    name varchar(50) not null
);

create table if not exists users.user_data (
    user_id uuid not null primary key,
    first_name varchar(60),
    last_name varchar(100),
    avatar_url varchar(300),
    department_id smallint,
    year_of_admission smallint,
    birth_date date,
    updated_at timestamptz not null default('now'),

    constraint fk_user_data_users
        foreign key (user_id)
            references users.users(id),
    constraint fk_user_data_auca_departments
        foreign key (department_id)
            references users.auca_departments(id),
    constraint ch_admission_year
        check ((year_of_admission > 1992 and year_of_admission < date_part('year', date('now')) + 1) or year_of_admission is null)
);
