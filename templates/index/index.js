document.addEventListener('DOMContentLoaded', function () {
  var ordersModal = new bootstrap.Modal(document.getElementById('ordersModal'));

  document.getElementById('openOrdersModal').addEventListener('click', function () {
    ordersModal.show();
  });
});

