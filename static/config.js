////////////////////////////////////////////////////////////
// Selector Control
////////////////////////////////////////////////////////////
var selected_value = "#0";
$("table:not("+selected_value+")").hide();

$("#selector").change(function() {
	selected_value = '#'+$("#selector").val();
	$(selected_value).show();
	$("table:not("+selected_value+")").hide();
});

////////////////////////////////////////////////////////////
// Controller Animations
////////////////////////////////////////////////////////////
$("input").mouseleave(function(){$("#image").attr("src", "static/hes.png");});
$(".A").mouseover(function(){$("#image").attr("src", "static/hes_A.png");});
$(".B").mouseover(function(){$("#image").attr("src", "static/hes_B.png");});
$(".Left").mouseover(function(){$("#image").attr("src", "static/hes_Left.png");});
$(".Right").mouseover(function(){$("#image").attr("src", "static/hes_Right.png");});
$(".Up").mouseover(function(){$("#image").attr("src", "static/hes_Up.png");});
$(".Down").mouseover(function(){$("#image").attr("src", "static/hes_Down.png");});
$(".Select").mouseover(function(){$("#image").attr("src", "static/hes_Select.png");});
$(".Start").mouseover(function(){$("#image").attr("src", "static/hes_Start.png");});

////////////////////////////////////////////////////////////
// Input Hotkeys
////////////////////////////////////////////////////////////

function add_to_input(focused, hotkey) {
	if(!$("#"+focused.id).val().includes(hotkey))
		$("#"+focused.id).val($("#"+focused.id).val()+hotkey);
}

var keys_dict = {};
keys_dict[9]="TAB";
keys_dict[13]="ENTER";
keys_dict[20]="CAPSLOCK";
keys_dict[27]="ESC";
keys_dict[32]="SPACE";
keys_dict[33]="PGUP";
keys_dict[34]="PGDN";
keys_dict[37]="LEFT";
keys_dict[38]="UP";
keys_dict[39]="RIGHT";
keys_dict[40]="DOWN";
keys_dict[45]="INS";
keys_dict[46]="DEL";
keys_dict[112]="F1";
keys_dict[144]="NUMLOCK";



$("input").click(function() {
	var focused = document.activeElement;
	var hotkey;
	$("#"+focused.id).keydown(function(event){
		//event.val().toUpperCase();
		if(event.ctrlKey) {
			event.preventDefault();
			hotkey = "CTRL";
		}
		else if(event.altKey) {
			event.preventDefault();
			hotkey = "ALT";
		}
		else if(event.shiftKey) {
			event.preventDefault();
			hotkey = "SHIFT";
		}
		else if(event.metaKey) {
			event.preventDefault();
			hotkey = "META";
		}
		else if(event.keyCode == 38) {
			event.preventDefault();
			hotkey = "UP";
		}
		else if(keys_dict[event.keyCode] != null) {
			event.preventDefault();
			hotkey = keys_dict[event.keyCode];
		}
		else return;
		add_to_input(focused, hotkey);
	});
});

