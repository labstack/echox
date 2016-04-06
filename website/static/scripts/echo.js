(function() {
  var menu = document.querySelectorAll('.menu a');

  for (var i = 0; i < menu.length; i++) {
    var m = menu[i];
    if (location.href === m.href) {
      m.className += 'active';
    }
  }

  // Add table class
  var tables = document.getElementsByTagName('table');

  for (var i = 0; i < tables.length; i++) {
    var t = tables[i];
    t.classList.add('mdl-data-table');
  }
})();

function search(e) {
  if (e.which === 13) {
    google.search.cse.element.getElement('standard0').execute(e.target.value);
  }
}
