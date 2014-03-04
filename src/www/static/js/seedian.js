$(document).ready(function(){
	//Manage tags
	$(".tm-input").tagsManager();
	var url = "localhost:5555"
	var tags1 = $("#tags1").tagsManager('tags');
	var tags2 = $("#tags2").tagsManager('tags');
	tags1 = ["gaming", 10, "funny"]
	processTags(tags1, 10, "box1")
	function processTags(tags, numTags, id) {
		$.ajax({
			type:"POST",
			dataType: "json",
			url: "http://"+url+"/processTags",
			traditional: true,
			data: {tags: tags, numTags:numTags},
			success: 
			function(data) {
					var items = [];
					$.each( data, function( key, val ) {
						items.push( "<a href='" + key + "'>" + val + "</a><br>" );
					});
					$( "<ul/>", {
					"class": "my-new-list",
					html: items.join( "" )
					}).appendTo( "#"+id );
				}
		});
	};
});