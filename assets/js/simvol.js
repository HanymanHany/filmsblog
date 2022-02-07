$(document).ready(function()
{
	$("#preview").keyup(function()
	{
		var box=$(this).val();
		var main = box.length *100;
		var value= (main / 800);
		var count= 800 - box.length;

	
			$('#count').html(count);
			$('#bar').animate(
			{
				"width": value+'%',
			}, 1);
		
		
	});
});