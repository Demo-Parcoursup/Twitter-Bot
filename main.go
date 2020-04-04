package main

// Import des packets utilisés
import (
	"encoding/json" // Utilisé pour le traitement des réponses en JSON des APIS de citations
	"fmt"
	"github.com/dghubble/go-twitter/twitter" // Pakcet pour l'API Twitter
	"github.com/dghubble/oauth1" // Packet pour une implémentation Golang de OAuth1 (nécessaire pour se connecter à l'API Twitter)
	"log"
	"math/rand"
	"net/http"
	"time"
)

/* Définition du client HTTP pour envoyer des requêtes
Avec un timeout déclaré à 10 secondes
 */
var HttpClient = &http.Client{Timeout: 10 * time.Second}

// Définition de la structure (dict) qui contient les clés d'API pour Twitter
type Parametres struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// Structure qui définie le JSON à récupérer de la première API de citations
type QuoteOne struct {
	Author	string	`json:"author"`
	Quote	string	`json:"quote"`
}

// Structure qui définie le JSON à récupérer de la seconde API de définition
type QuoteTwo struct {
	Author	string	`json:"author"`
	Quote	string	`json:"en"`
}

// Fonction qui envoie la requête à l'API et récupère le code JSON
func getJson(url string, target interface{}) error {
	r, err := HttpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

// Fonction qui met en place le Client Twitter afin d'envoyer des tweets
func getTwitterClient(creds *Parametres) (*twitter.Client, error) {
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	Client := config.Client(oauth1.NoContext, token)
	TwitterClient := twitter.NewClient(Client)

	// Verification des clés d'API
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	user, _, err := TwitterClient.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	log.Printf("Info utilisateur:\n%+v\n", user)
	return TwitterClient, nil
}

// Fonction qui envoie le tweet
func SendTweet(Quote string, Author string) (*twitter.Tweet, *http.Response) {
	// Normalement il est fortement recommendé de ne jamais inclure les clés API dans le code source principal
	// Mais j'y fais exception pour la simplicité de l'exemple
	creds := Parametres{
		AccessToken:       "XXXXXXXXXXXXXXXXXXXXXX", //A remplacer par les vrais clés API
		AccessTokenSecret: "XXXXXXXXXXXXXXXXXXXXXX",
		ConsumerKey:       "XXXXXXXXXXXXXXXXXXXXXX",
		ConsumerSecret:    "XXXXXXXXXXXXXXXXXXXXXX",
	}

	// Récupère le Client avec la fonction précédente
	client, err := getTwitterClient(&creds)
	if err != nil {
		log.Println("Erreur lors de l'intitation du Client Twitter")
		log.Println(err)
	}

	// Envoie le tweet
	tweet, response, err := client.Statuses.Update(fmt.Sprintf("“%s„\n\n- %s", Quote, Author), nil)
	if err != nil {
		log.Println(err)
	}

	return tweet, response
}

// Fonction principal
func main() {
	fmt.Println("Inititation du bot")

	// Mise en place d'un cycle (loop) infinie
	for {
		// Liste d'URL d'API de citations à utiliser
		URLs := []string{
			"https://programming-quotes-api.herokuapp.com/quotes/random",
			"http://quotes.stormconsultancy.co.uk/random.json",
		}
		// Mise en place d'un seed aléatoire afin de choisir au hasard une URL
		rand.Seed(time.Now().UnixNano())
		choix := URLs[rand.Intn(len(URLs))]
		fmt.Println(choix)
		// Si le choix tombe sur la première URL, utiliser la structure JSON adaptée
		if choix == "https://programming-quotes-api.herokuapp.com/quotes/random" {
			resp := QuoteTwo{}
			// Récupère le code JSON
			getJson(choix, &resp)
			fmt.Println(resp.Author)
			fmt.Println(resp.Quote)
			// Envoie le tweet à l'aide de la fonction précédente
			tweet, response := SendTweet(resp.Quote, resp.Author)
			log.Printf("%+v\n", response)
			log.Printf("%+v\n", tweet)
		} else {
			resp := QuoteOne{}
			// Récupère le code JSON
			getJson(choix, &resp)
			fmt.Println(resp.Author)
			fmt.Println(resp.Quote)
			// Envoie le tweet à l'aide de la fonction précédente
			tweet, response := SendTweet(resp.Quote, resp.Author)
			log.Printf("%+v\n", response)
			log.Printf("%+v\n", tweet)
		}
		// Met le script en pause le cycle pendant 24 heures
		time.Sleep(24 * time.Hour)

	}

}
