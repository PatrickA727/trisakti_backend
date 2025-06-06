package models

import "time"

type RegisterStudentPayload struct {
	NewStudent			StudentRegister			`json:"student_data" binding:"required"`
	StudentAcademicData	[]DataAkademik	`json:"student_academic_data"`
}

type EditStudentPayload struct {	// Unused
	StudentData	Students	`json:"student_data,omitempty"`
	StudentAcad	[]DataAkademik	`json:"student_academic_data,omitempty"`
}

type Students struct {
	ID        			uint    `gorm:"primaryKey"`
	Nama      			string  `json:"nama" binding:"required"`
	Alamat    			string  `json:"alamat" binding:"required"`
	Email     			string  `json:"email" binding:"required,email"`
	NomorHP	  			string	`json:"nomor_hp" binding:"required"`
	NoAnggota			string	`json:"no_anggota" binding:"required"`
	StatusMahasiswa		bool	`json:"status_mahasiswa" binding:"required"`
	// LembagaPendidikan	string	`json:"lembaga_pendidikan" binding:"required"`
	// Jurusan				string	`json:"jurusan" binding:"required"`
	NamaJurusan			string	`json:"nama_jurusan"`
	BidangUsaha			string	`json:"bidang_usaha"`
	AlamatUsaha			string	`json:"alamat_usaha"`
	JenisKelamin		string	`json:"jenis_kelamin" binding:"required"`
	TahunMasuk			string	`json:"tahun_masuk" binding:"required"`
	Semester			string	`json:"semester" binding:"required"`
	Ktp					string	`json:"ktp" binding:"required"`		// Alamat by KTP
	AsalDaerah			string	`json:"asal_daerah" binding:"required"`
	Agama				string	`json:"agama"`
	AsalSekolah			string	`json:"asal_sekolah" binding:"required"`
	BahasaAsing			string	`json:"bahasa_asing"`
	Hobi				string	`json:"hobi"`
	MediaSosial			string	`json:"media_sosial"`
	Keterampilan		string	`json:"keterampilan"`
	NoTelpUsaha			string	`json:"no_telp_usaha"`
	TempatLahir			string	`json:"tempat_lahir"`
	TanggalLahir		string	`json:"tanggal_lahir"`
	LinkMedsos			string	`json:"link_medsos"`
	FotoPath  			string  `json:"foto_path" binding:"required"`
	Sarjana				string	`json:"sarjana"`
	Satuan				string	`json:"satuan"`
	SatuanFK			uint	`json:"satuan_fk"`
	JurusanFK			uint	`json:"jurusan_fk"`
	Createdat 			time.Time	`gorm:"autoCreateTime"`
	// SatuanPendidikan    SatuanPendidikan `gorm:"foreignKey:SatuanFK"`
	DataAkademik 		[]DataAkademik `gorm:"foreignKey:StudentID"`
}

type StudentRegister struct {
	ID        			uint    `gorm:"primaryKey"`
	Nama      			string  `json:"nama" binding:"required"`
	Alamat    			string  `json:"alamat" binding:"required"`
	Email     			string  `json:"email" binding:"required,email"`
	NomorHP	  			string	`json:"nomor_hp" binding:"required"`
	NoAnggota			string	`json:"no_anggota" binding:"required"`
	StatusMahasiswa		bool	`json:"status_mahasiswa" binding:"required"`
	// LembagaPendidikan	string	`json:"lembaga_pendidikan" binding:"required"`
	Jurusan				string	`json:"jurusan"`
	BidangUsaha			string	`json:"bidang_usaha"`
	AlamatUsaha			string	`json:"alamat_usaha"`
	JenisKelamin		string	`json:"jenis_kelamin" binding:"required"`
	TahunMasuk			string	`json:"tahun_masuk" binding:"required"`
	Semester			string	`json:"semester" binding:"required"`
	Ktp					string	`json:"ktp" binding:"required"`		// Alamat by KTP
	AsalDaerah			string	`json:"asal_daerah" binding:"required"`
	Agama				string	`json:"agama"`
	AsalSekolah			string	`json:"asal_sekolah" binding:"required"`
	BahasaAsing			string	`json:"bahasa_asing"`
	Hobi				string	`json:"hobi"`
	MediaSosial			string	`json:"media_sosial"`
	Keterampilan		string	`json:"keterampilan"`
	NoTelpUsaha			string	`json:"no_telp_usaha"`
	TempatLahir			string	`json:"tempat_lahir"`
	TanggalLahir		string	`json:"tanggal_lahir"`
	LinkMedsos			string	`json:"link_medsos"`
	FotoPath  			string  `json:"foto_path" binding:"required"`
	Sarjana				string	`json:"sarjana"`
	SatuanFK			uint	`json:"satuan_fk"`
	JurusanFK			uint	`json:"jurusan_fk"`
	Createdat 			time.Time	`gorm:"autoCreateTime"`
	// SatuanPendidikan    SatuanPendidikan `gorm:"foreignKey:SatuanFK"`
	DataAkademik 		[]DataAkademik `gorm:"foreignKey:StudentID"`
}

