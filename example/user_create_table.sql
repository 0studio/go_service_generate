create table if not exists `User`(
`id` bigint NOT NULL DEFAULT 0,
`helloName` varchar(10) NOT NULL DEFAULT '',
`age` int NOT NULL DEFAULT 0,
`sex` tinyint NOT NULL DEFAULT 0,
`t` int NOT NULL DEFAULT 0,
`t2` timestamp NOT NULL DEFAULT 0
,primary key (id,helloName)
);
