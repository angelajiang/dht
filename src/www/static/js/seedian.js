$(document).ready(function(){
	//Manage tags
	$(".tm-input").tagsManager();
	var url = "localhost:5555"
	var tags = $("#tags").tagsManager('tags');
	var tags2 = $("#tags2").tagsManager('tags');
	processTags(tags)
	function processTags(tags) {
		$.ajax({
			url: "http://"+url+"/processTags",
			data: { inputVal: "hi"},
			success: function( data ) {
					alert( "Data Loaded: " + data );
			  		console.log("hello");
				    $("#output").append("<br>");
				    $("#output").append(data);
				},
			//headers: {"Access-Control-Allow-Origin": "http://localhost:8080"},
			crossDomain: true,
			error: function (xhr, ajaxOptions, thrownError) {
			    console.log(xhr.status);
	    		console.log(thrownError);
			},
			dataType: "jsonp"
		});
	}
});