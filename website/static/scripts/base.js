function toggleSidenav() {
    document.querySelector(".sidenav2").classList.toggle("w3-show"), document.querySelector(".hamburger").classList.toggle("active")
}

!function () {
    for (var e = document.querySelectorAll(".docs table"), a = 0; a < e.length; a++) {
        var s = e[a];
        s.classList.add("w3-table", "w3-bordered");
        var t = document.createElement("div");
        t.setAttribute("class", "w3-responsive"), s.parentNode.insertBefore(t, s), t.appendChild(s)
    }
    for (var l = document.getElementsByTagName("blockquote"), a = 0; a < l.length; a++) {
        var r = l[a];
        r.classList.add("w3-panel"), r.classList.add("w3-leftbar"), r.classList.add("w3-light-grey")
    }
    for (var n = document.getElementsByTagName("code"), a = 0; a < n.length; a++) {
        var d = n[a];
        "PRE" !== d.parentElement.nodeName && d.classList.add("w3-codespan")
    }
    window.Prism && (Prism.languages.sh = Prism.languages.bash, Prism.languages.curl = Prism.languages.sh)

    const versionSelect = document.getElementById("version-select")
    versionSelect.onchange = (event) => {
        event.target.options[event.target.selectedIndex].value && (window.location = event.target.options[event.target.selectedIndex].value);
    };
}();
