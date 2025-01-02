package main

import (
	"bufio"
	"fmt"
	"gopkg.in/headzoo/surf.v1"
	"os"
	"strings"
)

var (
	homeDir, _ = os.UserHomeDir()
	credPath = homeDir + "/credentials.bak"
)

func gimmeThoseLines() []string {
	file, err := os.Open(credPath)
	var logins []string
	if err != nil {
		return logins
	}
	defer file.Close()
	gin := bufio.NewScanner(file)
	for gin.Scan() {
		logins = append(logins, gin.Text())
	}
	if err := gin.Err(); err != nil {
		fmt.Println("[ERROR] Error reading credentials >:")
	}
	return logins
}

func writeThoseLines(logins []string) error {
	file, err := os.Create(credPath)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, str := range logins {
		_, err = fmt.Fprintf(file, "%s\n", str)
		if err != nil {
			return err
		}
	}
	return nil
}

func saveCredentials(usr string, pas string) error {
	logins := gimmeThoseLines()
	lmp := make(map[string]string)
	for _, str := range logins {
		cur := strings.Split(str, ",")
		lmp[cur[0]] = cur[1]
	}
	lmp[usr] = pas
	logins = nil
	for name, pass := range lmp {
		logins = append(logins, fmt.Sprintf("%s,%s", name, pass))
	}
	return writeThoseLines(logins)
}

func getCredentials() (string, string) {
	logins := gimmeThoseLines()
	fmt.Println("Here's a list of saved logins . . .")
	for i, str := range logins {
		fmt.Printf("   %d. %s\n", i+1, strings.Split(str, ",")[0])
	}
	fmt.Printf("   %d. New Account!\n", len(logins)+1)
	fmt.Print("Select your login credentials: ")
	var id int
	fmt.Scanf("%d\n", &id)
	if id-1 >= 0 && id-1 < len(logins) {
		data := strings.Split(logins[id-1], ",")
		return data[0], data[1]
	}
	var usr, pas string
	fmt.Print("Enter your username: ")
	fmt.Scanf("%s\n", &usr)
	fmt.Print("Enter your password: ")
	fmt.Scanf("%s\n", &pas)
	return usr, pas
}

func login(usr string, pas string) error {
	netoo := surf.NewBrowser()
	err := netoo.Open("https://net2.sharif.edu/login")
	if err != nil {
		return err
	}
	form, err := netoo.Form(".wrap-login100 > form:nth-child(1)")
	if err != nil {
		return fmt.Errorf("[ERROR] Could not find the login form, you probably are already logged in <:")
	}
	form.Input("username", usr)
	form.Input("password", pas)
	if form.Submit() != nil {
		return fmt.Errorf("[ERROR] Login failed, possibly a wrong password!? >:")
	}
	return nil
}

func main() {
	usr, pas := getCredentials()
	err := login(usr, pas)
	if err != nil {
		panic(err)
	}
	fmt.Println("Login was successful! <:")
	saveCredentials(usr, pas)
}
