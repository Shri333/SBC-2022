package main

type Group struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Items []Item `json:"items,omitempty"`
}

type Item struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Count  int    `json:"count"`
	Parent *Group `json:"group,omitempty"`
}
