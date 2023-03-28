

package controller

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"education/service"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var cuser User

func (app *Application) LoginView(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "login.html", nil)
}

func (app *Application) Index(w http.ResponseWriter, r *http.Request) {
	ShowView(w, r, "index.html", nil)
}

func (app *Application) Help(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		CurrentUser User
	}{
		CurrentUser: cuser,
	}
	ShowView(w, r, "help.html", data)
}

// 用户登录
func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	loginName := r.FormValue("loginName")
	password := r.FormValue("password")

	var flag bool
	for _, user := range users {
		if user.LoginName == loginName && user.Password == password {
			cuser = user
			flag = true
			break
		}
	}

	data := &struct {
		CurrentUser User
		Flag        bool
	}{
		CurrentUser: cuser,
		Flag:        false,
	}

	if flag {
		// 登录成功
		ShowView(w, r, "index.html", data)
	} else {
		// 登录失败
		data.Flag = true
		data.CurrentUser.LoginName = loginName
		ShowView(w, r, "login.html", data)
	}
}

// 用户登出
func (app *Application) LoginOut(w http.ResponseWriter, r *http.Request) {
	cuser = User{}
	ShowView(w, r, "login.html", nil)
}

// 显示添加信息页面
func (app *Application) AddEduShow(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		CurrentUser User
		Msg         string
		Flag        bool
	}{
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
	}
	ShowView(w, r, "addEdu.html", data)
}

// 添加信息
func (app *Application) AddEdu(w http.ResponseWriter, r *http.Request) {
	const MySecret string = "abc&1*~#^2^#s0^=)^^7%b34"
	hashValue := GetSha256(r.FormValue("ciphertext"))
	note, err := Encrypt(hashValue, MySecret)
	if err != nil {
		fmt.Println("error encrypting your classified text: ", err)
	}
	encText, err := Encrypt(r.FormValue("ciphertext"), MySecret)
	if err != nil {
		fmt.Println("error encrypting your classified text: ", err)
	}
	edu := service.Education{
		AssetName:  r.FormValue("assetName"),
		OwnerID:    r.FormValue("ownerID"),
		State:      r.FormValue("state"),
		Version:    r.FormValue("version"),
		CertNo:     r.FormValue("certNo"),
		Ciphertext: encText,
		Note:       note,
	}

	app.Setup.SaveEdu(edu)
	/*transactionID, err := app.Setup.SaveEdu(edu)

	data := &struct {
		CurrentUser User
		Msg string
		Flag bool
	}{
		CurrentUser:cuser,
		Flag:true,
		Msg:"",
	}

	if err != nil {
		data.Msg = err.Error()
	}else{
		data.Msg = "信息添加成功:" + transactionID
	}*/

	//ShowView(w, r, "addEdu.html", data)
	r.Form.Set("certNo", edu.CertNo)
	r.Form.Set("assetName", edu.AssetName)
	app.FindCertByNoAndName(w, r)
}

func (app *Application) QueryPage(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		CurrentUser User
		Msg         string
		Flag        bool
	}{
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
	}
	ShowView(w, r, "query.html", data)
}

