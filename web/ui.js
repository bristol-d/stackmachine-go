const inputarea = document.getElementById("inputarea");
const status = document.getElementById("footer-right");

function bind(id, action) {
    document.getElementById(id).onclick = action;
}

document.getElementById("close-message").onclick = function() {
    set_state("overlay", 0);
}

document.getElementById("menu-about").onclick = function() {
    document.getElementById("message").innerHTML = 
    "<p>Stack machine simulator, developed by David for the University of Bristol's 'Overview of Computer Architecture' unit.<p/>" +
    "<p>Design inspired by Windows 3.1 :)</p>"
    set_state("overlay", 1);
}

function hex(n) {
    const zeroes = '0000';
    var s = n.toString(16);
    if (s.length < 4) {
        s = zeroes.substr(0, 4 - s.length) + s
    }
    s = "0x" + s;
    return s
}

const MSTATE = [
    'OK',
    'ILLEGAL OPERATION',
    'STACK UNDERFLOW',
    'RETURN STACK UNDERFLOW',
    'STACK OVERFLOW',
    'RETURN STACK OVERFLOW',
    'ARITHMETIC ERROR',
    'HALTED'
]

document.getElementById("menu-assemble").onclick = function() {
    status.innerHTML = "Assembling ...";
    result = assemble(inputarea.value);
    status.innerHTML = "Assembly complete";
    set_state("machine", 0);
    set_state("edit", 0);
    document.getElementById("binary-container").innerHTML = result;
    dump();
}

function dump() {
    var dump = JSON.parse(dump_simulation());
    document.getElementById("controls-info").innerHTML = 
`PC: ${hex(dump.pc)}
Machine state: ${MSTATE[dump.err]}
Stack size: ${dump.n}, top: ${dump.n > 0 ? hex(dump.top) : "n/a"}
` ;
    if (dump.err > 0) {
        set_state("machine", 0);
    }
}

document.getElementById("menu-edit").onclick = function() {
    set_state("edit", 1);
}

function reset() {
    reset_simulation();
    status.innerHTML = "Ready to run."
    dump();
    set_state("machine", 1);
}

function step() {
    step_simulation();
    dump();
}

bind("menu-reset", reset);
bind("button-reset", reset);

bind("menu-step", step);
bind ("button-step", step);

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
    edit: 1,
    machine: 0, // can run?
}

window.onload = function() {

    var items = document.querySelectorAll("[x-display]");
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
        BINDINGS[key].push({val: value, elem: items[i], t: "enable"});
        if (typeof(STATE[key]) === 'undefined') {
            STATE[key] = 0;
        }
        if (STATE[key] != value) {
            items[i].style.display = "none";
        }
    }

    var d = document.querySelectorAll("[x-class]"); 
    for (i=0; i < d.length; i++) {
        data = d[i].attributes["x-class"].value.split(":");
        if (data.length != 3) {
            continue;
        }
        key = data[0];
        value = data[1];
        cls = data[2];
        if (typeof(BINDINGS[key]) === 'undefined') {
            BINDINGS[key]=[];
        }
        BINDINGS[key].push({val: value, elem: d[i], t: "class", c: cls});
        if (typeof(STATE[key]) === 'undefined') {
            STATE[key] = 0;
        }
        if (STATE[key] == value) {
            d[i].classList.add(cls);
        }
    }

    var e = document.querySelectorAll("[x-enable]");
    for (i=0; i < e.length; i++) {
        data = e[i].attributes["x-enable"].value.split(":");
        if (data.length != 2) {
            continue;
        }
        key = data[0];
        value = data[1];
        if (typeof(BINDINGS[key]) === 'undefined') {
            BINDINGS[key]=[];
        }
        BINDINGS[key].push({val: value, elem: e[i], t: "disable"});
        if (typeof(STATE[key]) === 'undefined') {
            STATE[key] = 0;
        }
        if (STATE[key] != value) {
            e[i].disabled = true;
        }
    }
};

function set_state(key, value) {
    STATE[key] = value;
    if (typeof(BINDINGS[key]) === 'undefined') {
        return;
    }
    var items = BINDINGS[key];
    for (i = 0; i < items.length; i++) {
        item = items[i];
        if (item.t == "enable") {
            item.elem.style.display = value == item.val ? "" : "none";
        } else if (item.t == "class") {
            if (value == item.val) {
                items[i].classList.add(item.c);
            } else {
                items[i].classList.remove(item.c);
            }
        } else if (item.t == "disable") {
            item.elem.disabled = !(value == item.val);
        }
    }
}