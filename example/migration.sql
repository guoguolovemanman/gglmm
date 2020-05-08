
drop table if exists example;

create table if not exists example (
  id bigint unsigned not null auto_increment,
  created_at timestamp null default null,
  updated_at timestamp null default null,
  deleted_at timestamp null default null,
  int_value int not null,
  float_value float(10, 2) not null,
  string_value varchar(255) not null,
  primary key (id)
);