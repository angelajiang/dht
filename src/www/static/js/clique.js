$(document).ready(function(){
	//Manage tags
	/*
	$(".tm-input").tagsManager();
	var tags1 = $("#tags1").tagsManager('tags');
	*/
	var url = "localhost:5555"

	//TODO: automate event binding
	$("#tags1").keypress(
		function(e){
			if (e.which == 13){
				processTags($(this).val(), 10, "box1");
			}
		});
	$("#tags2").keypress(
		function(e){
			if (e.which == 13){
				processTags($(this).val(), 10, "box2");
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
					$.each(data, function(tag, data ) {
						//TODO: sort data by post.Ups
						$.each(data, function(index) {
							var post = data[index];
							items.push(
								"<div class='clearfix'>"+
								"<div class='ld-btn inline pull-left'>"+
								"<div class='btn-group btn-block '>"+
								"<button type='button' class='btn btn-success'>"+post.Ups+"</button>"+
								"<button type='button' class='btn btn-danger'>"+post.Downs+"</button></div></div>"+
								"&nbsp;&nbsp;"+
								"<a class='inline' href='" + post.Url + "'>" + post.Title + "</a></div>");

						});
					});
					$( "<div/>", {
					"class": "links",
					"id": id+"-links",
					html: items.join( "" )
					}).appendTo( "#"+id );
				}
		});
	};
});
