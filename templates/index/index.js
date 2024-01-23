companyData = [];

function updateShares(tableBody, shares) {
  console.log('UpdateShares - tableBody:', tableBody);
  console.log('UpdateShares - shares:', shares);

  try {
    if (!tableBody || !(tableBody instanceof HTMLElement)) {
      throw new Error('Invalid table body.');
    }

    if (!Array.isArray(shares)) {
      throw new Error('Shares data is not an array.');
    }

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
      tableBody.innerHTML = '<p>No shares available</p>';
    }
  } catch (error) {
    console.error('Error in updateShares:', error.message);
  }
}

function fetchDataAndLog(tableBodies) {
  var currentTime = new Date();
  console.log('Fetching data at', currentTime);

  fetch('/getprices')
    .then(response => {
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      return response.json();
    })
    .then(data => {
      console.log('Received data at', new Date(), ':', data);
      companyData = data;
      tableBodies.forEach(tableBody => {
        console.log('Updating shares for table:', tableBody);
        updateShares(tableBody, data);
      });
    })
    .catch(error => {
      console.error('Error fetching data:', error);
    });
}

function updateCompanyList() {
  const companyList = document.getElementById('companyList');
  if (companyList && companyList.children.length === 0) {
    const uniqueCompanies = [...new Set(companyData.map(share => share.company))];
    uniqueCompanies.forEach(company => {
      const option = document.createElement('option');
      option.value = company;
      companyList.appendChild(option);
    });
  }
}

document.addEventListener("DOMContentLoaded", function () {
  var longTable = document.getElementById('long-shares-table'); 
  var shortTable = document.getElementById('short-shares-table'); 
  var modalTable = document.getElementById('shares-table-body');
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
      fetchDataAndLog([longTable, modalTable, shortTable]);
    }
  });

  modalShort._element.addEventListener('shown.bs.modal', function () {
    if (shouldShowModalTable) {
      modalTable.classList.remove('d-none');
      fetchDataAndLog([shortTable, modalTable, longTable]);
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
  fetchDataAndLog([tableBody, longTable, shortTable, modalTable]);

  document.getElementById('companyLong').addEventListener('click', function() {
    updateCompanyList();
  });

  document.getElementById('companyShort').addEventListener('click', function() {
    updateCompanyList();
  });

  updateCompanyList();

  var companyInput = document.getElementById('companyLong');

  companyInput.addEventListener('input', function(event) {
    var enteredCompany = event.target.value;
    var companyList = document.getElementById('companyList');
    var validCompanies = [...companyList.options].map(option => option.value);

    if (!validCompanies.includes(enteredCompany)) {
      event.target.value = '';
    }
  });

  setInterval(function () {
    fetchDataAndLog([tableBody, longTable, shortTable, modalTable]);
  }, 1500);
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
                '<td><button class="copy-btn" data-dealid="' + (position.dealid || '') + '">Copy ID</button></td>' +
                '</tr>';
        }).join('');
        tableBody.innerHTML = newHTML;
        var copyButtons = document.querySelectorAll('.copy-btn');
        copyButtons.forEach(function (button) {
            button.addEventListener('click', function () {
                var dealId = button.getAttribute('data-dealid');
                copyToClipboard(dealId);
            });
        });
    } else {
        tableBody.innerHTML = '<br><p>No unclosed positions available</p>';
    }
}

async function copyToClipboard(text) {
  try {
      await navigator.clipboard.writeText(text);
      alert('Deal ID copied to clipboard: ' + text);
  } catch (err) {
      console.error('Unable to copy to clipboard', err);
      alert('Unable to copy to clipboard. Please try again or copy Deal ID manually.');
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

function validateForm(positionType) {
  var stopLossInput = document.getElementById(positionType === 'long' ? 'stoplossLong' : 'stoplossShort');
  var takeProfitInput = document.getElementById(positionType === 'long' ? 'takeprofitLong' : 'takeprofitShort');

  if (!isValidNumericInput(stopLossInput) || !isValidNumericInput(takeProfitInput)) {
      alert("Please enter valid numeric values for stop-loss and take-profit");
      return false;
  }

  var stopLoss = parseFloat(stopLossInput.value);
  var takeProfit = parseFloat(takeProfitInput.value);

  if (positionType === "long" && stopLoss >= takeProfit) {
      alert("Stop-loss should be less than take-profit");
      return false; 
  } else if (positionType === "short" && takeProfit >= stopLoss) {
      alert("Take-profit should be less than stop-loss");
      return false; 
  }
  return true; 
}

function isValidNumericInput(inputElement) {
  var value = inputElement.value.trim();
  return value !== "" && !isNaN(parseFloat(value)) && isFinite(value);
}



