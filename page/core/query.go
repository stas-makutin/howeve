package core

import (
	"encoding/json"

	"github.com/stas-makutin/howeve/api"
)

func StringToQuery(s string) (*api.Query, error) {
	var q api.Query
	err := json.Unmarshal([]byte(s), &q)
	return &q, err
}

func QueryToString(q *api.Query) (string, error) {
	b, err := json.Marshal(q)
	return string(b), err
}

func FetchQuery(url string, r *api.Query, then func(r *api.Query), catch func(err string)) {
	var fi *FetchInit
	if r != nil {
		body, err := QueryToString(r)
		if err != nil {
			catch("Unexpected request format: " + err.Error())
			return
		}
		fi = &FetchInit{Method: "POST", Body: body}
	}
	Fetch(url, fi, func(response *FetchResponse) {
		var errMsg string
		if response.OK {
			q, err := StringToQuery(response.Body)
			if err == nil {
				then(q)
				return
			}
			errMsg = "Unexpected response format: " + err.Error()
		} else {
			errMsg = "Unexpected response: " + response.Body
		}
		catch(errMsg)
	}, func(err *FetchError) {
		errMsg := err.Message
		if errMsg == "" {
			errMsg = "Unable to collect the requested data"
		}
		catch(errMsg)
	})
}

func FetchQueryWithSocket(r *api.Query, then func(r *api.Query), catch func(err string)) {
	request, err := QueryToString(r)
	if err != nil {
		catch("Unexpected request format: " + err.Error())
		return
	}
	socket := NewWebSocket(WebSocketUrl())
	socket.OnOpen(func() {
		socket.Send(request)
	})
	socket.OnMessage(func(data string) {
		socket.Close()
		if q, err := StringToQuery(data); err == nil {
			then(q)
		} else {
			catch("Unexpected response format: " + err.Error())
		}
	})
	socket.OnError(func() {
		socket.Close()
		catch("Unable to collect the requested data")
	})
}
