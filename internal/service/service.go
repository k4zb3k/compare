package service

import (
	"compare/internal/models"
	"compare/internal/repository"
	"compare/pkg/logging"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var logger = logging.GetLogger()

type Service struct {
	Repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{Repo: repo}
}

// =================================================>

func (s *Service) HumoPayment(hp []models.HumoPayment) error {
	err := s.Repo.HumoPayment(hp)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (s *Service) PartnerPayment(pp []models.PartnerPayment) error {
	err := s.Repo.PartnerPayment(pp)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

const (
	BankALIF = "alif"
	BankIBT  = "ibt"
	BankDC   = "dc"
)

func (s *Service) SaveInput(pathDir string, pp []models.PartnerPayment) error {
	err := os.MkdirAll(pathDir, 0777)
	if err != nil {
		logger.Errorf("failed creating dir: %v", err)
		return err
	}

	date := time.Now().Format("02.01.2006")
	pathFile := fmt.Sprintf("%s/input_%v.csv", pathDir, date)

	_, err = os.Stat(pathFile)

	csvFile, err := os.Create(pathFile)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer func() {
		if err := csvFile.Close(); err != nil {
			logger.Error(err)
		}
	}()

	w := csv.NewWriter(csvFile)
	defer w.Flush()

	ppMtx := PartnerPaymentToMatrix(pp)

	if err = w.WriteAll(ppMtx); err != nil {
		return fmt.Errorf("writeAll: %v", err)
	}

	return nil
}

func PartnerPaymentToMatrix(payments []models.PartnerPayment) (mtx [][]string) {
	mtx = append(mtx, []string{"id", "amount", "currency", "account", "bank", "status", "created_at", "updated_at"})

	for _, el := range payments {
		mtx = append(mtx, []string{strconv.Itoa(el.Id), strconv.Itoa(el.Amount), el.Currency, el.Account,
			el.Bank, el.Status, fmt.Sprintf("%v", el.CreatedAt), fmt.Sprintf("%v", el.UpdatedAt)})
	}

	return mtx
}

const (
	DirAlifInput = "./files/alif/input"
	DirIBTInput  = "./files/ibt/input"
	DirDCInput   = "./files/dc/input"
)

func (s *Service) SaveToFile(pp []models.PartnerPayment) error {
	var bankAlif, bankDc, bankIBT []models.PartnerPayment

	for _, payment := range pp {
		switch strings.ToLower(payment.Bank) {
		case BankALIF:
			bankAlif = append(bankAlif, payment)
		case BankDC:
			bankDc = append(bankDc, payment)
		case BankIBT:
			bankIBT = append(bankIBT, payment)
		default:
			logger.Error("partner err: ", payment.Bank)
		}
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		if err := s.SaveInput(DirIBTInput, bankIBT); err != nil {
			logger.Error("Error in saving IBT payments: ", err)
			return
		}
		wg.Done()
	}()
	go func() {
		if err := s.SaveInput(DirAlifInput, bankAlif); err != nil {
			logger.Error("Error in saving IBT payments: ", err)
			return
		}
		wg.Done()
	}()
	go func() {
		if err := s.SaveInput(DirDCInput, bankDc); err != nil {
			logger.Error("Error in saving IBT payments: ", err)
			return
		}
		wg.Done()
	}()

	wg.Wait()
	return nil
}

func (s *Service) ReadCSVFile() ([]models.PartnerPayment, error) {
	filePath := "./files/alif/input/input_19.06.2023.csv"

	file, err := os.Open(filePath)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	logger.Infoln("Successfully opened the CSV file")

	defer func() {
		if err := file.Close(); err != nil {
			logger.Error(err)
		}
	}()

	fileReader := csv.NewReader(file)
	records, err := fileReader.ReadAll()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var payments []models.PartnerPayment

	for ind, record := range records {
		if ind == 0 {
			continue
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			log.Println()
			break
		}

		amount, err := strconv.Atoi(record[1])
		if err != nil {
			log.Println()
			break
		}

		startD, _ := time.Parse("2006-01-02", record[6])
		updateD, _ := time.Parse("2006-01-02", record[7])

		payments = append(payments, models.PartnerPayment{
			Id:        id,
			Amount:    amount,
			Currency:  record[2],
			Account:   record[3],
			Bank:      record[4],
			Status:    record[5],
			CreatedAt: &startD,
			UpdatedAt: &updateD,
		})
	}

	return payments, nil
}

func (s *Service) Compare() (output []models.Output, err error) {
	dataFromHumo, err := s.Repo.GetDataFromHumo()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	dataFromCSVFile, err := s.ReadCSVFile()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var exists bool

	for _, humoPayment := range dataFromHumo {
		for _, partnerPayment := range dataFromCSVFile {
			if humoPayment.AgentId == partnerPayment.Id && humoPayment.Amount == partnerPayment.Amount {
				exists = true
				continue
			}
		}
		if !exists {
			output = append(output, models.Output{
				Id:        humoPayment.Id,
				AgentId:   humoPayment.AgentId,
				Amount:    humoPayment.Amount,
				Currency:  humoPayment.Currency,
				Account:   humoPayment.Account,
				Status:    humoPayment.Status,
				CreatedAt: humoPayment.CreatedAt,
				UpdatedAt: humoPayment.UpdatedAt,
			})
		}
		exists = false
	}

	for _, partnerPayment := range dataFromCSVFile {
		for _, humoPayment := range dataFromHumo {
			if partnerPayment.Id == humoPayment.AgentId && partnerPayment.Amount == humoPayment.Amount {
				exists = true
				continue
			}
		}
		if !exists {
			output = append(output, models.Output{
				Id:        partnerPayment.Id,
				Amount:    partnerPayment.Amount,
				Currency:  partnerPayment.Currency,
				Account:   partnerPayment.Account,
				Bank:      partnerPayment.Bank,
				Status:    partnerPayment.Status,
				CreatedAt: partnerPayment.CreatedAt,
				UpdatedAt: partnerPayment.UpdatedAt,
			})
		}
		exists = false
	}

	err = s.SaveOutput("./files/output", output)
	if err != nil {
		logger.Error(err)
		return
	}

	return output, nil
}

func OutputToMtx(output []models.Output) (mtx [][]string) {
	mtx = append(mtx, []string{"id", "agent_id", "amount", "currency", "account", "bank", "status", "created_at", "updated_at"})

	for _, out := range output {
		mtx = append(mtx, []string{strconv.Itoa(out.Id), strconv.Itoa(out.AgentId), strconv.Itoa(out.Amount), out.Currency, out.Account,
			out.Bank, out.Status, fmt.Sprintf("%v", out.CreatedAt), fmt.Sprintf("%v", out.UpdatedAt)})
	}

	return mtx
}

func (s *Service) SaveOutput(pathDir string, output []models.Output) error {
	err := os.MkdirAll(pathDir, 0777)
	if err != nil {
		logger.Errorf("error while creating output dir: %v", err)
	}

	date := time.Now().Format("02.01.2006")
	pathFile := fmt.Sprintf("%s/output_%v.csv", pathDir, date)

	csvFile, err := os.Create(pathFile)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer func() {
		if err := csvFile.Close(); err != nil {
			logger.Error(err)
		}
	}()

	csvWriter := csv.NewWriter(csvFile)

	outputMtx := OutputToMtx(output)

	if err = csvWriter.WriteAll(outputMtx); err != nil {
		return err
	}

	return nil
}

func (s *Service) ActualizeReester() {
	//1. Get all partners

	//2.
}

func (s *Service) JobReestr() {
	var integrationDate models.Partner
	//var intervalType models.Partner
	_, err := s.CheckInterval(integrationDate.IntegrationDate)
	if err != nil {
		logger.Error("error in s.CheckInterval", err)
	}
}

func (s *Service) GetAllPartners() (partners []models.Partner, err error) {
	// 1. Get all partners
	partners, err = s.Repo.GetAllPartners()
	if err != nil {
		logger.Errorf("error while getting all partners from DB: %v", err)
		return
	}

	// 2. Iterate through each partner
	//for _, partner := range partners {
	//	//checkInterval(partner)
	//}

	for _, partner := range partners {
		switch partner.Name {
		case BankALIF:

		}
	}

	// 3. Check interval type

	// 4. Create registers from interval_date until curr_time depending from interval_type

	return
}

type TimeS struct {
	Days int
}

// intervalType models.Partner --> входные данные
func (s *Service) CheckInterval(integrationDate string) (*TimeS, error) {
	currentTime := time.Now()

	parseIntegrationDate, err := time.Parse("2006-01-02", integrationDate)
	if err != nil {
		logger.Error("error while parsing integration date", err)
		return nil, err
	}

	if parseIntegrationDate.After(currentTime) {
		logger.Error("integration date is greater that current date", err)
		return nil, err
	}

	difference := currentTime.Sub(parseIntegrationDate)

	days := int(difference.Hours() / 24)
	fmt.Printf("разница %0.d \n", days)

	var reestrs []models.Reestr

	var id models.Partner

	// must create 3 funcs for intervalTypes. it means if intervalType will be D, it will add one day etc.
	// must do switch cases for Partners. if partner will be alif, so alif has his own FOR LOOP to adding reestrs

	for i := 1; i <= days; i++ {
		parseIntegrationDate = parseIntegrationDate.AddDate(0, 0, 1)
		fileName := fmt.Sprintf("input_%v.csv", parseIntegrationDate.Format("2006-01-02"))

		reestr := models.Reestr{
			Id:        i,
			PartnerId: id.Id,
			FileName:  fileName,
		}
		reestrs = append(reestrs, reestr)
	}
	fmt.Printf("%+v\n", reestrs)

	err = s.Repo.NewReestrRecord(reestrs)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &TimeS{Days: days}, nil
}
