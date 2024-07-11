package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	desti := []string{
	
	}

	reader := bufio.NewReader(os.Stdin)

	// Demander l'adresse email de l'utilisateur
	fmt.Print("Entrez votre adresse mail : ")
	userMail, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	userMail = strings.TrimSpace(userMail)

	// Demander le mot de passe d'application à l'utilisateur
	fmt.Print("Entrez le mot de passe d'application : ")
	appPassword, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	appPassword = strings.TrimSpace(appPassword)

	// Demander le sujet de l'email
	fmt.Print("Entrez le sujet de l'email : ")
	subject, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	subject = strings.TrimSpace(subject)

	// Demander le body de l'email (lecture multi-lignes)
	fmt.Println("Entrez le body de l'email (appuyez sur Entrée deux fois pour terminer) :")
	var destiBody strings.Builder
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		destiBody.WriteString(line + "\n")
	}

	// Demander si l'utilisateur veut envoyer un fichier (cv ou autre)
	fmt.Print("Voulez-vous envoyer un fichier ? (Oui/Non) : ")
	sendFile, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	sendFile = strings.TrimSpace(sendFile)

	var fileName string
	if strings.ToLower(sendFile) == "oui" {
		fmt.Print("Entrez le nom du fichier à envoyer (avec l'extension, par exemple, 'CV_Luca.pdf') : ")
		fileName, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		fileName = strings.TrimSpace(fileName)
	}

	// Envoyer les emails
	d := gomail.NewDialer("smtp.gmail.com", 587, userMail, appPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // Activer la connexion TLS sans vérification du certificat

	for _, recipient := range desti {
		m := gomail.NewMessage()
		m.SetHeader("From", userMail)
		m.SetHeader("To", recipient)
		m.SetHeader("Subject", subject)
		m.SetBody("text/plain", destiBody.String())

		if strings.ToLower(sendFile) == "oui" {
			if _, err := os.Stat(fileName); err == nil {
				m.Attach(fileName)
			} else {
				log.Printf("Le fichier %s n'existe pas. Email non envoyé à %s", fileName, recipient)
				continue
			}
		}

		log.Printf("Envoi de l'email à %s...", recipient)
		if err := d.DialAndSend(m); err != nil {
			log.Println("Erreur lors de l'envoi de l'email:", err)
		} else {
			log.Println("Email envoyé avec succès à", recipient)
		}

		log.Println("Attente d'une minute avant d'envoyer le prochain email...")
		time.Sleep(1 * time.Minute)
	}

	log.Println("Tous les emails ont été envoyés.")
}
