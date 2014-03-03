$(document).ready(function(){
	//Manage tags
	$(".tm-input").tagsManager();
	var url = "localhost:5555"
	var tags1 = $("#tags1").tagsManager('tags');
	var tags2 = $("#tags2").tagsManager('tags');
	tags1 = ["test1", "test2"]
	processTags(tags1, "box1")
	function processTags(tags, id) {
		console.log(tags)
		$.ajax({
			dataType: "json",
			url: "http://"+url+"/processTags",
			traditional: true,
			data: {tags: tags},
			success: 
			function(data) {
					var items = [];
					$.each( data, function( key, val ) {
						items.push( "<a href='" + key + "'>" + val + "</a>" );
					});
					$( "<ul/>", {
					"class": "my-new-list",
					html: items.join( "" )
					}).appendTo( "#"+id );
				}
		});
	};
});