// One preference owner, loaded before paint and shared by controls and media.
(function () {
  var root = document.documentElement;
  root.classList.add("sourcey-js");
  var script = document.currentScript;
  var key = script && script.getAttribute('data-storage-key') || 'sourcey-theme';
  var initial = script && script.getAttribute('data-default') || 'system';
  var system = window.matchMedia('(prefers-color-scheme: dark)');
  var selected;
  try { selected = localStorage.getItem(key); } catch { /* Storage may be disabled by the browser. */ }
  if (selected !== 'light' && selected !== 'dark') selected = null;
  function apply() {
    var dark = (selected || (initial === 'system' ? (system.matches ? 'dark' : 'light') : initial)) === 'dark';
    root.classList.toggle('dark', dark);
    root.style.colorScheme = dark ? 'dark' : 'light';
    var button = document.getElementById('theme-toggle');
    if (button) {
      var label = button.getAttribute(dark ? 'data-light-label' : 'data-dark-label') || (dark ? 'Switch to light mode' : 'Switch to dark mode');
      button.setAttribute('aria-label', label); button.setAttribute('title', label);
    }
    document.querySelectorAll('[data-sourcey-picture] source').forEach(function (source) {
      var value = source.getAttribute(dark ? 'data-dark' : 'data-light');
      if (value && source.getAttribute('srcset') !== value) source.setAttribute('srcset', value);
    });
    document.dispatchEvent(new CustomEvent('sourcey:themechange', { detail: { theme: dark ? 'dark' : 'light' } }));
  }
  apply();
  system.addEventListener('change', function () { if (!selected && initial === 'system') apply(); });
  window.addEventListener('storage', function (event) {
    if (event.key !== key) return;
    selected = event.newValue === 'light' || event.newValue === 'dark' ? event.newValue : null; apply();
  });
  document.addEventListener('click', function (event) {
    if (!(event.target instanceof Element) || !event.target.closest('#theme-toggle')) return;
    selected = root.classList.contains('dark') ? 'light' : 'dark';
    try { localStorage.setItem(key, selected); } catch { /* Storage may be disabled by the browser. */ }
    apply();
  });
  if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', apply, { once: true });
  else apply();
})();
