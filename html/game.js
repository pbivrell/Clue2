var selected = [];

var deck = {
	"uuid": 0,
	"name": "cluedo",
	"description": "Clue board from the cluedo set",
	"cards": [
		{
			"uuid": 0,
			"category": "suspect",
			"value": "Colonel Mustard"
		},{
			"uuid": 1,
			"category": "suspect",
			"value": "Professor Plum"
		}, {
			"uuid": 2,
			"category": "suspect",
			"value":"Reverend Green" 
		},{
			"uuid": 3,
			"category": "suspect",
			"value":"Mrs Peacock"
		},{
			"uuid": 4,
			"category": "suspect",
			"value":"Miss Scarlett"
		},{
			"uuid": 5,
			"category": "suspect",
			"value":"Mrs White"
		},{
			"uuid": 7,
			"category": "weapon",
			"value":"Dagger"
		},{
			"uuid": 8,
			"category": "weapon",
			"value":"Candlestick" 
		},{
			"uuid": 9,
			"category": "weapon",
			"value":"Revolver"
		},{
			"uuid": 10,
			"category": "weapon",
			"value":"Rope"
		},{
			"uuid": 11,
			"category": "weapon",
			"value":"Lead Pipe"
		},{
			"uuid": 12,
			"category": "weapon",
			"value":"Spanner"
		},{
			"uuid": 13,
			"category": "room",
			"value":"Hall" 
		},{
			"uuid": 14,
			"category": "room",
			"value": "Lounge" 
		},{
			"uuid": 15,
			"category": "room",
			"value":"Dining Room" 
		},{
			"uuid": 16,
			"category": "room",
			"value": "Kitchen" 
		},{
			"uuid": 17,
			"category": "room",
			"value":"Ballroom" 
		},{
			"uuid": 18,
			"category": "room",
			"value": "Conservatory" 
		},{
			"uuid": 19,
			"category": "room",
			"value": "Billiard Room" 
		},{
			"uuid": 20,
			"category": "room",
			"value":"Library" 
		},{
			"uuid": 21,
			"category": "room",
			"value":"Study"
		}
	]
}

var sid;
var uid;

window.onload = function() {
	deck.cards.forEach(myFunction);

	sid=(window.location + '').split('/')[4];
	uid=readCookie(sid);
	if (uid==="na") {
		window.location.href = "/?sid=" + sid;
	}
	document.getElementById("uid").innerHTML = uid;
}

function insert(item) {
	if (selected.includes(item)) {
		removeItem(selected, item);
	} else {
		selected.push(item);
	}
	var text = ""
	for(var i in selected) {
		text += " " + selected[i].value + ",";
	}
	document.getElementById("selected").innerText = text;
}

function removeItem(array, item){
	for(var i in array){
		if(array[i]==item){
			array.splice(i,1);
			break;
		}
	}
}


function myFunction(item, index) {

	var btn = document.createElement("button");               
	btn.innerText = item.value;               
	btn.onclick = function () { insert(item);  };
	document.getElementById("cards").appendChild(btn); 
}

function onclickSubmit() {
	var deck = [];
	for(var i in selected) {
		deck.push({
			deckId: 0,
			cardId: selected[i].uuid,
		});
	}
	$.ajax("/input?sid="+sid+"&uid="+uid, {
		data : JSON.stringify({
			cards: deck,

		}),
		contentType : 'application/json',
		type : 'POST',
		error: function(XMLHttpRequest, textStatus, errorThrown) { 
			alert("Status: " + textStatus); alert("Error: " + errorThrown); 
		},
		success: function(){

		} ,
	})

}
