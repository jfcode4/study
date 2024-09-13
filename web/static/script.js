deck = ""
decks = [];
cards = [];
index = 0

function getDecks() {
	fetch('/api/decks')
	.then(response => {
		if (!response.ok) {
			throw new Error("HTTP error: " + response.status);
		}
		return response.json();
	})
	.then(json => {
		decks = json;
		for (var deck in decks) {
			var button = document.createElement('button')
			button.innerText = decks[deck].name
			button.onclick = getStudy(decks[deck].id)
			homePage.appendChild(button)
		}
	});

}
function getStudy(deck_id) {
	return () => study(deck_id)
}

function study(deck_id) {
	fetch('/api/deck/'+deck_id)
	.then(response => {
		if (!response.ok) {
			throw new Error("HTTP error: " + response.status);
		}
		return response.json();
	});
	homePage.hidden = true;
	studyPage.hidden = false;
	deck = deck_id
	nextQuestions();
}

function showAnswer() {
	answerText.hidden = false;
	showAnswerButton.hidden = true;
	ratingButtons.hidden = false;
}
function hideAnswer() {
	answerText.hidden = true
	showAnswerButton.hidden = false;
	ratingButtons.hidden = true;
}
function rate(i) {
	id = cards[index]['id']
	hideAnswer()
	fetch(`/api/rate/${id}/${i}`)
	.then( response => {
		if (!response.ok) {
			throw new Error("HTTP error: " + response.status);
		}
	})
	.catch(error => {
		console.error("Could not rate card: " + error);
	})
	.then( () => {
		index += 1
		nextQuestions()
	})
}

function nextQuestion() {
	questionText.innerHTML = cards[index]['question'];
	answerText.innerHTML = cards[index]['answer'];
	id = cards[index]['id']
	due = cards[index]['due']
	interval = cards[index]['interval']
	tags = cards[index]['tags']
	debugText.innerHTML = `id=${id}, due=${due}, interval=${interval}, tags=${tags}, ${index+1}/${cards.length}`
}

function nextQuestions() {
	if (index >= cards.length) {
		//response = await 
		fetch('/api/study/'+deck)
		.then(response => {
			if (!response.ok) {
				throw new Error("HTTP error: " + response.status);
			}
			return response.json();
		})
		.then(json => {
			cards = json;
			index = 0;
		})
		.then(() => nextQuestion())
		.catch(error => {
			console.error("Could not get next questions: " + error);
		});
	} else {
		nextQuestion();
	}
}

// Start
getDecks()
showAnswerButton.onclick = showAnswer;
badButton.onclick = () => rate(-1);
neutralButton.onclick = () => rate(0);
goodButton.onclick = () => rate(1);


document.addEventListener("keyup", (event) => {
	if (answerText.hidden) {
		if (event.key === "Enter") {
			showAnswer();
		}
	}
	else if (event.key === "1") {
		rate(-1)
	} else if (event.key === "2" || event.key === "Enter") {
		rate(0)
	} else if (event.key === "3") {
		rate(1)
	}
});
