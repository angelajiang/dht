$(document).ready(function(){
	//Manage tags
	/*
	$(".tm-input").tagsManager();
	var tags1 = $("#tags1").tagsManager('tags');
	*/
	var url = "localhost:5555"
	//$("#tags1").submit(processTags(tags1, 5, "box1"))
	var tags1 = ["chicago", "Northwestern"]
	$("#tags1").keypress(
		function(e){
			console.log('here')
			if (e.which == 13){
				processTags(tags1, 10, "links1");
			}
		});
	function processTags(tags, numLinks, id) {
		$.ajax({
			type:"POST",
			dataType: "json",
			url: "http://"+url+"/processTags",
			traditional: true,
			data: {tags: tags, numLinks:numLinks},
			success: 
			function(data) {
					var items = [];
					$.each(data, function( tag, data ) {
						$.each(data, function(index) {
							var post = data[index];
							items.push( "<button type='button' class='btn btn-success btn-xs'>like</button>"+
								"&nbsp;&nbsp;<a href='" + post.Url + "'>" + post.Title + "</a><br>");

						});
					});
					$( "<div/>", {
					"class": "new-list",
					html: items.join( "" )
					}).appendTo( "#"+id );
				}
		});
	};
});
