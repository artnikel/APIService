document.addEventListener('DOMContentLoaded', function () {
  var ordersModal = new bootstrap.Modal(document.getElementById('ordersModal'));

  document.getElementById('openOrdersModal').addEventListener('click', function () {
    ordersModal.show();
  });
});


function updateShares(shares) {
  var tableBody = document.getElementById('shares-table-body');
  if (shares.length > 0) {
      var newHTML = shares.map(function(share) {
          return '<tr><td>' + share.company + '</td><td>' + share.price + ' $</td></tr>';
      }).join('');
      tableBody.innerHTML = newHTML;
  } else {
      tableBody.innerHTML = '<p>No shares available.</p>';
  }
}

function fetchDataAndLog() {
  var currentTime = new Date();
  console.log('Fetching data at', currentTime);
  fetch('/getprices')
      .then(response => response.json())
      .then(data => {
          console.log('Received data at', new Date(), ':', data);
          updateShares(data);
      })
      .catch(error => {
          console.error('Error updating shares at', new Date(), ':', error);
      });
}

fetchDataAndLog();

setInterval(fetchDataAndLog, 3000);


