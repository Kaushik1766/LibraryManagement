create extension if not exists "uuid-ossp";

-- drop table if exists  transactions;

create table if not exists users (
    id uuid primary key default uuid_generate_v4(),
    name varchar(20) not null ,
    email varchar(20) unique not null ,
    password varchar(60) not null ,
    role int not null default 1
);

create table if not exists books(
    id uuid primary key default uuid_generate_v4(),
    title varchar(20) not null ,
    author varchar(20) not null
);

create table if not exists transactions(
    id uuid primary key default uuid_generate_v4(),
    book_id uuid references books(id) not null ,
    user_id uuid references users(id) not null ,
    issued_at timestamp default now() not null ,
    issued_till timestamp default now() + interval '1 day' not null ,
    returned_at timestamp default null,
    constraint check_issued_till check ( issued_at<issued_till )
);
