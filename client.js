/*
**	Данные соединения
*/
var SERVER_HOST = 'localhost'
var SERVER_PORT = '3000'
var HTTP_RESPONSE_OK = 200

/*
**	Устанавливаю ключ - значение в редис
*/
function SetValue() {
	var key = document.forms['redis']['key'].value;
	var value = document.forms['redis']['value'].value;

	console.log("SET tx: key=`"+ key + "` value=`" + value + "`");

	let xhr = new XMLHttpRequest();
	xhr.open("PUT", 'http://'+SERVER_HOST+':'+SERVER_PORT+'/');

	var request = {
		key: key,
		value: value
	};

	xhr.send(JSON.stringify(request));

	xhr.onload = function () {
		if (xhr.status == HTTP_RESPONSE_OK) {
			console.log('SET rx: OK')
			document.forms['redis']['key'].value = ""
			document.forms['redis']['value'].value = ""
		} else {
			console.log('SET rx: KO status=' + xhr.status)
		}
	}
	xhr.onerror = function() {
		console.log('SET rx: Error event')
	}
}

/*
**	Удаляю запись в редис по ключу
*/
function DropValue() {
	var key = document.forms['redis']['key'].value

	console.log("DROP tx: key=`"+ key + "`")

	let xhr = new XMLHttpRequest();
	xhr.open("DELETE", 'http://'+SERVER_HOST+':'+SERVER_PORT+'/');

	var request = {
		key: key
	};

	xhr.send(JSON.stringify(request));

	xhr.onload = function () {
		if (xhr.status == HTTP_RESPONSE_OK) {
			console.log('DROP rx: OK')
			document.forms['redis']['key'].value = ""
			document.forms['redis']['value'].value = ""
		} else {
			console.log('DROP rx: KO status=' + xhr.status)
		}
	}
	xhr.onerror = function() {
		console.log('DROP rx: Error event')
	}
}

/*
**	Получаю значение из редис по ключу
*/
function GetValue() {
	var key = document.forms['redis']['key'].value

	console.log("GET tx: key=`"+ key + "`")

	let xhr = new XMLHttpRequest();
	xhr.open("POST", 'http://'+SERVER_HOST+':'+SERVER_PORT+'/');

	var request = {
		key: key
	};

	xhr.send(JSON.stringify(request));

	xhr.onload = function () {
		if (xhr.status == HTTP_RESPONSE_OK) {
			let response = xhr.response
			console.log('GET rx: OK response=`' + response + '`')
			document.forms['redis']['value'].value = response
		} else {
			console.log('GET rx: KO status=' + xhr.status)
		}
	}
	xhr.onerror = function() {
		console.log('GET rx: Error event')
	}
}