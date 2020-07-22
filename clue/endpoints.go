package clue

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pbivrell/clue/decks"
	"github.com/pbivrell/clue/sessions"
)

const IndexPage = "./html/index.html"
const GamePage = "./html/game.html"
const BoardPage = "./html/better.html"

func FileServe(path string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}

type CreateReasponse struct {
	SessionId string `json:"sid"`
}

func Create(sm *sessions.Manager, dm *decks.Manager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		param, err := GetParameter("deck", r)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}

		deckNum, err := strconv.Atoi(param)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid deck parameter '%s':  %s", param, err), http.StatusInternalServerError)
			return

		}

		deck := dm.Retrieve(int64(deckNum))

		b := deck

		sessionid, err := sm.NewSession(NewState(b))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create new session: %s", err), http.StatusInternalServerError)
			return
		}

		resp := CreateReasponse{
			SessionId: sessionid,
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, fmt.Sprintf("Failed to create new session: %s", err), http.StatusInternalServerError)
			return
		}
	}
}

var RequiredParameterError = errors.New("required parameter")

func GetParameter(k string, r *http.Request) (string, error) {
	keys, ok := r.URL.Query()[k]

	if !ok || len(keys[0]) < 1 {
		return "", fmt.Errorf("%s '%s'", RequiredParameterError, k)
	}

	return keys[0], nil
}

type JoinResponse struct {
	UserId string `json:"uid"`
}

func Join(sm *sessions.Manager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		sessionid := vars["sid"]

		state, err := GetSessionState(sm, sessionid)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to join session: %s", err), http.StatusBadRequest)
			return
		}

		username, err := GetParameter("uid", r)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}

		userid := state.Join(username)
		resp := JoinResponse{
			UserId: userid,
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, fmt.Sprintf("Failed to join session: %s", err), http.StatusInternalServerError)
			return
		}

	}
}

type ReleaseRequest struct {
	Cards []decks.CardSet `json:"cards"`
}

func Release(sm *sessions.Manager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionid, err := GetParameter("sid", r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to release cards: %s", err), http.StatusBadRequest)
			return
		}

		userid, err := GetParameter("uid", r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to release cards: %s", err), http.StatusBadRequest)
			return
		}

		var req ReleaseRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to release cards: %s", err), http.StatusBadRequest)
			return
		}

		state, err := GetSessionState(sm, sessionid)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to release cards: %s", err), http.StatusBadRequest)
			return
		}

		err = state.Op(ReleaseHands(userid, req.Cards))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to release cards: %s", err), http.StatusBadRequest)
			return
		}
	}
}

type InputRequest struct {
	Cards []decks.CardSet `json:"cards"`
}

func Input(sm *sessions.Manager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionid, err := GetParameter("sid", r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to input cards: %s", err), http.StatusBadRequest)
			return
		}

		userid, err := GetParameter("uid", r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to input cards: %s", err), http.StatusBadRequest)
			return
		}

		var req InputRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to input cards: %s", err), http.StatusBadRequest)
			return
		}

		state, err := GetSessionState(sm, sessionid)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to input cards: %s", err), http.StatusBadRequest)
			return
		}

		err = state.Op(InputHands(userid, req.Cards))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to input cards: %s", err), http.StatusBadRequest)
			return
		}
	}
}

type RetrieveResponse struct {
	Deck decks.Deck `json:"deck"`
}

func Retrieve(sm *sessions.Manager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionid, err := GetParameter("sid", r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve cards: %s", err), http.StatusBadRequest)
			return
		}

		userid, err := GetParameter("uid", r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve cards: %s", err), http.StatusBadRequest)
			return
		}

		state, err := GetSessionState(sm, sessionid)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve cards: %s", err), http.StatusBadRequest)
			return
		}

		users, _ := state.Dump()

		user, ok := users[userid]
		if !ok {
			http.Error(w, fmt.Sprintf("Failed to retrieve cards: invalid user '%s'", userid), http.StatusBadRequest)
			return
		}

		resp := RetrieveResponse{
			Deck: user.Deck(0),
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, fmt.Sprintf("Failed to join session: %s", err), http.StatusInternalServerError)
			return
		}

	}
}

type SetupResponse struct {
	Bots     []string `json:"bots"`
	Solution string   `json:"solution"`
}

