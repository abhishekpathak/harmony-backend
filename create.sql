create table playlist(
id INT NOT NULL AUTO_INCREMENT,
videoid varchar(100),
name varchar(255),
length INT NOT NULL,
seek INT NOT NULL DEFAULT 0,
PRIMARY KEY(id)
);