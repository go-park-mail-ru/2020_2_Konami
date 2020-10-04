package main

func FillTestData()  {
	testUser := UserProfile{
		ImgSrc:    "testImg",
		Name:      "",
		City:      "",
		Telegram:  "",
		Vk:        "",
		Meetings:  nil,
		Interest:  "",
		Skills:    "",
		Education: "",
		Job:       "",
		Aims:      "",
	}
	card := userCard{
		CardId:    0,
		ImgSrc:    "",
		Name:      "",
		Job:       "",
		Interests: nil,
		Skills:    nil,
	}

	mapUser = make(map[int]UserProfile)
	mapSession = make(map[string]int)
	mapLoginPwd = make(map[string]string)
	UserCards = make([]userCard, 0)
	MeetingCards = make([]meetCard, 0)

	mapUser[51] = testUser
	UserCards = append(UserCards, card)
	card.CardId = 1
	UserCards = append(UserCards, card)
	card.CardId = 2
	UserCards = append(UserCards, card)
	card.CardId = 3
	UserCards = append(UserCards, card)
	card.CardId = 4
	UserCards = append(UserCards, card)
	card.CardId = 5

}
