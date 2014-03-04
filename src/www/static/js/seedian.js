$(document).ready(function(){
	//Manage tags
	$(".tm-input").tagsManager();
	var url = "localhost:5555"
	var tags1 = $("#tags1").tagsManager('tags');
	var tags2 = $("#tags2").tagsManager('tags');
	tags1 = ["gaming", "funny"]
	processTags(tags1, 10, "box1")
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
							items.push( "<a href='" + post.Url + "'>" + post.Title + "</a><br>" );
						});
					});
					$( "<ul/>", {
					"class": "my-new-list",
					html: items.join( "" )
					}).appendTo( "#"+id );
				}
		});
	};
});