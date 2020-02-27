set timezone = 'Europe/Moscow';

create user dbuser with encrypted password 'En9NR2b869';

create database calendar owner dbuser;
grant all privileges on database calendar to dbuser;

\connect calendar
create table events (
    id UUID primary key,
    title text not null,
    description text,
    owner text not null,
    start_time timestamptz not null,
    duration interval not null,
    notify interval
);

create index start_time_idx on events using btree (start_time);
