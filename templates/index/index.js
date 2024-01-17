function updateShares(tableBody, shares) {
  if (shares.length > 0) {
      var newHTML = '<thead style="font-size: 19px;">' +
                      '<tr>' +
                        '<th scope="col">Company</th>' +
                        '<th scope="col" class="text-end">Price</th>' +
                      '</tr>' +
                    '</thead>' +
                    '<tbody>' +
                      shares.map(function (share) {
                        return '<tr><td>' + share.company + '</td><td class="text-end">' + share.price + ' $</td></tr>';
                      }).join('') +
                    '</tbody>';
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

document.addEventListener("DOMContentLoaded", function () {
  var longTable = document.getElementById('long-shares-table'); // Исправлено
  var shortTable = document.getElementById('short-shares-table'); // Исправлено
  var modalTable = document.getElementById('modal-shares-table');
  var modalLong = new bootstrap.Modal(document.getElementById('longModal'));
  var modalShort = new bootstrap.Modal(document.getElementById('shortModal'));

  var shouldShowModalTable = true;

  modalLong._element.addEventListener('hidden.bs.modal', function () {
    if (!shouldShowModalTable) {
      modalTable.classList.add('d-none');
    }
  });

  modalShort._element.addEventListener('hidden.bs.modal', function () {
    if (!shouldShowModalTable) {
      modalTable.classList.add('d-none');
    }
  });

  modalLong._element.addEventListener('shown.bs.modal', function () {
    if (shouldShowModalTable) {
      modalTable.classList.remove('d-none');
      fetchDataAndLog(longTable);
    }
  });

  modalShort._element.addEventListener('shown.bs.modal', function () {
    if (shouldShowModalTable) {
      modalTable.classList.remove('d-none');
      fetchDataAndLog(shortTable);
    }
  });

  window.addEventListener('click', function (event) {
    var modalElementLong = document.getElementById('longModal'); 
    var modalElementShort = document.getElementById('shortModal'); 

    if (event.target == modalElementLong && shouldShowModalTable) {
      modalLong.hide();
    }

    if (event.target == modalElementShort && shouldShowModalTable) {
      modalShort.hide();
    }
  });

  var tableBody = document.getElementById('shares-table-body');
  fetchDataAndLog(tableBody);

  setInterval(function () {
    fetchDataAndLog(tableBody);
    if (shouldShowModalTable) {
      fetchDataAndLog(longTable);
      fetchDataAndLog(shortTable);
    }
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




