var selected_value = "#0";
$("table:not("+selected_value+")").hide();

$("#selector").change(function() {
	selected_value = '#'+$("#selector").val();
	$(selected_value).show();
	$("table:not("+selected_value+")").hide();
});

$("input").mouseleave(function(){$("#image").attr("src", "static/hes.png");});
$(".A").mouseover(function(){$("#image").attr("src", "static/hes_A.png");});
$(".B").mouseover(function(){$("#image").attr("src", "static/hes_B.png");});
$(".Left").mouseover(function(){$("#image").attr("src", "static/hes_Left.png");});
$(".Right").mouseover(function(){$("#image").attr("src", "static/hes_Right.png");});
$(".Up").mouseover(function(){$("#image").attr("src", "static/hes_Up.png");});
$(".Down").mouseover(function(){$("#image").attr("src", "static/hes_Down.png");});
$(".Select").mouseover(function(){$("#image").attr("src", "static/hes_Select.png");});
$(".Start").mouseover(function(){$("#image").attr("src", "static/hes_Start.png");});