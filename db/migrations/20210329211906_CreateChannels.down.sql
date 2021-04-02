alter table recordings
    drop foreign key recordings_channel_id_foreign,
    drop column channel_id;

drop table if exists epgstation_channels;

drop table if exists channels;
