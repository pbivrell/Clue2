var deck;
var sid;
var uid;

window.onload = function() {

	sid=(window.location + '').split('/')[4];

	$.get("/deck?sid="+sid,
		{},
		function(data, status) {
			console.log("data",data);
			deck = data;
			console.log("deck",deck);
			displayBoard();
			$.get("/get?sid="+sid+"&uid="+uid, {},
				function(data, status) {
					for (var i in data.deck.cards) {
						document.getElementById(data.deck.cards[i].uuid).classList.toggle("owned");
					}

					if(data.deck.cards.length === 0) {
						var ask = document.getElementById("askToggle");
						ask.style.display = "none";
						var submit = document.getElementById("submit");
						submit.onclick = inputCardsOnClick;
					}else{
						var submit = document.getElementById("submit");
						submit.style.display = "none";
					}
				});
		});

	$("#askToggle").click(function(){ 
		displayAsk();
	}, function(){});

	uid=readCookie(sid);
	if (uid==="na") {
		window.location.href = "/?sid=" + sid;
	}

}

function inputCardsOnClick() {
	console.log(selected);
	$.ajax("/input?sid="+sid+"&uid="+uid, {
		data : JSON.stringify({
			cards: selected,

		}),
		contentType : 'application/json',
		type : 'POST',
		error: function(XMLHttpRequest, textStatus, errorThrown) {
			alert("Status: " + textStatus); alert("Error: " + errorThrown);
		},
		success: function(){
			location.reload();
		},
	});
}

var selected = [];

function addOrRemove(array, value) {

	var index = -1;
	for (var i in array) {
		if (array[i].deckId == value.deckId && array[i].cardId == value.cardId) {
			index = i;
			break;
		}
	}

	if (index === -1) {
		array.push(value);
	} else {
		array.splice(index, 1);
	}
}

function cardOnClick(deck, card, div) {

	return function() {
		div.classList.toggle("selected");
		addOrRemove(selected, { "deckId": deck, "cardId": card });
	}
}

var queryUser;

function newUserOnClick(userid) {

	
	var enable = document.getElementById("enable");
	//enable.innerHTML= "? - " + userid.trim();

	return function() {
		queryUser = userid;
		var submit = document.getElementById("submit");
		submit.style.display = "flex";
		submit.onclick = queryCardsOnClick;
	}
}

function undo(name) {
	const elements = document.getElementsByClassName(name);
	while(elements.length > 0){
	    elements[0].classList.remove(name);
	}
}

function queryCardsOnClick() {
	console.log("Apple", selected);
	$.ajax("/query?sid="+sid+"&uid="+queryUser, {
		data : JSON.stringify({
			cards: selected,

		}),
		contentType : 'application/json',
		type : 'POST',
		error: function(XMLHttpRequest, textStatus, errorThrown) {
			alert("Status: " + textStatus); alert("Error: " + errorThrown);
		},
		success: function(data){
			undo("selected");
			undo("found");
			//highlightCards(data.card, "blue");
			console.log(data);
			var found = document.getElementById(data.card.uuid);
			found.classList.add("found");
			var submit = document.getElementById("submit");
			submit.style.display = "none";
			selected = [];
		},
	});
}

function displayAsk() {
	var ask = document.getElementById("ask");
	ask.innerHTML = "";
	$.get("/users?sid="+sid,
		{},
		function(data, status) {
			for(var i in data.users) {
				var newDiv = document.createElement("div");
				newDiv.innerHTML = data.users[i];
				newDiv.classList.add("user");
				newDiv.onclick = newUserOnClick(data.users[i]);
				ask.appendChild(newDiv);
			}
		});


}

function displayBoard() {
	console.log(deck);
	let categories = new Map()
	for (var i in deck.cards) {
		categories.set(deck.cards[i].category, true)
	}

	var div = document.getElementById("content");

	for (let [key, value] of categories) {

		var newDiv = document.createElement("div");
		newDiv.id = key;
		newDiv.innerHTML = key + "s";
		newDiv.classList.add("column");

		for (var i in deck.cards) {
			if(deck.cards[i].category === key) {
				var newDiv2 = document.createElement("div");
				newDiv2.id = deck.cards[i].uuid;
				newDiv2.classList.add("card");
				newDiv2.innerHTML = deck.cards[i].value;
				console.log("Hello");
				newDiv2.onclick = cardOnClick(deck.uuid, deck.cards[i].uuid, newDiv2)
				newDiv.appendChild(newDiv2);

			}
		}

		div.appendChild(newDiv);



	}
}
