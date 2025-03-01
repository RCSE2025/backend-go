//package service
//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"github.com/RCSE2025/backend-go/internal/repo"
//	"io/ioutil"
//	"net/http"
//)
//
//type BusinessService struct {
//	repo *repo.BusinessRepo
//}
//
//func NewBusinessService(repo *repo.BusinessRepo) *BusinessService {
//	return &BusinessService{
//		repo: repo,
//	}
//}
//

package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/RCSE2025/backend-go/internal/model"
	"github.com/RCSE2025/backend-go/internal/repo"
	"io/ioutil"
	"net/http"
)

var (
	ErrBusinessNotFound = errors.New("business not found")
	ErrBusinessExists   = errors.New("business already exists")
)

type BusinessService struct {
	repo     *repo.BusinessRepo
	userRepo *repo.UserRepo
}

func NewBusinessService(repo *repo.BusinessRepo, userRepo *repo.UserRepo) *BusinessService {
	return &BusinessService{
		repo:     repo,
		userRepo: userRepo,
	}
}

const (
	url   = "http://suggestions.dadata.ru/suggestions/api/4_1/rs/findById/party"
	token = "a88bd145b3272ed1e6876a27d2d592d8bb473865"
)

type Request struct {
	Query      string `json:"query"`
	BranchType string `json:"branch_type"`
}

func (s *BusinessService) GetBusinessInfoByINN(inn string) (any, error) {
	reqData := Request{
		Query:      inn,
		BranchType: "MAIN",
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(body)

	var response any
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *BusinessService) CreateBusiness(userID int64, business model.Business) error {
	return s.repo.CreateBusiness(userID, business)
}

func (s *BusinessService) GetAllBusinesses() ([]model.Business, error) {
	return s.repo.GetAllBusinesses()
}

func (s *BusinessService) GetBusinessByID(id int64) (model.Business, error) {
	exists, err := s.repo.BusinessExists(id)
	if err != nil {
		return model.Business{}, err
	}

	if !exists {
		return model.Business{}, ErrBusinessNotFound
	}

	return s.repo.GetBusinessByID(id)
}

func (s *BusinessService) UpdateBusiness(id int64, business model.Business) error {
	exists, err := s.repo.BusinessExists(id)
	if err != nil {
		return err
	}

	if !exists {
		return ErrBusinessNotFound
	}

	return s.repo.UpdateBusiness(id, business)
}

func (s *BusinessService) DeleteBusiness(id int64) error {
	exists, err := s.repo.BusinessExists(id)
	if err != nil {
		return err
	}

	if !exists {
		return ErrBusinessNotFound
	}

	return s.repo.DeleteBusiness(id)
}

func (s *BusinessService) GetBusinessByINN(inn int64) (model.Business, error) {
	business, err := s.repo.GetBusinessByINN(inn)
	if err != nil {
		return model.Business{}, ErrBusinessNotFound
	}

	return business, nil
}

func (s *BusinessService) GetBusinessByOGRN(ogrn int64) (model.Business, error) {
	business, err := s.repo.GetBusinessByOGRN(ogrn)
	if err != nil {
		return model.Business{}, ErrBusinessNotFound
	}

	return business, nil
}

func (s *BusinessService) GetUserBusinesses(userID int64) ([]model.Business, error) {
	exists, err := s.userRepo.UserExists(userID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrUserNotFound
	}

	return s.repo.GetUserBusinesses(userID)
}

func (s *BusinessService) GetBusinessUsers(businessID int64) ([]model.User, error) {
	exists, err := s.repo.BusinessExists(businessID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrBusinessNotFound
	}

	return s.repo.GetBusinessesUsers(businessID)
}

func (s *BusinessService) AddUserToBusiness(userID int64, businessID int64) error {
	userExists, err := s.userRepo.UserExists(userID)
	if err != nil {
		return err
	}

	if !userExists {
		return ErrUserNotFound
	}

	businessExists, err := s.repo.BusinessExists(businessID)
	if err != nil {
		return err
	}

	if !businessExists {
		return ErrBusinessNotFound
	}

	return s.repo.AddUserToBusiness(userID, businessID)
}

func (s *BusinessService) RemoveUserFromBusiness(userID int64, businessID int64) error {
	userExists, err := s.userRepo.UserExists(userID)
	if err != nil {
		return err
	}

	if !userExists {
		return ErrUserNotFound
	}

	businessExists, err := s.repo.BusinessExists(businessID)
	if err != nil {
		return err
	}

	if !businessExists {
		return ErrBusinessNotFound
	}

	return s.repo.RemoveUserFromBusiness(userID, businessID)
}
