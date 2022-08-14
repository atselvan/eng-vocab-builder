package main

import (
	"fmt"
	"net/http"

	"github.com/atselvan/ankiconnect"
	"github.com/atselvan/dictionaryapi"
	"github.com/gin-gonic/gin"
	"github.com/privatesquare/bkst-go-utils/utils/errors"
	"github.com/privatesquare/bkst-go-utils/utils/httputils"
	"github.com/privatesquare/bkst-go-utils/utils/slice"
)

type server struct {
	cnf *Config
}

func (s *server) GetWordsWord(c *gin.Context, word string) {
	meaning, restErr := getMeaningOfWord(word)
	if restErr != nil {
		c.JSON(restErr.StatusCode, restErr)
		return
	}
	c.JSON(http.StatusOK, meaning)
}

func (s *server) PostWordsWord(c *gin.Context, word string) {
	// check deck status
	restErr := checkDeck(s.cnf)
	if restErr != nil {
		c.JSON(restErr.StatusCode, restErr)
		return
	}

	// get meaning of the word
	meaning, restErr := getMeaningOfWord(word)
	if restErr != nil {
		c.JSON(restErr.StatusCode, restErr)
		return
	}

	// add word to anki deck
	restErr = addWordToAnki(s.cnf, meaning)
	if restErr != nil {
		c.JSON(restErr.StatusCode, restErr)
		return
	}

	c.JSON(http.StatusOK, httputils.RestMsg{Message: "Word added to anki successfully"})
}

// getMeaningOfWord gets the meaning of a word using dictionaryapi
func getMeaningOfWord(word string) (*dictionaryapi.Word, *errors.RestErr) {
	client := dictionaryapi.NewClient()
	return client.Word.Get(word)
}

// checkDeck checks if the deck exists, if not the functions creates the deck.
func checkDeck(cnf *Config,) *errors.RestErr {
	client := ankiconnect.NewClient().SetURL(cnf.AnkiConnectURL)
	decks, restErr := client.Decks.GetAll()
	if restErr != nil {
		return restErr
	}
	if !slice.EntryExists(*decks, cnf.AnkiDeckName) {
		restErr := client.Decks.Create(cnf.AnkiDeckName)
		if restErr != nil {
			return restErr
		}
	}
	return nil
}

// addWordToAnki adds the word to the anki deck
func addWordToAnki(cnf *Config, word *dictionaryapi.Word) *errors.RestErr {

	meanings := ""
	for _, meaning := range word.Meanings {

		for _, definition := range meaning.Definitions {
			meanings += fmt.Sprintf("<br><b>Part of Speech:</b> %s", meaning.PartOfSpeech)
			meanings += fmt.Sprintf("<br><b>Definition :</b> %s", definition.Definition)
			if definition.Example != "" {
				meanings += fmt.Sprintf("<br><b>Example :</b> %s", definition.Example)
			}
			if len(definition.Synonyms) > 0 {
				meanings += fmt.Sprintf("<br><b>Synonyms :</b> %v", definition.Synonyms)
			}
			meanings += "<br>"
		}
	}

	note := ankiconnect.Note{
		DeckName:  cnf.AnkiDeckName,
		ModelName: cnf.AnkiDeckModel,
		Fields: ankiconnect.Fields{
			Front: word.Word,
			Back:  meanings,
		},
		Audio: []ankiconnect.Audio{
			{
				URL:    word.Phonetics[0].Audio,
				Fields: []string{word.Phonetics[0].Text},
			},
		},
	}

	client := ankiconnect.NewClient().SetURL(cnf.AnkiConnectURL)
	restErr := client.Notes.Add(note)
	if restErr != nil {
		return restErr
	}
	return nil
}