// 根据证书编号与姓名查询信息
func (app *Application) FindCertByNoAndName(w http.ResponseWriter, r *http.Request) {
	const MySecret string = "abc&1*~#^2^#s0^=)^^7%b34"
	certNo := r.FormValue("certNo")
	name := r.FormValue("assetName")
	result, err := app.Setup.FindEduByCertNoAndName(certNo, name)
	var edu = service.Education{}
	hashValue2 := GetSha256(r.FormValue("ciphertext"))
	hashValue1, err := Decrypt(edu.Note, MySecret)
	if err != nil {
		fmt.Println("error decrypting your encrypted text: ", err)
	}
	if strings.Compare(hashValue2, hashValue1) == 0 {
		fmt.Println("Content validation succeeded!")
		json.Unmarshal(result, &edu)

		fmt.Println("根据证书编号与姓名查询信息成功：")
		fmt.Println(edu)

		data := &struct {
			Edu         service.Education
			CurrentUser User
			Msg         string
			Flag        bool
			History     bool
		}{
			Edu:         edu,
			CurrentUser: cuser,
			Msg:         "",
			Flag:        false,
			History:     false,
		}

		if err != nil {
			data.Msg = err.Error()
			data.Flag = true
		}

		ShowView(w, r, "queryResult.html", data)
	} else {
		fmt.Println("Content validation failed!")
		data := &struct {
			Edu         service.Education
			CurrentUser User
			Msg         string
			Flag        bool
			History     bool
		}{
			Edu:         edu,
			CurrentUser: cuser,
			Msg:         "",
			Flag:        false,
			History:     false,
		}

		if err != nil {
			data.Msg = err.Error()
			data.Flag = true
		}
		ShowView(w, r, "help.html", data)
	}

}

func (app *Application) QueryPage2(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		CurrentUser User
		Msg         string
		Flag        bool
	}{
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
	}
	ShowView(w, r, "query2.html", data)
}

// 根据身份证号码查询信息
func (app *Application) FindByID(w http.ResponseWriter, r *http.Request) {
	ownerID := r.FormValue("ownerID")
	result, err := app.Setup.FindEduInfoByEntityID(ownerID)
	var edu = service.Education{}
	json.Unmarshal(result, &edu)

	data := &struct {
		Edu         service.Education
		CurrentUser User
		Msg         string
		Flag        bool
		History     bool
	}{
		Edu:         edu,
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
		History:     true,
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}

	ShowView(w, r, "queryResult.html", data)
}

// 修改/添加新信息
func (app *Application) ModifyShow(w http.ResponseWriter, r *http.Request) {
	// 根据证书编号与姓名查询信息
	certNo := r.FormValue("certNo")
	name := r.FormValue("assetName")
	result, err := app.Setup.FindEduByCertNoAndName(certNo, name)

	var edu = service.Education{}
	json.Unmarshal(result, &edu)

	data := &struct {
		Edu         service.Education
		CurrentUser User
		Msg         string
		Flag        bool
	}{
		Edu:         edu,
		CurrentUser: cuser,
		Flag:        true,
		Msg:         "",
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}

	ShowView(w, r, "modify.html", data)
}

// 修改/添加新信息
func (app *Application) Modify(w http.ResponseWriter, r *http.Request) {
	edu := service.Education{
		AssetName:  r.FormValue("assetName"),
		OwnerID:    r.FormValue("ownerID"),
		State:      r.FormValue("state"),
		Version:    r.FormValue("version"),
		CertNo:     r.FormValue("certNo"),
		Ciphertext: r.FormValue("ciphertext"),
		Note:       r.FormValue("note"),
	}

	//transactionID, err := app.Setup.ModifyEdu(edu)
	app.Setup.ModifyEdu(edu)

	/*data := &struct {
		Edu service.Education
		CurrentUser User
		Msg string
		Flag bool
	}{
		CurrentUser:cuser,
		Flag:true,
		Msg:"",
	}

	if err != nil {
		data.Msg = err.Error()
	}else{
		data.Msg = "新信息添加成功:" + transactionID
	}

	ShowView(w, r, "modify.html", data)
	*/

	r.Form.Set("certNo", edu.CertNo)
	r.Form.Set("name", edu.AssetName)
	app.FindCertByNoAndName(w, r)
}

var bytess = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func Encrypt(text, MySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}

	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytess)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)

	return Encode(cipherText), nil
}
func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func GetSha256(str string) string {
	m := sha256.New()
	m.Write([]byte(str))
	sha256String := hex.EncodeToString(m.Sum(nil))
	return sha256String
}
func Decrypt(text, MySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}

	cipherText := Decode(text)
	cfb := cipher.NewCFBDecrypter(block, bytess)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)

	return string(plainText), nil
}
func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}
