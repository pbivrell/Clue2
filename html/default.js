function createCookie(name,value,days) {
	if (days) {
		var date = new Date();
		date.setTime(date.getTime()+(days*24*60*60*1000));
		var expires = "; expires="+date.toGMTString();
	}
	else var expires = "";
	document.cookie = name+"="+value+expires+"; path=/";
}


function readCookie(inName) {
	cookiearray = document.cookie.split(';');
        
	for(var i=0; i<cookiearray.length; i++) {
                  var name = cookiearray[i].split('=')[0].trim();
                  var value = cookiearray[i].split('=')[1];
		  if (name.localeCompare(inName) == 0) {
			return value
		  }
	}
	return "na";
}
