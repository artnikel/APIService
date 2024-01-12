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

document.getElementById('openOrdersModal').addEventListener('click', function() {
  fetchUnclosedPositions(); 
});

function updateUnclosedPositions(positions) {
        var tableBody = document.getElementById('unclosed-positions-table-body');
        if (positions.length > 0) {
            var newHTML = positions.map(function(position) {
                return '<tr>' +
                    '<td>' + position.DealID + '</td>' +
                    '<td>' + position.SharesCount + '</td>' +
                    '<td>' + position.Company + '</td>' +
                    '<td>' + position.PurchasePrice + '$</td>' +
                    '<td>' + position.StopLoss + '$</td>' +
                    '<td>' + position.TakeProfit + '$</td>' +
                    '<td>' + position.DealTime + '</td>' +
                    '</tr>';
            }).join('');
            tableBody.innerHTML = newHTML;
        } else {
            tableBody.innerHTML = '<br><p>No unclosed positions available.</p>';
        }
}

function fetchUnclosedPositions() {
  var currentTime = new Date();
  console.log('Fetching unclosed positions at', currentTime);
  fetch('/getunclosed')
  .then(response => {
      if (!response.ok) {
          console.error('Server returned an error. Status:', response.status);
          throw new Error('Network response was not ok');
      }
      return response.json();
      })
      .then(data => {
          console.log('Received unclosed positions at', new Date(), ':', data);
          updateUnclosedPositions(data);
      })
      .catch(error => {
          console.error('Error updating unclosed positions at', new Date(), ':', error);
      });
}