func Setup(sm *sessions.Manager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionid, err := GetParameter("sid", r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve cards: %s", err), http.StatusBadRequest)
			return
		}
		sBotCount, err := GetParameter("bots", r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve cards: %s", err), http.StatusBadRequest)
			return
		}

		botCount, err := strconv.Atoi(sBotCount)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid bot parameter '%s':  %s", sBotCount, err), http.StatusInternalServerError)
			return

		}

		bots := make([]string, botCount)

		for i := 0; i < botCount; i++ {

			state, err := GetSessionState(sm, sessionid)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to retrieve cards: %s", err), http.StatusBadRequest)
				return
			}

			bots[i] = state.Join(fmt.Sprintf("bot%d", i))
		}

		{
			state, err := GetSessionState(sm, sessionid)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to retrieve cards: %s", err), http.StatusBadRequest)
				return
			}
			state.Op(ShuffleDeck())
		}

		{
			state, err := GetSessionState(sm, sessionid)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to retrieve cards: %s", err), http.StatusBadRequest)
				return
			}
			state.Join("solution")

			state, err = GetSessionState(sm, sessionid)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to retrieve cards: %s", err), http.StatusBadRequest)
				return
			}
			err = state.Op(DrawSolution("solution", []string{"suspect", "room", "weapon"}))
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to retrieve cards: %s", err), http.StatusBadRequest)
				return
			}
		}

		{
			state, err := GetSessionState(sm, sessionid)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to retrieve cards: %s", err), http.StatusBadRequest)
				return
			}

			err = state.Op(DealHands(bots))
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to retrieve cards: %s", err), http.StatusBadRequest)
				return
			}
		}

		resp := SetupResponse{
			Bots:     bots,
			Solution: "solution",
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, fmt.Sprintf("Failed to join session: %s", err), http.StatusInternalServerError)
			return
		}

	}
}

type QueryRequest struct {
	Cards [3]decks.CardSet `json:"cards"`
}

type QueryResponse struct {
	Card decks.Card `json:"card"`
}

func Query(sm *sessions.Manager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		sessionid, err := GetParameter("sid", r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve cards: %s", err), http.StatusBadRequest)
			return
		}

		userid, err := GetParameter("uid", r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve cards: %s", err), http.StatusBadRequest)
			return
		}

		var req QueryRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to input cards: %s", err), http.StatusBadRequest)
			return
		}

		state, err := GetSessionState(sm, sessionid)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve cards: %s", err), http.StatusBadRequest)
			return
		}

		users, _ := state.Dump()

		user, ok := users[userid]
		if !ok {
			http.Error(w, fmt.Sprintf("Failed to retrieve cards: invalid user '%s'", userid), http.StatusBadRequest)
			return
		}

		rand.Shuffle(len(req.Cards), func(i, j int) { req.Cards[i], req.Cards[j] = req.Cards[j], req.Cards[i] })

		var found decks.Card
		for _, card := range user.Deck(0).Cards {
			for _, qCard := range req.Cards {
				if card.Uuid == qCard.CardId {
					found = card
				}
			}
		}

		resp := QueryResponse{
			Card: found,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, fmt.Sprintf("Failed to join session: %s", err), http.StatusInternalServerError)
			return
		}

	}
}

func Deck(sm *sessions.Manager, dm *decks.Manager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		sessionid, err := GetParameter("sid", r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve cards: %s", err), http.StatusBadRequest)
			return
		}

		state, err := GetSessionState(sm, sessionid)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve cards: %s", err), http.StatusBadRequest)
			return
		}

		_, gdeck := state.Dump()

		deckId := gdeck[0].Uuid

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(dm.Retrieve(deckId)); err != nil {
			http.Error(w, fmt.Sprintf("Failed to join session: %s", err), http.StatusInternalServerError)
			return
		}

	}
}

type UserResponse struct {
	Users []string `json:"users"`
}

func Users(sm *sessions.Manager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		sessionid, err := GetParameter("sid", r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve cards: %s", err), http.StatusBadRequest)
			return
		}

		state, err := GetSessionState(sm, sessionid)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve cards: %s", err), http.StatusBadRequest)
			return
		}

		users, _ := state.Dump()

		userKeys := []string{}
		for k, _ := range users {
			userKeys = append(userKeys, k)
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(UserResponse{Users: userKeys}); err != nil {
			http.Error(w, fmt.Sprintf("Failed to join session: %s", err), http.StatusInternalServerError)
			return
		}

	}
}

func Print(sm *sessions.Manager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Hello")
		sessionid, err := GetParameter("sid", r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve cards: %s", err), http.StatusBadRequest)
			return
		}

		state, err := GetSessionState(sm, sessionid)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve cards: %s", err), http.StatusBadRequest)
			return
		}

		fmt.Println("Hello2")
		fmt.Println(state.Dump())
	}
}
