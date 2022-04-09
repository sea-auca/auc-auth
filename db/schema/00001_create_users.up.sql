create schema if not exists users;

create extension if not exists "uuid-ossp";

create table if not exists users.users (
    id uuid primary key default uuid_generate_v4(),
    email varchar(100) not null UNIQUE,
    is_active boolean not null default 'false',
    is_validated boolean not null default 'false'
);