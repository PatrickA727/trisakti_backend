ALTER TABLE students
DROP COLUMN ttl;

ALTER TABLE students
ADD COLUMN tempat_lahir VARCHAR(255) NOT NULL DEFAULT 'empty';

ALTER TABLE students
ADD COLUMN tanggal_lahir DATE NOT NULL DEFAULT '2000-01-01';
