$(document).ready(function(){
	//Manage tags
	$(".tm-input").tagsManager();
	var url = "localhost:9999"
	var tags = $(".tm-input").tagsManager('tags');
	function processTags(tags) {
	  $.post(url+"/processTags", "", 
	  	function(data, status) {
		    $("#output").append("<br>");
		    $("#output").append(data);
  });
}
});