package db

import (
	"MdShorts/pkg/api"
	"MdShorts/pkg/entity"
	"MdShorts/pkg/repository"
	"encoding/base64"
	"math/rand"
	"strconv"
	"strings"

	"MdShorts/pkg/repository/user_login_history_repository"
	"errors"
	"time"

	"github.com/aekam27/trestCommon"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var (
	repo = repository.NewProfileRepository("users")
)

var (
	loginhistoryrepo = user_login_history_repository.NewUserLoginHistoryRepository("usersLoginHistory")
)

type accountService struct{}

func NewSignUpService(repository repository.ProfileRepository) AccountService {
	repo = repository
	return &accountService{}
}
func (*accountService) SignUp(cred Credentials) (string, entity.ProfileDB, error) {
	if cred.Password == "" {
		err := errors.New("password missing")
		trestCommon.ECLog2("sign up failed no password", err)
		return "", entity.ProfileDB{}, err
	}
	_, err := checkUser(cred.Email)

	if err != nil {
		trestCommon.ECLog2("sign up user not found", err)
		if err.Error() == "mongo: no documents in result" {
			var serv *accountService
			tokenString, err := serv.hashAndInsertData(cred)
			var profile entity.ProfileDB
			profile.Email = cred.Email
			profile.FirstName = cred.FirstName
			profile.LastName = cred.LastName
			profile.PhoneNo = cred.PhoneNo
			profile.CountryCode = cred.CountryCode
			profile.Speciality = cred.Speciality
			if err != nil {
				trestCommon.ECLog3("sign up not successful", err, logrus.Fields{"email": cred.Email})
			}
			return tokenString, profile, err
		} else {
			return "", entity.ProfileDB{}, err

		}
	}

	return "", entity.ProfileDB{}, errors.New("email already registered")
}

func (*accountService) SendVerificationEmail(email, pemail, uid string) (string, error) {
	emailSentTime := time.Now()
	if email == "" {
		return "email sent successfully", nil
	}
	verificationCode := trestCommon.GetRandomString(16)
	sendCode, err := trestCommon.Encrypt(email + ":" + verificationCode)
	if err != nil {
		trestCommon.ECLog2("send verification email encryption failed", err)
		return "", err
	}
	_, err = trestCommon.SendVerificationCode(email, sendCode)
	if err != nil {
		trestCommon.ECLog2("send verification email failed", err)
		return "", err
	}
	if pemail != "" {
		_, err = repo.UpdateOne(bson.M{"email": pemail}, bson.M{"$set": bson.M{"email": email, "email_sent_time": emailSentTime, "verification_code": verificationCode}})
		if err != nil {
			trestCommon.ECLog2("send verification email update failed", err)
			return "", err
		}
		return trestCommon.CreateToken(uid, email, "", "")
	} else {
		_, err = repo.UpdateOne(bson.M{"email": email}, bson.M{"$set": bson.M{"email_sent_time": emailSentTime, "verification_code": verificationCode}})
		if err != nil {
			trestCommon.ECLog2("send verification email update failed", err)
			return "", err
		}
		return "email sent successfully", nil
	}

}

func (*accountService) SendOTP(email, mnumber string) (string, error) {
	otpSentTime := time.Now()
	auth := viper.GetString("clicksend.username") + ":" + viper.GetString("clicksend.password")
	bauth := base64.StdEncoding.EncodeToString([]byte(auth))
	randomOTP := 1000 + rand.Intn(9999-1000)
	_, err := api.ClickSend(bauth, mnumber, randomOTP)
	if err != nil {
		trestCommon.ECLog2("send verification sms failed", err)
		return "", err
	}
	_, err = repo.UpdateOne(bson.M{"email": email}, bson.M{"$set": bson.M{"otp_sent_time": otpSentTime, "otp_code": strconv.Itoa(randomOTP)}})
	if err != nil {
		trestCommon.ECLog2("send verification OTP update failed", err)
		return "", err
	}
	return "OTP sent successfully", nil
}

func (*accountService) ResendOTP(cred OTP) (string, error) {
	otpSentTime := time.Now()
	auth := viper.GetString("clicksend.username") + ":" + viper.GetString("clicksend.password")
	bauth := base64.StdEncoding.EncodeToString([]byte(auth))
	randomOTP := 1000 + rand.Intn(9999-1000)
	userData, _ := checkUser(cred.Email)
	_, err := api.ClickSend(bauth, userData.PhoneNo, randomOTP)
	if err != nil {
		trestCommon.ECLog2("send verification sms failed", err)
		return "", err
	}
	_, err = repo.UpdateOne(bson.M{"email": cred.Email}, bson.M{"$set": bson.M{"otp_sent_time": otpSentTime, "otp_code": strconv.Itoa(randomOTP)}})
	if err != nil {
		trestCommon.ECLog2("send verification OTP update failed", err)
		return "", err
	}
	return "OTP sent successfully", nil
}

// https://firebase.flutter.dev/docs/analytics/usage
// https://medium.flutterdevs.com/firebase-analytics-2044e865efc4

func (*accountService) Login(cred Credentials) (string, entity.ProfileDB, error) {
	if cred.Password == "" {
		err := errors.New("password missing")
		trestCommon.ECLog2("login failed no password", err)
		return "", entity.ProfileDB{}, err
	}
	salt := viper.GetString("salt")
	userData, err := checkUser(cred.Email)
	if err != nil {
		trestCommon.ECLog2("login failed user not found", err)
		return "", entity.ProfileDB{}, err
	}
	if userData.Status == "created" {
		err = errors.New("user not verified")
		trestCommon.ECLog2("login failed user has not verified his/her email address", err)
		return "", entity.ProfileDB{}, err
	}
	if userData.Status == "deleted" || userData.Status == "archived" {
		err = errors.New("unauthorized user")
		trestCommon.ECLog2("login failed user has deleted or archived his profile", err)
		return "", entity.ProfileDB{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(cred.Password+salt))
	if err != nil {
		trestCommon.ECLog2("login failed password hash doesn't match", err)
		return "", entity.ProfileDB{}, err
	}
	tokenString, err := trestCommon.CreateToken(userData.ID.Hex(), cred.Email, "", userData.Status)
	if err != nil {
		trestCommon.ECLog3("login failed unable to create token", err, logrus.Fields{"email": cred.Email, "name": userData.FirstName, "status": userData.Status})
		return "", entity.ProfileDB{}, err
	}
	repo.UpdateOne(bson.M{"_id": userData.ID}, bson.M{"$set": bson.M{"login_time": time.Now(), "last_login_device_info": cred.LastLoginDeviceInfo, "last_login_location": cred.LastLoginLocation}})
	loginhistoryrepo.InsertOne(bson.M{"user_id": userData.ID.Hex(), "login_time": time.Now(), "last_login_device_info": cred.LastLoginDeviceInfo, "last_login_location": cred.LastLoginLocation})
	userData.Password = ""
	userData.PasswordResetCode = ""
	userData.VerificationCode = ""
	newUrl := createPreSignedDownloadUrl(userData.UrlToProfileImage)
	userData.UrlToProfileImage = newUrl
	return tokenString, userData, nil
}
func createPreSignedDownloadUrl(url string) string {
	s := strings.Split(url, "?")
	if len(s) > 0 {
		o := strings.Split(s[0], "/")
		if len(o) > 3 {
			fileName := o[4]
			path := o[3]
			downUrl, _ := trestCommon.PreSignedDownloadUrlAWS(fileName, path)
			return downUrl
		}
	}
	return ""
}
func (*accountService) VerifyEmail(cred Credentials) (string, error) {

	userData, err := checkUser(cred.Email)
	if err != nil {
		trestCommon.ECLog3("verify user not found", err, logrus.Fields{"email": cred.Email})
		return "", err
	}

	if cred.VerificationCode != userData.VerificationCode {
		err = errors.New("unauthorized user")
		trestCommon.ECLog3("verify user verification code didn't match", err, logrus.Fields{"email": cred.Email, "db verify code": userData.VerificationCode, "code provided by user": cred.VerificationCode})
		return "", err
	}
	if userData.Status == "verified" {
		err = errors.New("user already verified")
		trestCommon.ECLog3("verify user verification user already verified", err, logrus.Fields{"email": cred.Email})
		return "", err
	}
	_, err = repo.UpdateOne(bson.M{"_id": userData.ID}, bson.M{"$set": bson.M{"verified_time": time.Now(), "status": "verified"}})
	if err != nil {
		trestCommon.ECLog3("verify user unable to update status", err, logrus.Fields{"email": cred.Email})
		return "", err
	}
	return "verified", nil
}

func (*accountService) VerifyOTP(cred OTP) (string, error) {

	userData, err := checkUser(cred.Email)
	if err != nil {
		trestCommon.ECLog3("verify user not found", err, logrus.Fields{"email": cred.Email})
		return "", err
	}

	if cred.OTP != userData.OTP {
		err = errors.New("unauthorized user")
		trestCommon.ECLog3("verify user verification code didn't match", err, logrus.Fields{"email": cred.Email, "db verify code": userData.OTP, "code provided by user": cred.OTP})
		return "", err
	}
	_, err = repo.UpdateOne(bson.M{"_id": userData.ID}, bson.M{"$set": bson.M{"verified_time": time.Now(), "phone_status": "verified"}})
	if err != nil {
		trestCommon.ECLog3("verify user unable to update status", err, logrus.Fields{"email": cred.Email})
		return "", err
	}
	return "verified", nil
}

func checkUser(email string) (entity.ProfileDB, error) {
	var userData entity.ProfileDB
	if email == "" {
		err := errors.New("email missing")
		trestCommon.ECLog2("check user failed no email", err)
		return userData, err
	}
	if !trestCommon.ValidateEmail(email) {
		err := errors.New("invalid email")
		trestCommon.ECLog2("check user failed invalid email", err)
		return userData, err
	}

	return repo.FindOne(bson.M{"email": email}, bson.M{})

}

func (*accountService) hashAndInsertData(cred Credentials) (string, error) {
	salt := viper.GetString("salt")

	hash, err := bcrypt.GenerateFromPassword([]byte(cred.Password+salt), 5)
	if err != nil {
		trestCommon.ECLog3("hash password", err, logrus.Fields{"email": cred.Email})
		return "", err
	}
	cred.Password = string(hash)
	cred.CreatedDate = time.Now()
	cred.Status = "created"
	userid, err := repo.InsertOne(cred)
	if err != nil {
		trestCommon.ECLog3("hashAndInsertData Insert failed", err, logrus.Fields{"email": cred.Email})
		return "", nil
	}
	var serv accountService
	_, err = serv.SendVerificationEmail(cred.Email, "", "")
	if err != nil {
		trestCommon.ECLog3("hashAndInsertData Insert failed", err, logrus.Fields{"email": cred.Email})
	}
	_, err = serv.SendOTP(cred.Email, cred.PhoneNo)
	if err != nil {
		trestCommon.ECLog3("hashAndInsertData Insert failed", err, logrus.Fields{"email": cred.Email})
	}
	return trestCommon.CreateToken(userid, cred.Email, "", cred.Status)
}
func (*accountService) SendResetLink(email string) (string, error) {
	var cred Credentials
	cred.Email = email
	_, err := checkUser(cred.Email)
	if err != nil {
		trestCommon.ECLog2("user not found", err)
		return "", err

	}
	emailSentTime := time.Now()
	verificationCode := trestCommon.GetRandomString(16)
	resetCode, err := trestCommon.Encrypt(email + ":" + verificationCode)
	if err != nil {
		trestCommon.ECLog2("send reset link encryption failed", err)
		return "", err
	}
	_, err = trestCommon.SendResetPasswordLink(email, resetCode)
	if err != nil {
		trestCommon.ECLog2("send reset password link failed", err)
		return "", err
	}
	_, err = repo.UpdateOne(bson.M{"email": email}, bson.M{"$set": bson.M{"email_sent_time": emailSentTime, "password_reset_code": verificationCode}})
	if err != nil {
		trestCommon.ECLog2("send reset link update failed", err)
		return "", err
	}
	return "Reset link sent successfully", nil
}

func (*accountService) VerifyResetLink(cred Credentials) (string, string, error) {

	userData, err := checkUser(cred.Email)
	if err != nil {
		trestCommon.ECLog3("verify user not found", err, logrus.Fields{"email": cred.Email})
		return "", "", err
	}

	if cred.PasswordResetCode != userData.PasswordResetCode {
		err = errors.New("unauthorized user")
		trestCommon.ECLog3("verify user password reset code didn't match", err, logrus.Fields{"email": cred.Email, "db verify code": userData.PasswordResetCode, "code provided by user": cred.PasswordResetCode})
		return "", "", err
	}
	_, err = repo.UpdateOne(bson.M{"_id": userData.ID}, bson.M{"$set": bson.M{"password_reset_time": time.Now()}})
	if err != nil {
		trestCommon.ECLog3("verify user unable to update status", err, logrus.Fields{"email": cred.Email})
		return "", "", err
	}

	return "verified", userData.Email, nil
}

func (*accountService) UpdatePassword(cred Credentials) (string, error) {
	if cred.Password == "" {
		err := errors.New("password missing")
		trestCommon.ECLog2("update password failed no password", err)
		return "", err
	}

	_, err := checkUser(cred.Email)

	if err != nil {
		trestCommon.ECLog2("update password failed user not found", err)
		return "", err
	}

	salt := viper.GetString("salt")

	hash, err := bcrypt.GenerateFromPassword([]byte(cred.Password+salt), 5)
	if err != nil {
		trestCommon.ECLog3("hash password", err, logrus.Fields{"email": cred.Email})
		return "", err
	}
	cred.Password = string(hash)
	cred.CreatedDate = time.Now()
	_, err = repo.UpdateOne(bson.M{"email": cred.Email}, bson.M{"$set": bson.M{"password": cred.Password, "update_time": time.Now()}})
	if err != nil {
		trestCommon.ECLog2("password update failed", err)
		return "", err
	}

	return "password updated successfully", nil
}
