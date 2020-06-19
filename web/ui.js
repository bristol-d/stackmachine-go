const inputarea = document.getElementById("inputarea");
const status = document.getElementById("footer-right");

document.getElementById("close-message").onclick = function() {
    document.getElementById("overlay").style.display = "none";
}

document.getElementById("menu-about").onclick = function() {
    document.getElementById("message").innerHTML = 
    "<p>Stack machine simulator, developed by David for the University of Bristol's 'Overview of Computer Architecture' unit.<p/>" +
    "<p>Design inspired by Windows 3.1 :)</p>"
    document.getElementById("overlay").style.display = "block";
}

document.getElementById("menu-assemble").onclick = function() {
    inputarea.style.display = "none";
    status.innerHTML = "Assembling ...";
    
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