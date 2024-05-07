package types

type QuestionAnswer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func getQuestionAnswers(paramType string) []QuestionAnswer {
	questionAnswerList := map[string][]QuestionAnswer{
		"sales": {
			{Question: "Plan A", Answer: "Plan A is ..."},
			{Question: "Plan B", Answer: "Plan B is ..."},
		},
		"partner": {
			{Question: "Advertisement", Answer: "To advertise in here ..."},
			{Question: "Integration", Answer: "Here is our criteria, please contact ..."},
		},
		"support": {
			{Question: "Mobile App", Answer: "To solve there problem, please go here ..."},
			{Question: "Admin Portal", Answer: "To solve there problem, please go here ..."},
		},
		"": {
			{Question: "default", Answer: "This is default"},
		},
	}

	questionAnswers, ok := questionAnswerList[paramType]

	if !ok {
		return []QuestionAnswer{}
	}

	return questionAnswers
}
