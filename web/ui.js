const inputarea = document.getElementById("inputarea");
const status = document.getElementById("footer-right");

document.getElementById("close-message").onclick = function() {
    set_state("overlay", 0);
}

document.getElementById("menu-about").onclick = function() {
    document.getElementById("message").innerHTML = 
    "<p>Stack machine simulator, developed by David for the University of Bristol's 'Overview of Computer Architecture' unit.<p/>" +
    "<p>Design inspired by Windows 3.1 :)</p>"
    set_state("overlay", 1);
}

document.getElementById("menu-assemble").onclick = function() {
    status.innerHTML = "Assembling ...";
    result = assemble(inputarea.value);
    status.innerHTML = "Assembly complete";
    set_state("edit", 0);
    document.getElementById("binary-container").innerHTML = result;
}

document.getElementById("menu-edit").onclick = function() {
    set_state("edit", 1);
}

function cursorPosition() {
    start = inputarea.selectionStart;
    end = inputarea.selectionEnd;
    pos = start;
    if (start != end) {
        pos = inputarea.selectionDirection == "forward" ? end : start;
    }
    val = inputarea.value;
    line = 1;
    col = 1;
    for (i = 0; i < pos; i++) {
        if (val[i] == "\n") {
            line += 1;
            col = 1;
        } else {
            col += 1;
        }
    }
    status.innerHTML = "Line " + line + " column " + col;
}
inputarea.onkeyup = cursorPosition;
inputarea.onmousedown = cursorPosition;
inputarea.ontouchstart = cursorPosition;
inputarea.oninput = cursorPosition;
inputarea.onpaste = cursorPosition;
inputarea.oncut = cursorPosition;
inputarea.onselect = cursorPosition;
inputarea.onselectstart = cursorPosition;

BINDINGS = {

}

STATE = {
    edit: 1
}

items = document.querySelectorAll("[x-display]");
for (i=0; i < items.length; i++) {
    data = items[i].attributes["x-display"].value.split(":");
    if (data.length != 2) {
        continue;
    }
    key = data[0];
    value = data[1];
    if (typeof(BINDINGS[key]) === 'undefined') {
        BINDINGS[key]=[];
    }
    BINDINGS[key].push({val: value, elem: items[i]});
    if (typeof(STATE[key]) === 'undefined') {
        STATE[key] = 0;
    }
    if (STATE[key] != value) {
        items[i].style.display = "none";
    }
}

function set_state(key, value) {
    STATE[key] = value;
    if (typeof(BINDINGS[key]) === 'undefined') {
        return;
    }
    items = BINDINGS[key];
    for (i = 0; i < items.length; i++) {
        item = items[i];
        item.elem.style.display = value == item.val ? "" : "none";
    }
}