create table playlist(
id INT NOT NULL AUTO_INCREMENT,
videoid varchar(100),
name varchar(255),
length INT NOT NULL,
seek INT NOT NULL DEFAULT 0,
PRIMARY KEY(id)
);

alter table playlist add column added_by varchar(255) not null default "system";

alter table playlist add column thumbnail varchar(255);

create table play_history(
videoid varchar(100) PRIMARY KEY,
name varchar(255),
playcount int,
last_played datetime
);

create table user_history(
user varchar(255),
videoid varchar(100),
last_played datetime,
PRIMARY KEY(user,videoid)
);

create table song_details(
name varchar(255) not null default '',
duration int not null,
thumbnail varchar(255) not null default 'NA',
views bigint not null,
likes int not null,
dislikes int not null,
favourites int not null,
comments int not null
);
