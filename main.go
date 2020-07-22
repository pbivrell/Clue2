package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pbivrell/clue/clue"
	"github.com/pbivrell/clue/decks"
	"github.com/pbivrell/clue/sessions"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	ip := "192.168.86.248"
	//ip := "192.168.86.239"
	//ip := "localhost"
	//ip := "192.168.86.195"
	//ip := "71.237.93.61"
	//ip := "10.234.26.69"
	// Read deck
	deck, err := decks.FromFile("./html/cludo.json")
	if err != nil {
		panic(err)
	}

	// Track Existing Decks
	dm := decks.NewManager()

	// Add deck to manager
	dm.Insert(deck)

	// manage sessions
	sm := sessions.NewManager()

	r := mux.NewRouter()

	// API ep's
	r.HandleFunc("/deck", clue.Deck(sm, dm))
	r.HandleFunc("/users", clue.Users(sm))
	r.HandleFunc("/join/{sid}/", clue.Join(sm))
	r.HandleFunc("/create", clue.Create(sm, dm))
	r.HandleFunc("/input", clue.Input(sm))
	r.HandleFunc("/release", clue.Release(sm))
	r.HandleFunc("/get", clue.Retrieve(sm))
	r.HandleFunc("/setup", clue.Setup(sm))
	r.HandleFunc("/query", clue.Query(sm))
	r.HandleFunc("/print", clue.Print(sm))

	// Web page eps
	r.PathPrefix("/html/").Handler(http.StripPrefix("/html/", http.FileServer(http.Dir("./html/"))))
	r.HandleFunc("/game/{sid}/", clue.FileServe(clue.GamePage))
	r.HandleFunc("/board/{sid}/", clue.FileServe(clue.BoardPage))
	r.HandleFunc("/", clue.FileServe(clue.IndexPage))

	fmt.Println(http.ListenAndServe(ip+":8000", r))
}

/*func main() {

	deck, err := decks.FromFile("./clue/cludo.json")
	if err != nil {
		panic(err)
	}

	manager := decks.NewManager()

	manager.Insert(deck)

	session := sessions.NewManager()

	d := manager.Retrieve(0)

	state := clue.NewState(&d)

	sessionid, err := session.NewSession(state)
	if err != nil {
		panic(err)
	}

	activeState, err := clue.GetSessionState(session, sessionid)
	if err != nil {
		panic(err)
	}
	userid := activeState.Join()

	fmt.Println(sessionid, userid)

	activeState, err = clue.GetSessionState(session, sessionid)
	if err != nil {
		panic(err)
	}

	cards := []decks.CardSet{
		{
			DeckId: 0,
			CardId: 1,
		},
		{
			DeckId: 0,
			CardId: 4,
		},
		{
			DeckId: 0,
			CardId: 20,
		},
	}

	err = activeState.Op(clue.InputHands(userid, cards))
	if err != nil {
		panic(err)
	}

	fmt.Println("Got hands")

	activeState, err = clue.GetSessionState(session, sessionid)
	if err != nil {
		panic(err)
	}
	solution := activeState.Join()

	fmt.Println("Created solution")

	activeState, err = clue.GetSessionState(session, sessionid)
	if err != nil {
		panic(err)
	}
	activeState.Op(clue.ShuffleDeck())
	fmt.Println("Shuffling")

	activeState, err = clue.GetSessionState(session, sessionid)
	if err != nil {
		panic(err)
	}
	err = activeState.Op(clue.DrawSolution(solution, []string{"suspect", "room", "weapon"}))
	if err != nil {
		panic(err)
	}

	fmt.Println("Drawn solution")

	activeState, err = clue.GetSessionState(session, sessionid)
	if err != nil {
		panic(err)
	}
	u1 := activeState.Join()

	activeState, err = clue.GetSessionState(session, sessionid)
	if err != nil {
		panic(err)
	}
	u2 := activeState.Join()

	userids := []string{u1, u2}

	activeState, err = clue.GetSessionState(session, sessionid)
	if err != nil {
		panic(err)
	}
	activeState.Op(clue.DealHands(userids))

	activeState, err = clue.GetSessionState(session, sessionid)
	if err != nil {
		panic(err)
	}
	a, _ := activeState.Dump()

	fmt.Printf("%+v\n", a[u1].Deck(0))
	fmt.Printf("%+v\n", a[u2].Deck(0))
	fmt.Printf("%+v\n", a[solution].Deck(0))
	fmt.Printf("%+v\n", a[userid].Deck(0))

}*/
