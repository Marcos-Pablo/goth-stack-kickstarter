document.addEventListener('click', function (e) {
  var btn = e.target.closest('[data-close-alert]');
  if (btn) {
    var alert = btn.closest('.alert');
    if (alert) {
      alert.remove();
    }
  }
});
