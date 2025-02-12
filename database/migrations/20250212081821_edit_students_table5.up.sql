ALTER TABLE students 
ADD CONSTRAINT unique_nama_email UNIQUE (nama, email),
ALTER COLUMN nomor_hp SET DATA TYPE TEXT USING nomor_hp::TEXT;