package main

type User struct {
	StatusCode int
	Success    bool     `json:"success"` // unneeded field, http response code covers this
	Data       UserData `json:"data,omitempty"`
	AppVersion string   `json:"appVersion"` //  "appVersion": "4.206.0" // should not be a field, could be a header
}

type UserData struct {
	Stats   UserStats `json:"stats,omitempty"`
	Flags   UserFlags `json:"flags,omitempty"`
	Armoire Armoire   `json:"armoire,omitempty"`
	//    "auth":
	//    "preferences":
	//    "_id": "6...9",
	//    "notifications":
	//    "id": "6...9"
}

type UserStats struct {
	Gold float64 `json:"gp,omitempty"` // "gp": 58845.45291132366
}

type UserFlags struct {
	// todo: if the armoire is empty?? terminate??? no, empty just means
	//       there are no fancy drops, we could just say, the thing is empty
	//       so we will just be getting food and xp for now....
	ArmoireEmpty bool `json:"armoireEmpty"`
	// todo: if armoire is not enabled, terminate
	ArmoireEnabled bool `json:"armoireEnabled"`
	ArmoireOpened  bool `json:"armoireOpened"`
}

type Armoire struct {
	Type        string `json:"type"`
	DropKey     string `json:"dropKey"`
	DropArticle string `json:"dropArticle"`
	DropText    string `json:"dropText"`
	Value       int    `json:"value"`
}

type ArmoireChecker struct {
	InitialGold   float64
	SpendingLimit float64
	DropsMap      map[string][]Armoire
	DropsCount    int64 // todo: track lifetime drops between sessions
	*Requester
}
