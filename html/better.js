var deck;
var sid;
var uid;

window.onload = function() {

		// Get session id from path
		sid=(window.location + '').split('/')[4];
		
		uid=readCookie(sid);
		if (uid==="na") {
				window.location.href = "/?sid=" + sid;
		}
		document.getElementById("user").innerHTML = uid;

		$.get("/deck?sid="+sid,
				{},
				function(data, status) {
						console.log("data",data);
						deck = data;
						console.log("deck",deck);
						displayBoard();
						displayAsk();
						$.get("/get?sid="+sid+"&uid="+uid, {},
								function(data, status) {
									console.log(data);
										for (var i in data.deck.cards) {
												document.getElementById(data.deck.cards[i].uuid).classList.toggle("owned");
										}
								});
				});


}

function setupOnClick() {
	var bots = parseInt(document.getElementById("setup-text").value);
        $.get("/setup?sid="+sid+"&bots="+bots,
                {},
                function(data, status) {
			console.log(data);
		}
	);
}

function displayAsk() {
        var ask = document.getElementById("people");
        $.get("/users?sid="+sid,
                {},
                function(data, status) {
                        for(var i in data.users) {
                                var newDiv = document.createElement("div");
                                newDiv.innerHTML = data.users[i];
                                newDiv.classList.add("user");
				newDiv.onclick = newUserOnClick(newDiv);
                                ask.appendChild(newDiv);
                        }
                });


}

function newUserOnClick(div) {

        return function() {
                div.classList.toggle("selected-user");
        }
}


function releaseOnClick() {
	var selected = [];
	var selectedItems = document.getElementsByClassName('selected');
	for (var i = 0; i < selectedItems.length; ++i) {
		selected.push( { "deckId": deck.uuid, "cardId": parseInt(selectedItems[i].id)});
	}

        console.log(selected);
        $.ajax("/release?sid="+sid+"&uid="+uid, {
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

function inputCardsOnClick() {
	var selected = [];
	var selectedItems = document.getElementsByClassName('selected');
	for (var i = 0; i < selectedItems.length; ++i) {
		selected.push( { "deckId": deck.uuid, "cardId": parseInt(selectedItems[i].id)});
	}

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


function queryCardsOnClick() {
	var selected = [];
	var selectedItems = document.getElementsByClassName('selected');
	for (var i = 0; i < selectedItems.length; ++i) {
		selected.push( { "deckId": deck.uuid, "cardId": parseInt(selectedItems[i].id)});
	}
        $.ajax("/query?sid="+sid+"&uid="+document.getElementsByClassName('selected-user')[0].innerHTML, {
                data : JSON.stringify({
                        cards: selected,

                }),
                contentType : 'application/json',
                type : 'POST',
                error: function(XMLHttpRequest, textStatus, errorThrown) {
                        alert("Status: " + textStatus); alert("Error: " + errorThrown);
                },
                success: function(data){
                        console.log(data);
                        var found = document.getElementById(data.card.uuid);
                        found.classList.add("found");
                },
        });
}


function displayBoard() {
        console.log(deck);
        let categories = new Map()
        for (var i in deck.cards) {
                categories.set(deck.cards[i].category, true)
        }

        var div = document.getElementById("cards");

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

function cardOnClick(deck, card, div) {

        return function() {
		div.classList.remove("found");
                div.classList.toggle("selected");
        }
}