type StudentsPayload struct {
	ID        			*uint    `gorm:"primaryKey,omitempty"`
	Nama      			*string  `json:"nama,omitempty"`
	Alamat    			*string  `json:"alamat,omitempty"`
	Email     			*string  `json:"email,omitempty"`
	NomorHP	  			*string	`json:"nomor_hp,omitempty"`
	NoAnggota			*string	`json:"no_anggota,omitempty"`
	StatusMahasiswa		*bool	`json:"status_mahasiswa,omitempty"`
	// LembagaPendidikan	string	`json:"lembaga_pendidikan" binding:"required"`
	NamaJurusan				*string	`json:"nama_jurusan,omitempty"`
	BidangUsaha			*string	`json:"bidang_usaha,omitempty"`
	AlamatUsaha			*string	`json:"alamat_usaha,omitempty"`
	JenisKelamin		*string	`json:"jenis_kelamin,omitempty"`
	TahunMasuk			*string	`json:"tahun_masuk,omitempty"`
	Semester			*string	`json:"semester,omitempty"`
	Ktp					*string	`json:"ktp,omitempty"`		// Alamat by KTP
	AsalDaerah			*string	`json:"asal_daerah,omitempty"`
	Agama				*string	`json:"agama,omitempty"`
	AsalSekolah			*string	`json:"asal_sekolah,omitempty"`
	BahasaAsing			*string	`json:"bahasa_asing,omitempty"`
	Hobi				*string	`json:"hobi,omitempty"`
	MediaSosial			*string	`json:"media_sosial,omitempty"`
	Keterampilan		*string	`json:"keterampilan,omitempty"`
	NoTelpUsaha			*string	`json:"no_telp_usaha,omitempty"`
	TempatLahir			*string	`json:"tempat_lahir,omitempty"`
	TanggalLahir		*string	`json:"tanggal_lahir,omitempty"`
	LinkMedsos			*string	`json:"link_medsos,omitempty"`
	FotoPath  			*string  `json:"foto_path,omitempty"`
	Sarjana				*string	`json:"sarjana,omitempty"`
	Satuan				*string	`json:"satuan,omitempty"`	
}

type StudentUpdate struct {
	Nama      			*string  `json:"nama,omitempty"`	// omitempty to remove struct field if empty(no entry from request)
	Alamat    			*string  `json:"alamat,omitempty"`
	Email     			*string  `json:"email,omitempty"`
	NomorHP	  			*string	`json:"nomor_hp,omitempty"`
	NoAnggota			*string	`json:"no_anggota,omitempty"`
	StatusMahasiswa		*bool	`json:"status_mahasiswa,omitempty"`
	// LembagaPendidikan	*string	`json:"lembaga_pendidikan,omitempty"`
	// Jurusan				*string	`json:"jurusan,omitempty"`
	SatuanFK			*uint	`json:"satuan_fk,omitempty"`
	JurusanFK			*uint	`json:"jurusan_fk,omitempty"`
	BidangUsaha			*string	`json:"bidang_usaha,omitempty"`
	AlamatUsaha			*string	`json:"alamat_usaha,omitempty"`
	JenisKelamin		*string	`json:"jenis_kelamin,omitempty"`
	TahunMasuk			*string	`json:"tahun_masuk,omitempty"`
	Semester			*string	`json:"semester,omitempty"`
	AlamatKTP			*string	`json:"ktp,omitempty"`
	AsalDaerah			*string	`json:"asal_daerah,omitempty"`
	Agama				*string	`json:"agama,omitempty"`
	AsalSekolah			*string	`json:"asal_sekolah,omitempty"`
	BahasaAsing			*string	`json:"bahasa_asing,omitempty"`
	Hobi				*string	`json:"hobi,omitempty"`
	MediaSosial			*string	`json:"media_sosial,omitempty"`
	Keterampilan		*string	`json:"keterampilan,omitempty"`
	NoTelpUsaha			*string	`json:"no_telp_usaha,omitempty"`
	TempatLahir			*string	`json:"tempat_lahir,omitempty"`
	TanggalLahir		*string	`json:"tanggal_lahir,omitempty"`
	LinkMedsos			*string	`json:"link_medsos,omitempty"`
	FotoPath  			*string  `json:"foto_path,omitempty"`
}

type DataAkademik struct {
	ID        	uint    	`gorm:"primaryKey"`
	StudentID	uint		`json:"student_id"`
	// Student   	Students 	`gorm:"foreignKey:StudentID;references:ID"`
	NamaPrestasi	string	`json:"nama_prestasi"`
	ContentLink	string		`json:"content_link" binding:"required"`
	Tipe		string		`json:"tipe" binding:"required"`
	Createdat 	time.Time	`gorm:"autoCreateTime"`
}

type DataAkademikUpdate struct {
	StudentID	uint		`json:"student_id,omitempty"`
	NamaPrestasi	string	`json:"nama_prestasi,omitempty"`
	ContentLink	string		`json:"content_link,omitempty"`
	Tipe		string		`json:"tipe,omitempty"`
}

type SatuanPendidikan struct {
	ID		uint    	`gorm:"primaryKey"`
	Satuan	string		`json:"satuan"`
}

type SatuanPendidikanPayload struct {
	Satuan	string	`json:"satuan"`
}

type Jurusan struct {
	SatuanID	uint	`json:"satuan_id"`
	NamaJurusan	string	`json:"nama_jurusan"`
}

type JurusanGet struct {
	ID			uint	`json:"id"`
	NamaJurusan	string	`json:"nama_jurusan"`
}

type FilePayload struct {
	Extension		string	`json:"extension" binding:"required"`
	ContentType		string	`json:"content_type" binding:"required"`
}

type GetFilePayload struct {
	FileKey		string	`json:"file_key" binding:"required"`
}
