/**
  @Author : hanxiaodong
*/

package web

import (
	"net/http"
	"fmt"
	"wu/web/controller"
)


// Start the web service and specify routing information
func WebStart(app controller.Application) {

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Specify routing information (matching requests)
	http.HandleFunc("/", app.LoginView)
	http.HandleFunc("/login", app.Login)
	http.HandleFunc("/loginout", app.LoginOut)

	http.HandleFunc("/index", app.Index)
	http.HandleFunc("/help", app.Help)

	http.HandleFunc("/addCertInfo", app.AddCertShow) // show AddInfo page
	http.HandleFunc("/addCert", app.AddCert)         // Submit an information request

	http.HandleFunc("/queryPage", app.QueryPage)       // Go to the search information page by certificate number and name
	http.HandleFunc("/query", app.FindCertByNoAndName) // Search information by certificate number and name

	http.HandleFunc("/queryPage2", app.QueryPage2) // Go to Search Information by ID Number page
	http.HandleFunc("/query2", app.FindByID)       // Search information by ID number

	http.HandleFunc("/modifyPage", app.ModifyShow) // Modify information page
	http.HandleFunc("/modify", app.Modify)         //  Modify information

	http.HandleFunc("/upload", app.UploadFile)

	fmt.Println("Start the web service, listening on port number: 9000")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Printf("Web service start failure: %v", err)
	}

}
