package models

import "time"

type RegisterStudentPayload struct {
	NewStudent			Students			`json:"student_data" binding:"required"`
	StudentAcademicData	[]DataAkademik	`json:"student_academic_data" binding:"required"`
}

type Students struct {
	ID        			uint    `gorm:"primaryKey"`
	Nama      			string  `json:"nama" binding:"required"`
	TTL       			string  `json:"ttl" binding:"required"`
	Alamat    			string  `json:"alamat" binding:"required"`
	Kota	  			string	`json:"kota" binding:"required"`
	Provinsi  			string  `json:"provinsi" binding:"required"`
	Email     			string  `json:"email" binding:"required,email"`
	KodePos	  			uint	`json:"kode_pos" binding:"required"`
	NomorHP	  			string	`json:"nomor_hp" binding:"required"`
	Warganegara			string	`json:"warganegara" binding:"required"`	
	NoAnggota			string	`json:"no_anggota" binding:"required"`
	StatusMahasiswa		bool	`json:"status_mahasiswa" binding:"required"`
	LembagaPendidikan	string	`json:"lembaga_pendidikan" binding:"required"`
	Jurusan				string	`json:"jurusan" binding:"required"`
	Fakultas			string	`json:"fakultas" binding:"required"`
	BidangUsaha			string	`json:"bidang_usaha" binding:"required"`
	AlamatUsaha			string	`json:"alamat_usaha" binding:"required"`
	FotoPath  			string  `json:"foto_path" binding:"required"`
	Createdat 			time.Time	`gorm:"autoCreateTime"`
	DataAkademik 		[]DataAkademik `gorm:"foreignKey:StudentID"`
}

type StudentUpdate struct {
	Nama      			*string `json:"nama,omitempty"`
	TTL       			*string `json:"ttl,omitempty"`
	Alamat    			*string `json:"alamat,omitempty"`
	Kota	  			*string	`json:"kota,omitempty"`
	Provinsi  			*string `json:"provinsi,omitempty"`
	Email     			*string `json:"email,omitempty"`
	KodePos	  			*uint	`json:"kode_pos,omitempty"`
	NomorHP	  			*string	`json:"nomor_hp,omitempty"`
	Warganegara			*string	`json:"warganegara,omitempty"`	
	NoAnggota			*string	`json:"no_anggota,omitempty"`
	StatusMahasiswa		*bool	`json:"status_mahasiswa,omitempty"`
	LembagaPendidikan	*string	`json:"lembaga_pendidikan,omitempty"`
	Jurusan				*string	`json:"jurusan,omitempty"`
	Fakultas			*string	`json:"fakultas,omitempty"`
	BidangUsaha			*string	`json:"bidang_usaha,omitempty"`
	AlamatUsaha			*string	`json:"alamat_usaha,omitempty"`
	FotoPath  			*string `json:"foto_path,omitempty"`
}

type DataAkademik struct {
	ID        	uint    	`gorm:"primaryKey"`
	StudentID	uint		`gorm:"index"`
	Student   	Students 	`gorm:"foreignKey:StudentID;references:ID"`
	ContentLink	string		`json:"content_link" binding:"required"`
	Tipe		string		`json:"tipe" binding:"required"`
	Createdat 	time.Time	`gorm:"autoCreateTime"`
}
