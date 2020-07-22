window.onload = function() {
	$(document).ready(function(){
		$("button").click(function(){
			$.get("/create?deck=0",
				{},
				function(data,status){
					window.location.href = "/html/game.html";
					createCookie("test","suga", 1);

				});
		});
	});

}

function createCookie(name,value,days) {
	if (days) {
		var date = new Date();
		date.setTime(date.getTime()+(days*24*60*60*1000));
		var expires = "; expires="+date.toGMTString();
	}
	else var expires = "";
	document.cookie = name+"="+value+expires+"; path=/";
}

