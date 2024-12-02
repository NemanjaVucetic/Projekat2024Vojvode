package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"projectService/domain"
)

type Client struct {
	address string
}

func NewClient(host, port string) Client {
	return Client{
		address: fmt.Sprintf("http://%s:%s", host, port),
	}
}

func (client Client) Get(id string) (*domain.User, error) {
	requestURL := client.address + "/" + id
	httpReq, err := http.NewRequest(http.MethodGet, requestURL, nil)

	if err != nil {
		log.Println(err)
		return nil, errors.New("error while getting user info")
	}

	res, err := http.DefaultClient.Do(httpReq)

	if err != nil || res.StatusCode != http.StatusOK {
		log.Println(err)
		log.Println(res.StatusCode)
		return nil, errors.New("error while getting user info")
	}

	user := &domain.User{}
	json.NewDecoder(res.Body).Decode(user)

	return user, nil
}
