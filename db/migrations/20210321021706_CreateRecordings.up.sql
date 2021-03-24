create table recordings
(
    id           bigint auto_increment not null unique primary key,
    program_name varchar(255)          not null,
    started_at   datetime              not null,
    ends_at      datetime              not null,
    ended        boolean               not null default false
);

create table epgstation_recordings
(
    epgstation_id bigint not null unique primary key,
    recording_id  bigint not null unique,
    foreign key epgstation_recordings_recordings_id_foreign (recording_id) references recordings (id)
);

create table discord_recordings
(
    recording_id bigint       not null unique primary key,
    foreign key discord_recordings_recordings_id_foreign (recording_id) references recordings (id),
    channel_id   varchar(255) not null,
    message_id   varchar(255) not null
);
