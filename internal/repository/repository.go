package repository

import (
	"compare/internal/models"
	"compare/pkg/logging"
	"gorm.io/gorm"
)

var logger = logging.GetLogger()

type Repository struct {
	Conn *gorm.DB
}

func NewRepository(conn *gorm.DB) *Repository {
	return &Repository{Conn: conn}
}

// ============================================>

func (r Repository) HumoPayment(hp []models.HumoPayment) error {
	err := r.Conn.Create(hp).Error
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (r Repository) PartnerPayment(pp []models.PartnerPayment) error {
	err := r.Conn.Create(pp).Error
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (r Repository) GetDataFromHumo() (humoPayments []models.HumoPayment, err error) {
	if err = r.Conn.Find(&humoPayments).Error; err != nil {
		logger.Error(err)
		return
	}

	return
}

func (r Repository) GetAllPartners() (partners []models.Partner, err error) {
	if err = r.Conn.Find(&partners).Error; err != nil {
		logger.Errorf("can not get all partners from DB: %v", err)
		return
	}

	return
}

//func (r *Repository) NewReestrRecord(partner models.Partner, date *time.Time) error {
//	var newRecord models.Reestr
//	dateFormat := date.Format("02.01.2006")
//
//	newRecord.FileName = fmt.Sprintf("input_%v.csv", dateFormat)
//	newRecord.DateStart = date.Add(24*time.Hour + time.Duration(partner.IntervalStart)*time.Hour)
//	newRecord.DateStop = date.Add(24*time.Hour + time.Duration(partner.IntervalStop)*time.Hour)
//
//	if partner.IntervalDate == "d" {
//		newRecord.DateStart = date.Add(24 * time.Hour)
//		newRecord.DateStop = date.Add(24*time.Hour + 5 + time.Minute)
//	}
//
//	newRecord.PartnerId = partner.Id
//
//	if err := r.Conn.Create(&newRecord).Error; err != nil {
//		return err
//	}
//
//	return nil
//}

func (r Repository) NewReestrRecord(reestr []models.Reestr) error {
	if err := r.Conn.Create(&reestr).Error; err != nil {
		logger.Error("can not create New Reestr", err)
		return err
	}

	return nil
}
