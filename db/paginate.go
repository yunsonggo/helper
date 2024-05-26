package db

import "gorm.io/gorm"

func Paginate(page, pageSize, max int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if max <= 0 {
			max = 100
		}
		switch {
		case pageSize > max:
			pageSize = max
		case pageSize <= 0:
			pageSize = 1
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
