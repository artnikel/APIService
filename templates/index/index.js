function updateShares(tableBody, shares) {
    if (shares.length > 0) {
      var newHTML = shares.map(function(share) {
        return '<tr><td>' + share.company + '</td><td class="text-end">' + share.price + ' $</td></tr>';
      }).join('');
      tableBody.innerHTML = newHTML;
    } else {
      tableBody.innerHTML = '<p>No shares available.</p>';
    }
}

function fetchDataAndLog(tableBody) {
    var currentTime = new Date();
    console.log('Fetching data at', currentTime);
    fetch('/getprices')
      .then(response => response.json())
      .then(data => {
        console.log('Received data at', new Date(), ':', data);
        updateShares(tableBody, data);
      })
      .catch(error => {
        console.error('Error updating shares at', new Date(), ':', error);
      });
}

function toggleTableVisibility() {
    var checkbox = document.getElementById('showTableCheckbox');
    var table = document.getElementById('livePricesTable');
    var tableBody = document.getElementById('shares-table-body');

    if (checkbox.checked) {
      table.style.display = 'block';
      fetchDataAndLog(tableBody);
    } else {
      table.style.display = 'none';
    }
}

document.addEventListener("DOMContentLoaded", function() {
    var tableBody = document.getElementById('shares-table-body');
    fetchDataAndLog(tableBody);
    setInterval(function() {
      fetchDataAndLog(tableBody);
    }, 3000);
});

document.getElementById('openOrdersModal').addEventListener('click', function() {
  fetchUnclosedPositions(); 
});

document.getElementById('openHistoryModal').addEventListener('click', function() {
    fetchClosedPositions(); 
  });

function updateUnclosedPositions(positions) {
    var tableBody = document.getElementById('unclosed-positions-table-body');
    if (positions.length > 0) {
        var newHTML = positions.map(function (position) {
            return '<tr>' +
                '<td>' + (position.dealid || '') + '</td>' +
                '<td>' + (position.sharescount || '') + '</td>' +
                '<td>' + (position.company || '') + '</td>' +
                '<td>' + (position.purchaseprice ? position.purchaseprice + '$' : '') + '</td>' +
                '<td>' + (position.stoploss ? position.stoploss + '$' : '') + '</td>' +
                '<td>' + (position.takeprofit ? position.takeprofit + '$' : '') + '</td>' +
                '<td>' + (position.dealtime ? formatTimeString(position.dealtime) : '') + '</td>' +
                '<td><button class="copy-btn" data-clipboard-text="' + (position.dealid || '') + '">Copy ID</button></td>' +
                '</tr>';
        }).join('');
        tableBody.innerHTML = newHTML;
    } else {
        tableBody.innerHTML = '<br><p>No unclosed positions available</p>';
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

function updateClosedPositions(positions) {
    var tableBody = document.getElementById('closed-positions-table-body');
    if (positions.length > 0) {
        var newHTML = positions.map(function (position) {
            return '<tr>' +
                '<td>' + (position.dealid || '') + '</td>' +
                '<td>' + (position.sharescount || '') + '</td>' +
                '<td>' + (position.company || '') + '</td>' +
                '<td>' + (position.purchaseprice ? position.purchaseprice + '$' : '') + '</td>' +
                '<td>' + (position.stoploss ? position.stoploss + '$' : '') + '</td>' +
                '<td>' + (position.takeprofit ? position.takeprofit + '$' : '') + '</td>' +
                '<td>' + (position.dealtime ? formatTimeString(position.dealtime) : '') + '</td>' +
                '<td>' + (position.profit ? position.profit + '$' : '') + '</td>' +
                '<td>' + (position.enddealtime ? formatTimeString(position.enddealtime) : '') + '</td>' +
                '</tr>';
        }).join('');
        tableBody.innerHTML = newHTML;
    } else {
        tableBody.innerHTML = '<br><p>History is clear</p>';
    }
}


function fetchClosedPositions() {
  var currentTime = new Date();
  console.log('Fetching closed positions at', currentTime);
  fetch('/getclosed')
  .then(response => {
      if (!response.ok) {
          console.error('Server returned an error. Status:', response.status);
          throw new Error('Network response was not ok');
      }
      return response.json();
      })
      .then(data => {
          console.log('Received closed positions at', new Date(), ':', data);
          updateClosedPositions(data);
      })
      .catch(error => {
          console.error('Error updating closed positions at', new Date(), ':', error);
      });
}

function formatTimeString(timeString) {
    if (!timeString) {
        return ''; 
    }
    const options = { year: 'numeric', month: 'numeric', day: 'numeric', hour: 'numeric', minute: 'numeric', second: 'numeric', timeZoneName: 'short' };
    const formattedTime = new Date(timeString).toLocaleString('en-US', options);
    return formattedTime;
}




