$(document).ready(function(){
	//Manage tags
	$(".tm-input").tagsManager();
	var url = "localhost:8080"
	var tags = $("#tags").tagsManager('tags');
	var tags2 = $("#tags2").tagsManager('tags');
	processTags(tags)
	function processTags(tags) {
	  $.post("http://"+url+"/processTags", "mydata", 
	  	function(data, status) {
	  		console.log(data)
		    $("#output").append("<br>");
		    $("#output").append(data);
  });
}
});