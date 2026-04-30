document.addEventListener('click', function (e) {
  var btn = e.target.closest('[data-close-alert]');
  if (btn) {
    var alert = btn.closest('.alert');
    if (alert) {
      alert.remove();
    }
  }
});

(function () {
  var toggle = document.getElementById('theme-toggle');
  if (!toggle) return;

  // Set initial checkbox state
  var currentTheme = document.documentElement.getAttribute('data-theme') || 'light';
  toggle.checked = currentTheme === 'dark';

  toggle.addEventListener('change', function () {
    var theme = toggle.checked ? 'dark' : 'light';
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('theme', theme);
  });
})();
