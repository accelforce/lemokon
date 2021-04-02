create table channels
(
    id                    bigint auto_increment not null unique primary key,
    name                  varchar(255)          not null unique,
    remote_control_key_id integer default null
);

create table epgstation_channels
(
    epgstation_id bigint  not null unique primary key,
    channel_id    bigint  not null unique,
    constraint epgstation_channels_channel_id_foreign foreign key (channel_id) references channels (id),
    has_logo_data boolean not null default false
);

alter table recordings
    add column
        (
        channel_id bigint,
        constraint recordings_channel_id_foreign foreign key (channel_id) references channels (id)
        );
