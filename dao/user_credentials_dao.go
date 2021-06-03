package dao

import (
	"apuscorp.com/core/auth-service/client"
	"apuscorp.com/core/auth-service/database"
	"apuscorp.com/core/auth-service/dto/apperror"
	"apuscorp.com/core/auth-service/dto/request"
	"apuscorp.com/core/auth-service/dto/response"
	"apuscorp.com/core/auth-service/i18n"
	"apuscorp.com/core/auth-service/model"
	"apuscorp.com/core/auth-service/util"
	"apuscorp.com/core/auth-service/util/constant"
	"apuscorp.com/core/auth-service/util/errorcode"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

type UserCredentialsDao struct {
	Db     *database.Database
	I18n   *i18n.I18n
	Client *client.Client
}

// Login Chức năng login
func (r *UserCredentialsDao) Login(c *gin.Context, loginData *request.UserLoginDto) (*response.UserLoginRes, error) {
	// validate
	if loginData.Mobile == nil && loginData.Email == nil && loginData.Username == nil {
		return nil, apperror.NewAppErrorWithStatus(http.StatusBadRequest, errorcode.AUT1008, "", c.GetHeader(constant.HeaderAcceptLanguage), nil)
	}

	var cred *string = nil
	credType := constant.Mobile
	if loginData.Mobile != nil {
		cred = loginData.Mobile
		credType = constant.Mobile
	} else if loginData.Email != nil {
		cred = loginData.Email
		credType = constant.Email
	} else if loginData.Username != nil {
		cred = loginData.Username
		credType = constant.Username
	}
	if cred == nil {
		return nil, apperror.NewAppErrorWithStatus(http.StatusBadRequest, errorcode.AUT1008, "", c.GetHeader(constant.HeaderAcceptLanguage), nil)
	}

	// Lấy danh sách user_credential có username
	creds := make([]model.UserCredential, 0)
	if err := r.Db.DB.Where("type = ? and credential = ?", credType, *cred).Find(&creds).Error; err != nil {
		return nil, apperror.AppError{
			Status:       http.StatusUnauthorized,
			ErrorCode:    errorcode.AUT1001,
			ErrorMessage: r.I18n.MustLocalize("vi", errorcode.AUT1001, nil),
		}
	}
	if len(creds) == 0 {
		return nil, apperror.AppError{
			Status:       http.StatusUnauthorized,
			ErrorCode:    errorcode.AUT1001,
			ErrorMessage: r.I18n.MustLocalize("vi", errorcode.AUT1001, nil),
		}
	}

	// iterate danh sách user_credentials và so sánh password hash
	for _, cred := range creds {
		err := bcrypt.CompareHashAndPassword([]byte(cred.Password), []byte(loginData.Password))
		if err == nil {
			var user model.User
			if err := r.Db.DB.Where("id = ?", cred.UserId).First(&user).Error; err != nil {
				return nil, apperror.AppError{
					Status:       http.StatusUnauthorized,
					ErrorCode:    errorcode.AUT1002,
					ErrorMessage: r.I18n.MustLocalize("vi", errorcode.AUT1002, nil),
				}
			}
			// get company ID
			res, err := r.Client.MembershipClient.R().
				SetResult(&response.Response{}).
				SetError(&response.Response{}).
				SetQueryParam("user-id", strconv.Itoa(int(cred.UserId))).
				Get("/api/v1/user-companies/my-company-id")
			util.Must(err)
			if res.IsError() {
				return nil, util.MustResponse(res.Error().(*response.Response))
			}

			return &response.UserLoginRes{
				UserId:     user.ID,
				CompanyId:  uint64(res.Result().(*response.Response).Data.(float64)),
				RememberMe: loginData.RememberMe,
			}, nil
		} else if !errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, apperror.AppError{
				Status:       http.StatusUnauthorized,
				ErrorCode:    errorcode.AUT1001,
				ErrorMessage: r.I18n.MustLocalize("vi", errorcode.AUT1001, nil),
			}
		}
	}

	return nil, apperror.AppError{
		Status:       http.StatusUnauthorized,
		ErrorCode:    errorcode.AUT1001,
		ErrorMessage: r.I18n.MustLocalize("vi", errorcode.AUT1001, nil),
	}
}
